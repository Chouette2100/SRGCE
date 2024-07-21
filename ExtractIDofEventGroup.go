package main

import (
	"fmt"
	"log"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
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

	sqlstmt := "select eventid from " + tevent + " where achk = ?"
	stmt, srdblib.Dberr = srdblib.Db.Prepare(sqlstmt)
	if srdblib.Dberr != nil {
		err = fmt.Errorf("row.Priepare(): %w", srdblib.Dberr)
		return
	}
	defer stmt.Close()

	rows, srdblib.Dberr = stmt.Query(mode)
	if srdblib.Dberr != nil {
		err = fmt.Errorf("stmt.Query(): %w", srdblib.Dberr)
		return
	}
	defer rows.Close()

	idofeventgroup = make([]string, 0)

	id := ""
	for rows.Next() {
		srdblib.Dberr = rows.Scan(&id)
		if srdblib.Dberr != nil {
			err = fmt.Errorf("rows.Scan(): %w", srdblib.Dberr)
			return
		}
		log.Printf(" parent id = %s\n", id)
		idofeventgroup = append(idofeventgroup, id)
	}
	return
}
