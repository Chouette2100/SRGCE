package main

import (
	//	"database/sql"
	"fmt"
	"log"
	// "strconv"
	// "strings"

	"net/http"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srapi/v2"
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
		/*
			blockinflist, err = GetEventidOfBlockEvent(eid)
			if err != nil {
				err = fmt.Errorf("exsrapi.GetEventidOfEventBox(): %w", err)
				//	return
				continue
			}
		*/
		client := &http.Client{}
		var er *srapi.EventRanking
		er, err = srapi.ApiEventRanking(client, eid, 1)
		log.Printf("**Evnetid = %s\n", eid)
		for _, ebl := range er.EventBlockList {
			log.Printf("**  ShwRankLabel = %s\n", ebl.ShowRankLabel)
			var blockinf BlockInf
			blockinf.Show_rank_label = ebl.ShowRankLabel
			blockinf.Block_list = make([]Block, 0)
			for _, eb := range ebl.BlockList {
				log.Printf("**    BlockID = %d, Label = %s\n", eb.BlockID, eb.Label)
				var block Block
				block.Label = eb.Label
				block.Block_id = eb.BlockID
				blockinf.Block_list = append(blockinf.Block_list, block)
			}
			blockinflist.Blockinf = append(blockinflist.Blockinf, blockinf)
		}

		//	if len(blocklist) == 0 {
		if len(blockinflist.Blockinf) == 0 {
			//	子のイベントが検出できていない。
			log.Printf("** ブロックイベントの子のイベントが検出できません。 eventid=%s\n", eid)
			continue
		}

		eventinflist := make([]exsrapi.Event_Inf, 0)
		for _, blockinf := range blockinflist.Blockinf {
			// var blockname_org string
			blockname := blockinf.Show_rank_label
			/*
				if strings.Contains(blockname, "\"") {
					blockname = strings.Replace(blockname, "\"", "", -1)
				} else {
					blockname_org = blockname
					blockname = "Overall"
				}
			*/
			blocklist := blockinf.Block_list
			for _, block := range blocklist {
				label := block.Label
				blockid := fmt.Sprintf("%d", block.Block_id)
				eidb := eid + "?block_id=" + blockid

				/*
					if strings.Contains(label, "\"") {
						label = strings.Replace(label, "\"", "", -1)
					} else {
						if label == blockname_org {
							label = blockname
						} else {
							label = strconv.Itoa(block.Block_id % 10)
						}
					}
				*/
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
						I_Event_ID: weventinf.Ieventid,
						Event_name: weventinf.Event_name + "[" + blockname + "][" + label + "](" + blockid + ")",
						Period:     weventinf.Period,
						Start_time: weventinf.Starttime,
						End_time:   weventinf.Endtime,
						Rstatus:    "",
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
		for _, eventinf := range eventinflist {
			if !eventinf.Valid {
				log.Printf("  **Inserted[%s]: %s\n", eventinf.Event_ID, eventinf.Event_name)
			} else {
				var intf interface{}
				wev := srdblib.Wevent{}
				if tevent != "wevent" {
					err = fmt.Errorf("ExpandBlockEventIntoEvent(): tevent != wevent")
				} else {
					intf, err = srdblib.Dbmap.Get(&wev, eventinf.Event_ID)
				}
				if err != nil || intf == nil || intf.(*srdblib.Wevent).Event_name == eventinf.Event_name {
					log.Printf("  **Ignored[%s]: %s\n", eventinf.Event_ID, eventinf.Event_name)
				} else {
					wev = *intf.(*srdblib.Wevent)
					wev.Event_name = eventinf.Event_name
					_, err = srdblib.Dbmap.Update(&wev)
					if err != nil {
						log.Printf("Update(): %v", err)
					}
					log.Printf("  **Updated[%s]: %s\n", eventinf.Event_ID, eventinf.Event_name)
				}
			}
		}
		_, err = srdblib.Db.Exec("UPDATE "+tevent+" SET achk = ? where eventid = ?", BlockEvent%4, eid)
		log.Printf("  %s is BlockEvent. Number of Child Event is %d\n", eid, len(eventinflist))

	}

	return

}
