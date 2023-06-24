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

func GetRoominfAll() (
	err error,
) {
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

	tnow := time.Now().Truncate((time.Second))	// 秒より下を切り捨てる。時刻の比較を整数の比較と同等にする。
	tday1 := tnow.Truncate(24 * time.Hour).Add(-57 * time.Hour)
	tday2 := tnow.Truncate(24 * time.Hour).Add(-33 * time.Hour)
	//	tday1 := tnow.Truncate(24 * time.Hour).Add(-33 * time.Hour)
	//	tday2 := tnow.Truncate(24 * time.Hour).Add(-9 * time.Hour)
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

	idofeventbox := make([]string, 0)

	id := ""
	for rows.Next() {
		srdblib.Dberr = rows.Scan(&id)
		if srdblib.Dberr != nil {
			err = fmt.Errorf("rows.Scan(): %w", srdblib.Dberr)
			return
		}
		idofeventbox = append(idofeventbox, id)
	}

	log.Printf("==================================\n")
	for _, id := range idofeventbox {
		var eventinf exsrapi.Event_Inf
		var roominfolist RoomInfoList
		if ! strings.Contains(id, "?") {
			GetEventInfAndRoomList(id, 1, 30, &eventinf, &roominfolist)
		} else {
			GetEventInfAndRoomListBR(client, id, 1, 30, &eventinf, &roominfolist)
		}
		for _, room := range roominfolist {
			InsertIntoEventuser(eventinf.Event_ID, room)
			InsertIntoOrUpdateUser(tnow, eventinf.Event_ID, room)
			log.Printf("%s %s %d %d %d\n", eventinf.Event_ID, room.Account, room.Userno, room.Irank, room.Point)
		}
	}

	return
}
