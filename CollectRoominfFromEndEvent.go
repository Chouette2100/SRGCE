package main

import (
	"fmt"
	"log"

	//	"strings"
	//	"strconv"
	"time"

	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi/v2"
	//	"github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

type Uinf struct {
	Userno   int
	Username string
	Point    int
	Rank     int
}

// 結果が発表されたイベントの順位と獲得ポイントを取得する
// これはイベント終了日の翌日12時から翌々日12時までのあいだに行う必要がある)
// 前日終了のイベントのデータを取得するか、前々日のものを取得するかは実行時刻に応じて判断される。
func CollectRoominfFromEndEvent(
	client *http.Client,
	tevent string,
	teventuser string,
	tuser string,
	tuserhistory string,
) (
	err error,
) {

	fn := exsrapi.PrtHdr()
	defer exsrapi.PrintExf("", fn)()

	var stmt *sql.Stmt
	var rows *sql.Rows

	tnow := time.Now().Truncate((time.Second)) // 秒より下を切り捨てる。時刻の比較を整数の比較と同等にする。
	hh, _, _ := tnow.Clock()

	//	現時点で結果が表示されるイベントの終了時刻の範囲を求める。
	//	tnow.Truncate(24 * time.Hour)はUTCで計算されている模様。
	//	特にg時まではtnow.Truncate(24 * time.Hour)の結果は1日前になることに注意
	var tday1, tday2 time.Time
	if hh < 9 {
		//	12時までは前々日に終了したイベントの結果が表示されるが、9時までは前日相当となる。
		tday1 = tnow.Truncate(24 * time.Hour).Add(-33 * time.Hour)
		tday2 = tnow.Truncate(24 * time.Hour).Add(-9 * time.Hour)
	} else if hh < 12 {
		//	12時までは前々日に終了したイベントの結果が表示されている。
		tday1 = tnow.Truncate(24 * time.Hour).Add(-57 * time.Hour)
		//	tday1 = tnow.Truncate(24 * time.Hour).Add(-153 * time.Hour)
		tday2 = tnow.Truncate(24 * time.Hour).Add(-33 * time.Hour)
	} else {
		//	12時をすぎると前日に終了したイベントの結果が表示される。
		tday1 = tnow.Truncate(24 * time.Hour).Add(-33 * time.Hour)
		//	day1 = tnow.Truncate(24 * time.Hour).Add(-57 * time.Hour)
		tday2 = tnow.Truncate(24 * time.Hour).Add(-9 * time.Hour)
	}
	log.Printf("tday:  %s\n", tnow.Format("2006-01-02 15:04:05 MST"))
	log.Printf("tday1: %s\n", tday1.Format("2006-01-02 15:04:05 MST"))
	log.Printf("tday2: %s\n", tday2.Format("2006-01-02 15:04:05 MST"))

	sqlstmt := "select eventid from " + tevent + " where achk = 0 and endtime > ? and endtime < ?"
	stmt, err = srdblib.Db.Prepare(sqlstmt)
	if err != nil {
		err = fmt.Errorf("srdblib.Db.Prepare(): %w", err)
		return
	}
	defer stmt.Close()

	rows, err = stmt.Query(tday1, tday2)
	if err != nil {
		err = fmt.Errorf("stmt.Query(): %w", err)
		return
	}
	defer rows.Close()

	idofevent := make([]string, 0)

	eid := ""
	for rows.Next() {
		err = rows.Scan(&eid)
		if err != nil {
			err = fmt.Errorf("rows.Scan(): %w", err)
			return
		}
		idofevent = append(idofevent, eid)
	}

	log.Printf("==================================\n")
	for _, eid = range idofevent {
		log.Printf("eventid: %s\n", eid)

		// =============================================== ここから　CollectOneRoominfFromEndEvent()

		//	err = CollectOneRoominfFromEndEvent(client, tevent, teventuser, tuser, tuserhistory, tnow, eid)
		err = CollectOneRoominfFromEndEvent(client, tnow, eid)
		if err != nil {
			err = fmt.Errorf("CollectOneRoominfFromEndEvent(): %w", err)
			continue
		}

		// =============================================== ここまで　CollectOneRoominfFromEndEvent()

		/*
			//	取得すべきデータの存在チェック（取得済みかのチェック）
			nrow := 0
			sqlsc := "select count(*) from " + teventuser + " where eventid = ?"
			srdblib.Db.QueryRow(sqlsc, eid).Scan(&nrow)
			if nrow > 0 {
				//	取得済み
				log.Printf("    data exists. skip.\n")
				continue
			}

			//	イベントの詳細を得る、ここではIeventidが必要である
			row, err := srdblib.Dbmap.Get(srdblib.Event{}, eid)
			if err != nil {
				err = fmt.Errorf("Dbmap.Get(): %w", err)
				return err
			}
			event := row.(*srdblib.Event)

			//	イベントに参加しているルームを取得する
			roomlistinf, err := srapi.GetRoominfFromEventByApi(client, event.Ieventid, 1, 1)
			if err != nil {
				err = fmt.Errorf("GetRoominfFromEventByApi(): %w", err)
				return err
			}
			roomid := roomlistinf.RoomList[0].Room_id

			//	イベント結果を取得する
			bid := 0
			if strings.Contains(event.Eventid, "block_id") {
				eida := strings.Split(event.Eventid, "=")
				bid, _ = strconv.Atoi(eida[1])
			}
			pranking, err := srapi.ApiEventsRanking(client, (event).Ieventid, roomid, bid)
			if err != nil {
				err = fmt.Errorf("ApiEventsRanking(): %w", err)
				return err
			}

			uinflist := make([]Uinf, 0)

			if len(pranking.Ranking) != 0 {
				for _, ranking := range pranking.Ranking {
					uinflist = append(uinflist, Uinf{Userno: ranking.Room.RoomID, Username: ranking.Room.Name, Point: ranking.Point, Rank: ranking.Rank})
				}
			} else {
				var eventinf exsrapi.Event_Inf
				var roominfolist exsrapi.RoomInfoList
				//	if !strings.Contains(id, "?") {
				exsrapi.GetEventinfAndRoomList(eid, 1, 30, &eventinf, &roominfolist)
				//	} else {
				//		exsrapi.GetEventinfAndRoomListBR(client, id, 1, 30, &eventinf, &roominfolist)
				//	}
				for _, roominf := range roominfolist {
					uinflist = append(uinflist, Uinf{Userno: roominf.Userno, Username: roominf.Name, Point: roominf.Point, Rank: roominf.Irank})
				}
			}

			for _, uinf := range uinflist {
				status := ""
				status, err = CreateEventuserFromEventinf(teventuser, eid, uinf)
				if err != nil {
					err = fmt.Errorf("CreateEventuserFromEventinf(): %w", err)
					status = "**error"
					log.Printf("  %-10s %-25s%10d%4d%10d %s\n", status, uinf.Username, uinf.Userno, uinf.Rank, uinf.Point, eid)
					return err
				} else {
					log.Printf("  %-10s %-25s%10d%4d%10d %s\n", status, uinf.Username, uinf.Userno, uinf.Rank, uinf.Point, eid)
					if status == "ignored." {
						continue
					}
				}
				//		err = InsertIntoOrUpdateUser(tuser, tuserhistory, tnow, id, ranking)
				wuser := new(srdblib.Wuser)
				wuser.Userno = uinf.Userno
				err = srdblib.UpinsWuserSetProperty(client, tnow, wuser, 1440 * 5, 1000)
				if err != nil {
					err = fmt.Errorf("InsertIntoOrUpdateUser(): %w", err)
					return err
				}
			}
		*/

		//	========================================================= ここまで CollectOneRoominfFromEndEvent()
	}

	return
}
