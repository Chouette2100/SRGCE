package main

import (
	//	"database/sql"
	"fmt"
	"log"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
)

func ExpandBlockEventIntoEvent() (
	err error,
) {

	fn := exsrapi.PrtHdr()
	defer exsrapi.PrintExf("", fn)()

	idofblockevent, err := ExtractIDofEventGroup(BlockEvent)
	if err != nil {
		err = fmt.Errorf("SelectIDofEventGroup(BlockEvent): %w", err)
		return
	}

	for _, eid := range idofblockevent {

		//	var blocklist []exsrapi.Block
		//	blocklist, err = exsrapi.GetEventidOfBlockEvent(eid)
		var blockinflist exsrapi.BlockInfList
		blockinflist, err = exsrapi.GetEventidOfBlockEvent(eid)
		if err != nil {
			err = fmt.Errorf("exsrapi.GetEventidOfEventBox(): %w", err)
			return
		}

		//	if len(blocklist) == 0 {
		if len(blockinflist.Blockinf) == 0 {
			//	子のイベントが検出できていない。
			log.Printf("** ブロックイベントの子のイベントが検出できません。 eventid=%s\n", eid)
			continue
		}

		eventinflist := make([]exsrapi.Event_Inf, 0)
		for _, blockinf := range blockinflist.Blockinf {
			blockname := blockinf.Show_rank_label
			blocklist := blockinf.Block_list
		for _, block := range blocklist {
			var eventinf exsrapi.Event_Inf
			blockid := fmt.Sprintf("%d", block.Block_id)
			eidb := eid + "?block_id=" + blockid
			err = exsrapi.GetEventinf(eidb, &eventinf)
			if err != nil {
				log.Printf("GetEventinf(): %v", err)
				//	return fmt.Errorf("GetEventinf(): %v", status)
			} else {
				eventinf.Event_ID = eidb
				eventinf.Event_name += "[" + blockname + "][" + block.Label + "](" + blockid + ")"
				eventinflist = append(eventinflist, eventinf)
			}
		}
	}
		err = srdblib.InsertEventinflistToEvent(&eventinflist, true)
		if err != nil {
			err = fmt.Errorf("srdblib.InsertEventinflistToEvent(): %w", err)
			return
		}
		_, err = srdblib.Db.Exec("UPDATE "+srdblib.Tevent+" SET achk = ? where eventid = ?", BlockEvent%4, eid)
		log.Printf("  %s is BlockEvent. Number of Child Event is %d\n", eid, len(eventinflist))

	}

	return

}
