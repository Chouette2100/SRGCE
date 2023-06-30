package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
)

func InsertIntoOrUpdateUser(
	tnow time.Time, eventid string, roominf exsrapi.RoomInfo,
) (
	err error,
) {

	isnew := false

	userno, _ := strconv.Atoi(roominf.ID)
	//	log.Printf("  *** InsertIntoOrUpdateUser() *** userno=%d\n", userno)

	//	レコードがすでに存在するか？
	nrow := 0
	err = srdblib.Db.QueryRow("select count(*) from " + srdblib.Tuser + " where userno =" + roominf.ID).Scan(&nrow)

	if err != nil {
		err = fmt.Errorf("QueryRow(): %w", err)
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
		//	存在しない。

		isnew = true

		//	log.Printf("insert into " + srdblib.Tuserhistory + "(*new*) userno=%d rank=<%s> nrank=<%s> prank=<%s> level=%d, followers=%d, fans=%d, fans_lst=%d\n",
		//		userno, roominf.Rank, roominf.Nrank, roominf.Prank, roominf.Level, roominf.Followers, fans, fans_lst)

		sqli := "INSERT INTO " + srdblib.Tuser + " (userno, userid, user_name, longname, shortname, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, ts, currentevent)"
		sqli += " VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

		//	log.Printf("sql=%s\n", sql)
		var stmti *sql.Stmt
		stmti, err = srdblib.Db.Prepare(sqli)
		if err != nil {
			//	log.Printf("InsertIntoOrUpdateUser() error() (INSERT/Prepare) err=%s\n", err.Error())
			err = fmt.Errorf("Prepare(): %w", err)
			return
		}
		defer stmti.Close()

		lenid := len(roominf.ID)
		_, err = stmti.Exec(
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
			//	log.Printf("error(InsertIntoOrUpdateUser() INSERT/Exec) err=%s\n", err.Error())
			//	status = -2
			_, err = stmti.Exec(
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
				//	log.Printf("error(InsertIntoOrUpdateUser() INSERT/Exec) err=%s\n", err.Error())
				err = fmt.Errorf("Exec(): %w", err)
			}
		}
	} else {
		//	存在する。
		sqls := "select user_name, genre, `rank`, nrank, prank, level, followers, fans, fans_lst from " + srdblib.Tuser + "  where userno = ?"
		err = srdblib.Db.QueryRow(sqls, userno).Scan(&name, &genre, &rank, &nrank, &prank, &level, &followers, &fans, &fans_lst)
		if err != nil {
			log.Printf("err=[%s]\n", err.Error())
			err = fmt.Errorf("WueryRow().Scan(): %w", err)
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

			//	log.Printf("insert into userhistory(*changed*) userno=%d level=%d, followers=%d, fans=%d\n",
			//		userno, roominf.Level, roominf.Followers, roominf.Fans)
			sqlu := "update " + srdblib.Tuser + " set userid=?,"
			sqlu += "user_name=?,"
			sqlu += "genre=?,"
			sqlu += "`rank`=?,"
			sqlu += "nrank=?,"
			sqlu += "prank=?,"
			sqlu += "level=?,"
			sqlu += "followers=?,"
			sqlu += "fans=?,"
			sqlu += "fans_lst=?,"
			sqlu += "ts=?,"
			sqlu += "currentevent=? "
			sqlu += "where userno=?"
			var stmtu *sql.Stmt
			stmtu, err = srdblib.Db.Prepare(sqlu)

			if err != nil {
				log.Printf("InsertIntoOrUpdateUser() error(Update/Prepare) err=%s\n", err.Error())
				err = fmt.Errorf("Prepare(): %w", err)
				return
			}
			defer stmtu.Close()

			_, err = stmtu.Exec(
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
				err = fmt.Errorf("Exec(): %w", err)
			}
		}
		/* else {
			//	log.Printf("not insert into userhistory(*same*) userno=%d level=%d, followers=%d\n", userno, roominf.Level, roominf.Followers)
		}
		*/

	}

	if isnew {
		sqli := "INSERT INTO " + srdblib.Tuserhistory + "(userno, user_name, genre, `rank`, nrank, prank, level, followers, fans, fans_lst, ts)"
		sqli += " VALUES(?,?,?,?,?,?,?,?,?,?,?)"
		//	log.Printf("sql=%s\n", sql)
		var stmti *sql.Stmt
		stmti, err = srdblib.Db.Prepare(sqli)
		if err != nil {
			log.Printf("error(INSERT into userhistory/Prepare) err=%s\n", err.Error())
			err = fmt.Errorf("(userhistory) Prepare(): %w", err)
			return
		}
		defer stmti.Close()

		_, err = stmti.Exec(
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
			_, err = stmti.Exec(
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
				//	log.Printf("error(Insert Into into userhistory INSERT/Exec) err=%s\n", err.Error())
				err = fmt.Errorf("Exec(): %w", err)
			}
		}

	}

	return

}
