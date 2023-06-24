package main
import (
	"strconv"
	"log"

	"github.com/Chouette2100/srdblib"
)
func InsertIntoEventuser(eventid string, roominf RoomInfo) (status int) {

	status = 0

	userno, _ := strconv.Atoi(roominf.ID)

	nrow := 0
	/*
		sql := "select count(*) from eventuser where "
		sql += "userno =" + roominf.ID + " and "
		//	sql += "eventno = " + fmt.Sprintf("%d", eventno)
		sql += "eventid = " + eventid
		//	log.Printf("sql=%s\n", sql)
		err := Db.QueryRow(sql).Scan(&nrow)
	*/
	sql := "select count(*) from " + srdblib.Teventuser + " where userno =? and eventid = ?"
	err := srdblib.Db.QueryRow(sql, roominf.ID, eventid).Scan(&nrow)

	if err != nil {
		log.Printf("select count(*) from user ... err=[%s]\n", err.Error())
		status = -1
		return
	}
	/*
	Colorlist := Colorlist2
	if Event_inf.Cmap == 1 {
		Colorlist = Colorlist1
	}
	*/

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
		/*
			} else {
				_, err = stmt.Exec(
					eventid,
					userno,
					"Y",	//	"N"から変更する＝順位に関わらず獲得ポイントデータを取得する。
					"N",
					Colorlist[i%len(Colorlist)].Name,
					"N",
					roominf.Point,
				)
			}
		*/

		if err != nil {
			log.Printf("error(InsertIntoOrUpdateUser() INSERT/Exec) err=%s\n", err.Error())
			status = -2
		}
		status = 1
	}
	return

}
