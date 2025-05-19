package main

import (
	"fmt"
	"log"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

func StoreEventinflistInEvent(
	tevent string,
	eventinflist []exsrapi.Event_Inf,
	) (
	err error,
) {

	fn := exsrapi.PrtHdr()
	defer exsrapi.PrintExf("", fn)()

	var stmts, stmtu *sql.Stmt

	//	既存データの変化をチェックする必要があるカラムの抽出用SQL
	sqls := "select endtime, period, noentry, achk from " + tevent + " where eventid = ?"
	stmts, srdblib.Dberr = srdblib.Db.Prepare(sqls)
	if srdblib.Dberr != nil {
		err = fmt.Errorf("Prepare(sqls): %w", srdblib.Dberr)
		return
	}
	defer stmts.Close()

	//	データが変更されたカラムの更新用SQL
	sqlu := "UPDATE " + tevent + " SET endtime = ?, period = ?, noentry = ?, achk = ? WHERE eventid = ?"
	stmtu, srdblib.Dberr = srdblib.Db.Prepare(sqlu)
	if srdblib.Dberr != nil {
		err = fmt.Errorf("Prepare(sqlu): %w", srdblib.Dberr)
		return
	}
	defer stmtu.Close()

	var endtime time.Time
	var noentry, achk int
	var period string
	for i, eventinf := range eventinflist {
		if eventinf.Event_ID == "safaripark_showroom" {
			//	block_id=0 が存在することに対する一時的回避処理
			continue
		}
		//	存在確認
		srdblib.Dberr = stmts.QueryRow(eventinf.Event_ID).Scan(&endtime, &period, &noentry, &achk)
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
			reason := ""
			if eventinf.End_time.Sub(endtime)  < time.Second * -1 || eventinf.End_time.Sub(endtime) > time.Second{
				//	終了時刻が変更された。
				reason += "E"
			} else {
				reason += " "
			}
			if eventinf.Period != period {
				//	期間が変更（修正）された。
				reason += "P"
			} else {
				reason += " "
			}
			if eventinf.NoEntry != noentry {
				//	イベント参加者数が変化した。
				reason += "N"
			} else {
				reason += " "
			}
			if eventinf.Achk%4 != achk {
				//	イベント種別（ブロックイベント、イベントボックス、それら以外）が変化した。
				reason += "A"
			} else {
				reason += " "
			}
			if reason != "    " {
				if eventinf.Achk%4 == achk {
					//	ここでこの条件が成り立つのはendtime, noentry, period が変化したケースでイベントグループの子イベントの登録が終わっている場合。
					//	その場合子イベントの登録状態（achk < 4)を変更してはならない。Achkを現在の状態に戻してから保存する。
					eventinf.Achk = achk
				}
				stmtu.Exec(eventinf.End_time, eventinf.Period, eventinf.NoEntry, eventinf.Achk, eventinf.Event_ID)
				log.Printf("  **Updated[%s]: %-30s %s\n", reason, eventinf.Event_ID, eventinf.Event_name)

			} else {
				//	イベント状態に変化がない。
				log.Printf("  **Ignored: %-30s %s\n", eventinf.Event_ID, eventinf.Event_name)
			}
		}
	}

	if len(eventinflist) != 0 {
		err = srdblib.InsertEventinflistToEvent(tevent, &eventinflist, false)
		if err != nil {
			err = fmt.Errorf("InsertEventinflistToEvent(): %w", srdblib.Dberr)
			return
		}
	}

	return
}
