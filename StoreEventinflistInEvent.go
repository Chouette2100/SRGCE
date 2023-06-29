package main

import (
	"fmt"
	//	"log"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
)

func StoreEventinflistInEvent(eventinflist []exsrapi.Event_Inf) (
	err error,
) {

	fn := exsrapi.PrtHdr()
	defer exsrapi.PrintExf("", fn)()

	var stmts, stmtu *sql.Stmt

	sqls := "select endtime, noentry, achk from " + srdblib.Tevent + " where eventid = ?"
	stmts, srdblib.Dberr = srdblib.Db.Prepare(sqls)
	if srdblib.Dberr != nil {
		err = fmt.Errorf("Prepare(sqls): %w", srdblib.Dberr)
		return
	}
	defer stmts.Close()

	sqlu := "UPDATE " + srdblib.Tevent + " SET endtime = ?, noentry = ?, achk = ? WHERE eventid = ?"
	stmtu, srdblib.Dberr = srdblib.Db.Prepare(sqlu)
	if srdblib.Dberr != nil {
		err = fmt.Errorf("Prepare(sqlu): %w", srdblib.Dberr)
		return
	}
	defer stmtu.Close()

	var endtime time.Time
	var noentry, achk int
	for i, eventinf := range eventinflist {
		//	存在確認
		srdblib.Dberr = stmts.QueryRow(eventinf.Event_ID).Scan(&endtime, &noentry, &achk)
		switch {
		case srdblib.Dberr == sql.ErrNoRows:
			//	存在しない。
			//	後続の処理でinsertする。
			eventinflist[i].Valid = true
		case srdblib.Dberr != nil:
			//	エラー
			err = fmt.Errorf("QueryRow(event.Event_ID).Scan(): %w", srdblib.Dberr)
			return
		default:
			//	存在する。endtime、achkが違うならupdateする。
			if eventinf.End_time != endtime || eventinf.NoEntry != noentry ||eventinf.Achk%4 != achk {
				if eventinf.Achk%4 == achk {
					//	ここでこの条件が成り立つのはendtime, noentry が変化したケースでイベントグループの子イベントの登録が終わっている場合。
					//	その場合子イベントの登録状態（achk < 4)を変更してはならない。
					eventinf.Achk = achk
				}
				stmtu.Exec(eventinf.End_time, eventinf.NoEntry, eventinf.Achk, eventinf.Event_ID)
			}
		}
	}

	if len(eventinflist) != 0 {
		err = srdblib.InsertEventinflistToEvent(&eventinflist, false)
		if err != nil {
			err = fmt.Errorf("InsertEventinflistToEvent(): %w", srdblib.Dberr)
			return
		}
	}

	return
}
