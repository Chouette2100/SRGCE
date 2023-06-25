package main
import (
	"time"
	"log"
	"strconv"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"

)
func InsertIntoOrUpdateUser(tnow time.Time, eventid string, roominf exsrapi.RoomInfo) (status int) {

	status = 0

	isnew := false

	userno, _ := strconv.Atoi(roominf.ID)
	log.Printf("  *** InsertIntoOrUpdateUser() *** userno=%d\n", userno)

	nrow := 0
	err := srdblib.Db.QueryRow("select count(*) from " + srdblib.Tuser + " where userno =" + roominf.ID).Scan(&nrow)

	if err != nil {
		log.Printf("select count(*) from user ... err=[%s]\n", err.Error())
		status = -1
		return
	}

	name := ""
	genre := ""
	rank := ""
	nrank := ""
	prank := ""
	level := 0
	followers := 0
	fans := -1
	fans_lst := -1

	if nrow == 0 {

		isnew = true

		log.Printf("insert into " + srdblib.Tuserhistory + "(*new*) userno=%d rank=<%s> nrank=<%s> prank=<%s> level=%d, followers=%d, fans=%d, fans_lst=%d\n",
			userno, roominf.Rank, roominf.Nrank, roominf.Prank, roominf.Level, roominf.Followers, fans, fans_lst)

		sql := "INSERT INTO " +srdblib.Tuser + " (userno, userid, user_name, longname, shortname, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, ts, currentevent)"
		sql += " VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

		//	log.Printf("sql=%s\n", sql)
		stmt, err := srdblib.Db.Prepare(sql)
		if err != nil {
			log.Printf("InsertIntoOrUpdateUser() error() (INSERT/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}
		defer stmt.Close()

		lenid := len(roominf.ID)
		_, err = stmt.Exec(
			userno,
			roominf.Account,
			roominf.Name,
			//	roominf.ID,
			roominf.Name,
			roominf.ID[lenid-2:lenid],
			roominf.Genre,
			roominf.Rank,
			roominf.Nrank,
			roominf.Prank,
			roominf.Level,
			roominf.Followers,
			roominf.Fans,
			roominf.Fans_lst,
			tnow,
			eventid,
		)

		if err != nil {
			log.Printf("error(InsertIntoOrUpdateUser() INSERT/Exec) err=%s\n", err.Error())
			//	status = -2
			_, err = stmt.Exec(
				userno,
				roominf.Account,
				roominf.Account,
				roominf.ID,
				roominf.ID[lenid-2:lenid],
				roominf.Genre,
				roominf.Rank,
				roominf.Nrank,
				roominf.Prank,
				roominf.Level,
				roominf.Followers,
				roominf.Fans,
				roominf.Fans_lst,
				tnow,
				eventid,
			)
			if err != nil {
				log.Printf("error(InsertIntoOrUpdateUser() INSERT/Exec) err=%s\n", err.Error())
				status = -2
			}
		}
	} else {

		sql := "select user_name, genre, `rank`, nrank, prank, level, followers, fans, fans_lst from " + srdblib.Tuser + "  where userno = ?"
		err = srdblib.Db.QueryRow(sql, userno).Scan(&name, &genre, &rank, &nrank, &prank, &level, &followers, &fans, &fans_lst)
		if err != nil {
			log.Printf("err=[%s]\n", err.Error())
			status = -1
		}
		//	log.Printf("current userno=%d name=%s, nrank=%s, prank=%s level=%d, followers=%d\n", userno, name, nrank, prank, level, followers)

		if roominf.Genre != genre ||
			roominf.Rank != rank ||
			//	roominf.Nrank != nrank ||
			//	roominf.Prank != prank ||
			roominf.Level != level ||
			roominf.Followers != followers ||
			roominf.Fans != fans {

			isnew = true

			log.Printf("insert into userhistory(*changed*) userno=%d level=%d, followers=%d, fans=%d\n",
				userno, roominf.Level, roominf.Followers, roominf.Fans)
			sql := "update " + srdblib.Tuser + " set userid=?,"
			sql += "user_name=?,"
			sql += "genre=?,"
			sql += "`rank`=?,"
			sql += "nrank=?,"
			sql += "prank=?,"
			sql += "level=?,"
			sql += "followers=?,"
			sql += "fans=?,"
			sql += "fans_lst=?,"
			sql += "ts=?,"
			sql += "currentevent=? "
			sql += "where userno=?"
			stmt, err := srdblib.Db.Prepare(sql)

			if err != nil {
				log.Printf("InsertIntoOrUpdateUser() error(Update/Prepare) err=%s\n", err.Error())
				status = -1
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(
				roominf.Account,
				roominf.Name,
				roominf.Genre,
				roominf.Rank,
				roominf.Nrank,
				roominf.Prank,
				roominf.Level,
				roominf.Followers,
				roominf.Fans,
				roominf.Fans_lst,
				tnow,
				eventid,
				roominf.ID,
			)

			if err != nil {
				log.Printf("error(InsertIntoOrUpdateUser() Update/Exec) err=%s\n", err.Error())
				status = -2
			}
		}
		/* else {
			//	log.Printf("not insert into userhistory(*same*) userno=%d level=%d, followers=%d\n", userno, roominf.Level, roominf.Followers)
		}
		*/

	}

	if isnew {
		sql := "INSERT INTO " + srdblib.Tuserhistory + "(userno, user_name, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, ts)"
		sql += " VALUES(?,?,?,?,?,?,?,?,?,?,?)"
		//	log.Printf("sql=%s\n", sql)
		stmt, err := srdblib.Db.Prepare(sql)
		if err != nil {
			log.Printf("error(INSERT into userhistory/Prepare) err=%s\n", err.Error())
			status = -1
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			userno,
			roominf.Name,
			roominf.Genre,
			roominf.Rank,
			roominf.Nrank,
			roominf.Prank,
			roominf.Level,
			roominf.Followers,
			roominf.Fans,
			roominf.Fans_lst,
			tnow,
		)

		if err != nil {
			log.Printf("error(Insert Into into userhistory INSERT/Exec) err=%s\n", err.Error())
			//	status = -2
			_, err = stmt.Exec(
				userno,
				roominf.Account,
				roominf.Genre,
				roominf.Rank,
				roominf.Nrank,
				roominf.Prank,
				roominf.Level,
				roominf.Followers,
				roominf.Fans,
				roominf.Fans_lst,
				tnow,
			)
			if err != nil {
				log.Printf("error(Insert Into into userhistory INSERT/Exec) err=%s\n", err.Error())
				status = -2
			}
		}

	}

	return

}
