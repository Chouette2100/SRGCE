// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package main

import (
	"fmt"
	"log"
	//	"strconv"
	//	"strings"
	"time"

	//	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi"
	//	"github.com/Chouette2100/srapi"
	"github.com/Chouette2100/srdblib"
)

func CollectOneRoominfFromEndEvent(
	client *http.Client,
	//	tevent string,
	//	teventuser string,
	//	tuser string,
	//	tuserhistory string,
	tnow time.Time,
	eid string,
) (
	err error,
) {

	//	srdblib.Dbmap.AddTableWithName(srdblib.Wuser{}, tuser).SetKeys(false, "Userno")
	//	srdblib.Dbmap.AddTableWithName(srdblib.Userhistory{}, tuserhistory).SetKeys(false, "Userno", "Ts")
	//	srdblib.Dbmap.AddTableWithName(srdblib.Event{}, "wevent").SetKeys(false, "Eventid")
	//	srdblib.Dbmap.AddTableWithName(srdblib.Eventuser{}, teventuser).SetKeys(false, "Eventid", "Userno")

	log.Printf("eventid: %s\n", eid)

	// 取得すべきデータの存在チェック（取得済みかのチェック）
	nrow := 0
	//	sqlsc := "select count(*) from " + teventuser + " where eventid = ?"
	sqlsc := "select count(*) from wevent where eventid = ?"
	srdblib.Db.QueryRow(sqlsc, eid).Scan(&nrow)
	if nrow > 0 {
		//	取得済み
		log.Printf("    data exists. skip.\n")
		err = fmt.Errorf("CollectOneRoominfFromEndEvent(): data exists")
		return
	}

	//	====================================== ここから GetEventsRankingByApi()

	pranking, err := srdblib.GetEventsRankingByApi(client, eid, 2)
	if err != nil {
		err = fmt.Errorf("GetEventsRankingByApi(): %w", err)
		return err
	}

	/*
	// イベントの詳細を得る、ここではIeventidが必要である
	row, err := srdblib.Dbmap.Get(srdblib.Event{}, eid)
	if err != nil {
		err = fmt.Errorf("Dbmap.Get(): %w", err)
		return err
	}
	event := row.(*srdblib.Event)

	// イベントに参加しているルームを一つ取得する
	//	(ApiEventsRanking()にはイベントIDとともにルームIDが必要だから)
	roomlistinf, err := srapi.GetRoominfFromEventByApi(client, event.Ieventid, 1, 1)
	if err != nil {
		err = fmt.Errorf("GetRoominfFromEventByApi(): %w", err)
		return err
	}
	roomid := roomlistinf.RoomList[0].Room_id

	// イベント結果を取得する
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
	*/

	//	====================================== ここまで GetEventsRankingByApi()

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
		// status, err = CreateEventuserFromEventinf(teventuser, eid, uinf)
		status, err = CreateEventuserFromEventinf(eid, uinf)
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
	return
}
