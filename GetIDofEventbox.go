package main

import (
	"fmt"
	"log"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/srdblib"
)

func GetIDofEventbox() (
	idofeventbox []string,
	err error,
) {

	var stmt *sql.Stmt
	var rows *sql.Rows

	sqlstmt := "select eventid from " + srdblib.Tevent + " where achk = 6"
	stmt, srdblib.Dberr = srdblib.Db.Prepare(sqlstmt)
	if srdblib.Dberr != nil {
		err = fmt.Errorf("row.Priepare(): %w", srdblib.Dberr)
		return
	}
	defer stmt.Close()

	rows, srdblib.Dberr = stmt.Query()
	if srdblib.Dberr != nil {
		err = fmt.Errorf("stmt.Query(): %w", srdblib.Dberr)
		return
	}
	defer rows.Close()

	idofeventbox = make([]string, 0)

	id := ""
	for rows.Next() {
		srdblib.Dberr = rows.Scan(&id)
		if srdblib.Dberr != nil {
			err = fmt.Errorf("rows.Scan(): %w", srdblib.Dberr)
			return
		}
		log.Printf(" parent id = %s\n", id)
		idofeventbox = append(idofeventbox, id)
	}
	return
}
