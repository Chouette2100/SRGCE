package main
import (
	"strconv"
	"log"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
)
func InsertIntoEventuser(eventid string, roominf exsrapi.RoomInfo) (status int) {

	status = 0

	userno, _ := strconv.Atoi(roominf.ID)

	nrow := 0
	sql := "select count(*) from " + srdblib.Teventuser + " where userno =? and eventid = ?"
	err := srdblib.Db.QueryRow(sql, roominf.ID, eventid).Scan(&nrow)
	if err != nil {
		log.Printf("select count(*) from user ... err=[%s]\n", err.Error())
		status = -1
		return
	}

	if nrow == 0 {
		sql := "INSERT INTO " + srdblib.Teventuser + "(eventid, userno, point, vld) VALUES(?,?,?,?)"
		stmt, err := srdblib.Db.Prepare(sql)
		if err != nil {
			log.Printf("error(INSERT/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}
		defer stmt.Close()

		//	if i < 10 {
		_, err = stmt.Exec(
			eventid,
			userno,
			roominf.Point,
			roominf.Irank,
		)

		if err != nil {
			log.Printf("error(InsertIntoOrUpdateUser() INSERT/Exec) err=%s\n", err.Error())
			status = -2
		}
		status = 1
	}
	return

}
