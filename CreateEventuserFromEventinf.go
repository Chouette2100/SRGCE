package main

import (
	"fmt"
	"log"
	//	"strconv"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	//	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

func CreateEventuserFromEventinf(
	//	teventuser string,
	eventid string,
	//	roominf exsrapi.RoomInfo,
	uinf Uinf,
) (
	status string,
	err error,
) {

	//	fn := exsrapi.PrtHdr()
	//	defer exsrapi.PrintExf("", fn)()

	//	userno, _ := strconv.Atoi(roominf.ID)
	userno := uinf.Userno

	//	レコードがすでに存在するか？
	nrow := 0
	status = "ignored."
	//	sqls := "select count(*) from " + teventuser + " where userno =? and eventid = ?"
	sqls := "select count(*) from weventuser where userno =? and eventid = ?"
	//	err = srdblib.Db.QueryRow(sqls, roominf.ID, eventid).Scan(&nrow)
	err = srdblib.Db.QueryRow(sqls, userno, eventid).Scan(&nrow)
	if err != nil {
		//	log.Printf("select count(*) from user ... err=[%s]\n", err.Error())
		err = fmt.Errorf("QueryRow().Scan(): %w", err)
		return
	}

	if nrow == 0 {
		//	存在しない。
		var stmti *sql.Stmt
		//	sqli := "INSERT INTO " + teventuser + "(eventid, userno, point, vld) VALUES(?,?,?,?)"
		sqli := "INSERT INTO weventuser (eventid, userno, point, vld) VALUES(?,?,?,?)"
		stmti, err = srdblib.Db.Prepare(sqli)
		if err != nil {
			//	log.Printf("error(INSERT/Prepare) err=%s\n", err.Error())
			err = fmt.Errorf("Prepare(): %w", err)
			return
		}
		defer stmti.Close()

		//	if i < 10 {
		_, err = stmti.Exec(
			eventid,
			userno,
			//	roominf.Point,
			uinf.Point,
			//	roominf.Irank,
			uinf.Rank,
		)

		if err != nil {
			log.Printf("error(InsertIntoOrUpdateUser() INSERT/Exec) err=%s\n", err.Error())
			err = fmt.Errorf("Exec(): %w", err)
		}
		status = "inserted."
	}
	return

}
