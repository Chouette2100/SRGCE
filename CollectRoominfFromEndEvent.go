package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
)

func CollectRoominfFromEndEvent() (
	err error,
) {

	fn := exsrapi.PrtHdr()
	defer exsrapi.PrintExf("", fn)()

	var stmt *sql.Stmt
	var rows *sql.Rows

	//      cookiejarがセットされたHTTPクライアントを作る
	client, jar, err := exsrapi.CreateNewClient("ShowroomCGI")
	if err != nil {
		log.Printf("CreateNewClient: %s\n", err.Error())
		return
	}
	//      すべての処理が終了したらcookiejarを保存する。
	defer jar.Save()

	tnow := time.Now().Truncate((time.Second)) // 秒より下を切り捨てる。時刻の比較を整数の比較と同等にする。
	hh, _, _ := tnow.Clock()

	var tday1, tday2 time.Time
	if hh < 12 {
		//	12時までは前々日に終了したイベントの結果が表示されている。
		tday1 = tnow.Truncate(24 * time.Hour).Add(-57 * time.Hour)
		tday2 = tnow.Truncate(24 * time.Hour).Add(-33 * time.Hour)

	} else {
		//	12時をすぎると前日に終了したイベントの結果が表示される。
		tday1 = tnow.Truncate(24 * time.Hour).Add(-33 * time.Hour)
		tday2 = tnow.Truncate(24 * time.Hour).Add(-9 * time.Hour)

	}
	log.Printf("tday:  %s\n", tnow.Format("2006-01-02 15:04:05 MST"))
	log.Printf("tday1: %s\n", tday1.Format("2006-01-02 15:04:05 MST"))
	log.Printf("tday2: %s\n", tday2.Format("2006-01-02 15:04:05 MST"))

	sqlstmt := "select eventid from " + srdblib.Tevent + " where achk = 0 and endtime > ? and endtime < ?"
	stmt, srdblib.Dberr = srdblib.Db.Prepare(sqlstmt)
	if srdblib.Dberr != nil {
		err = fmt.Errorf("row.Priepare(): %w", srdblib.Dberr)
		return
	}
	defer stmt.Close()

	rows, srdblib.Dberr = stmt.Query(tday1, tday2)
	if srdblib.Dberr != nil {
		err = fmt.Errorf("stmt.Query(): %w", srdblib.Dberr)
		return
	}
	defer rows.Close()

	idofevent := make([]string, 0)

	id := ""
	for rows.Next() {
		srdblib.Dberr = rows.Scan(&id)
		if srdblib.Dberr != nil {
			err = fmt.Errorf("rows.Scan(): %w", srdblib.Dberr)
			return
		}
		idofevent = append(idofevent, id)
	}

	log.Printf("==================================\n")
	for _, id := range idofevent {
		log.Printf("eventid: %s\n", id)
		var eventinf exsrapi.Event_Inf
		var roominfolist exsrapi.RoomInfoList
		if !strings.Contains(id, "?") {
			exsrapi.GetEventinfAndRoomList(id, 1, 30, &eventinf, &roominfolist)
		} else {
			exsrapi.GetEventinfAndRoomListBR(client, id, 1, 30, &eventinf, &roominfolist)
		}
		for _, room := range roominfolist {
			err = CreateEventuserFromEventinf(eventinf.Event_ID, room)
			if err != nil {
				err = fmt.Errorf("CreateEventuserFromEventinf(): %w", err)
				return
			}
			status := ""
			status, err = InsertIntoOrUpdateUser(tnow, eventinf.Event_ID, room)
			if err != nil {
				err = fmt.Errorf("InsertIntoOrUpdateUser(): %w", err)
				return
			}
			log.Printf("  %-10s %-25s%10d%4d%10d %s\n", status, room.Account, room.Userno, room.Irank, room.Point, eventinf.Event_ID)
		}
	}

	return
}
