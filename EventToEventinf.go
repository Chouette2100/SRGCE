package main

import (
	"time"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srapi"
	//	"github.com/Chouette2100/srdblib"
)

/*
ApiEventSearch()で得られるイベント情報構造体 srapi.Eventを
データベースのeventテーブルに対応したexsrapi.Event_Infに変換する
*/
func EventToEventinf(
	event *srapi.Event,	// srapi.ApiEventSearch()で得られるイベント構造体
) (
	Eventinf *exsrapi.Event_Inf,	// exsrapi、srdblibで用いるイベント構造体
) {
	eventinf := exsrapi.Event_Inf{}	// ローカルな変数だが外部から参照可能であるあいだは存在する。

	eventinf.Event_ID = event.Event_url_key	//	イベントページURLの最後のフィールド、イベントIDとして扱う
	eventinf.I_Event_ID = event.Event_id	//	こちらが本来のイベントIDらしいが、あまり目にすることはないが、/api/event/block_rankingには必要jj
	eventinf.Event_name = event.Event_name
	eventinf.Start_time = time.Unix(event.Started_at, 0)
	eventinf.Sstart_time = eventinf.Start_time.Format("2006-01-02 15:04:05")	//	互換性のため
	eventinf.End_time = time.Unix(event.Ended_at, 0)
	eventinf.Send_time = eventinf.End_time.Format("2006-01-02 15:04:05")	//	互換性のため

	eventinf.Rstatus = "NowSaved"

	if event.Is_event_Block {
		eventinf.Achk = 5		//	子ブロックが展開されていないイベントボックス、展開されたら 1 とする。
	} else if event.Is_box_event {
		eventinf.Achk = 6	//	子ブロックが展開されていないブロックイベント、展開されたら 2 とする。
	}

	return &eventinf
}
