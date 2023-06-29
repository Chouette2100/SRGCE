package main

import (
	"fmt"
	"log"
	"time"
	//	"sort"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srapi"
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

func CreateCurrentEventList() (
	top *OngoingEvent,
	err error,
) {

	fn := exsrapi.PrtHdr()
	defer exsrapi.PrintExf("", fn)()

	client, cookiejar, err := exsrapi.CreateNewClient("")
	if err != nil {
		log.Printf("exsrapi.CeateNewClient(): %s", err.Error())
		return //	エラーがあれば、ここで終了
	}
	defer cookiejar.Save()

	// テンプレートに埋め込むデータ（ポイントやランク）を作成する
	top = new(OngoingEvent)
	top.TimeNow = time.Now().Unix()

	top.Eventlist, err = srapi.MakeEventListByApi(client)
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
