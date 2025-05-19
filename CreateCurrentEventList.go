package main

import (
	"fmt"
	"log"
	"time"
	//	"sort"
	"net/http"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srapi/v2"
)

type OngoingEvent struct {
	TimeNow    int64
	Totalcount int
	ErrMsg     string
	Eventlist  []srapi.Event
}

/*
var Db *sql.DB
var Err error
*/

/*
	開催中あるいは開催予定のイベントの一覧を取得する。
	CreateCurrentEventList(status)
	 status: 1: 開催中(デフォルト)、 3: 開催予定、 4: 終了済み
*/
func CreateCurrentEventList(
	client *http.Client,
	status int,
) (
	top *OngoingEvent,
	err error,
) {

	fn := exsrapi.PrtHdr()
	defer exsrapi.PrintExf("", fn)()

	// テンプレートに埋め込むデータ（ポイントやランク）を作成する
	top = new(OngoingEvent)
	top.TimeNow = time.Now().Unix()

	top.Eventlist, err = srapi.MakeEventListByApi(client, status)
	if err != nil {
		err = fmt.Errorf("MakeListOfPoints(): %w", err)
		log.Printf("MakeListOfPoints() returned error %s\n", err.Error())
		top.ErrMsg = err.Error()
	}
	top.Totalcount = len(top.Eventlist)

	//	ソートが必要ないときは次の行とimportの"sort"をコメントアウトする。
	//	無名関数のリターン値でソート条件を変更できます。
	//	sort.Slice(top.Eventlist, func(i, j int) bool { return top.Eventlist[i].Ended_at > top.Eventlist[j].Ended_at })

	log.Printf("  CreateCurrentEventList() returned %d events\n", top.Totalcount)

	return

}
