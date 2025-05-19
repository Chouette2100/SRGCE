package main

import (
	//	"database/sql"
	"fmt"
	"log"
	"strings"
	"strconv"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

func ExpandBlockEventIntoEvent(
	tevent string,
) (
	err error,
) {

	fn := exsrapi.PrtHdr()
	defer exsrapi.PrintExf("", fn)()

	idofblockevent, err := ExtractIDofEventGroup(tevent, BlockEvent)
	if err != nil {
		err = fmt.Errorf("SelectIDofEventGroup(BlockEvent): %w", err)
		return
	}

	for _, eid := range idofblockevent {

		//	var blocklist []exsrapi.Block
		//	blocklist, err = exsrapi.GetEventidOfBlockEvent(eid)
		// var blockinflist exsrapi.BlockInfList
		var blockinflist BlockInfList
		// blockinflist, err = exsrapi.GetEventidOfBlockEvent(eid)
		blockinflist, err = GetEventidOfBlockEvent(eid)
		if err != nil {
			err = fmt.Errorf("exsrapi.GetEventidOfEventBox(): %w", err)
			//	return
			continue
		}

		//	if len(blocklist) == 0 {
		if len(blockinflist.Blockinf) == 0 {
			//	子のイベントが検出できていない。
			log.Printf("** ブロックイベントの子のイベントが検出できません。 eventid=%s\n", eid)
			continue
		}

		eventinflist := make([]exsrapi.Event_Inf, 0)
		for _, blockinf := range blockinflist.Blockinf {
			var blockname_org string
			blockname := blockinf.Show_rank_label
			if strings.Contains(blockname, "\"") {
			blockname = strings.Replace(blockname, "\"", "", -1)
			} else {
				blockname_org = blockname
				blockname = "Overall"
			}
			blocklist := blockinf.Block_list
			for _, block := range blocklist {
				label := block.Label
				blockid := fmt.Sprintf("%d", block.Block_id)
				eidb := eid + "?block_id=" + blockid

				if strings.Contains(label, "\"") {
					label = strings.Replace(label, "\"", "", -1)
				} else {
					if label == blockname_org {
						label = blockname
					} else {
						label = strconv.Itoa(block.Block_id % 10)
					}
				}
				var eventinf exsrapi.Event_Inf
				// err = exsrapi.GetEventinf(eidb, &eventinf)
				var weventinf srdblib.Wevent
				var intf interface{}
				intf, err = srdblib.Dbmap.Get(&weventinf, eid)
				if err != nil {
					log.Printf("GetEventinf(): %v", err)
					//	return fmt.Errorf("GetEventinf(): %v", status)
				} else {
					weventinf = *intf.(*srdblib.Wevent)
					eventinf = exsrapi.Event_Inf{
						Event_ID:   eidb,
						I_Event_ID:   weventinf.Ieventid,
						Event_name: weventinf.Event_name + "[" + blockname + "][" + label + "](" + blockid + ")",
						Period:    weventinf.Period,
						Start_time: weventinf.Starttime,
						End_time:   weventinf.Endtime,
						Rstatus: "",
					}
					// eventinf.Event_ID = eidb
					// eventinf.Event_name += "[" + blockname + "][" + block.Label + "](" + blockid + ")"
					eventinflist = append(eventinflist, eventinf)
				}
			}
		}
		err = srdblib.InsertEventinflistToEvent(tevent, &eventinflist, true)
		if err != nil {
			err = fmt.Errorf("srdblib.InsertEventinflistToEvent(): %w", err)
			return
		}
		_, err = srdblib.Db.Exec("UPDATE "+tevent+" SET achk = ? where eventid = ?", BlockEvent%4, eid)
		log.Printf("  %s is BlockEvent. Number of Child Event is %d\n", eid, len(eventinflist))

	}

	return

}
