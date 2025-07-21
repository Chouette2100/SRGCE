package main

import (
	"fmt"
	"log"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

const EventBox = 6
const BlockEvent = 5

func ExtractIDofEventGroup(
	tevent string,
	mode int,
) (
	idofeventgroup []string,
	err error,
) {

	fn := exsrapi.PrtHdr()
	defer exsrapi.PrintExf("", fn)()

	var stmt *sql.Stmt
	var rows *sql.Rows

	sqlstmt := "select eventid from " + tevent + " where starttime > Now() and (achk = ? or achk = ?) "
	if targetEvent != "" {
		// 特定のブロックイベントの展開を行う場合
		sqlstmt = "select eventid from " + tevent + " where eventid = ? "
	}

	stmt, err = srdblib.Db.Prepare(sqlstmt)
	if err != nil {
		err = fmt.Errorf("row.Priepare(): %w", err)
		return
	}
	defer stmt.Close()

	if targetEvent == "" {
		rows, err = stmt.Query(mode, mode&3)
	} else {
		// 特定のブロックイベントの展開を行う場合
		rows, err = stmt.Query(targetEvent)
	}
	if err != nil {
		err = fmt.Errorf("stmt.Query(): %w", err)
		return
	}
	defer rows.Close()

	idofeventgroup = make([]string, 0)

	id := ""
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			err = fmt.Errorf("rows.Scan(): %w", err)
			return
		}
		log.Printf(" parent id = %s\n", id)
		idofeventgroup = append(idofeventgroup, id)
	}
	return
}
