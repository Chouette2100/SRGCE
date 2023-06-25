package main

import (
	"fmt"
	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
	"log"
)

func InsertBlockeventToEvent() (
	err error,
) {

	idofblockevent, err := SelectIDofEventGroup(BlockEvent)
	if err != nil {
		err = fmt.Errorf("SelectIDofEventGroup(BlockEvent): %w", err)
		return
	}

	eventinflist := make([]exsrapi.Event_Inf, 0)
	for _, eid := range idofblockevent {

		blocklist, err := exsrapi.GetEventidOfBlockEvent(eid)
		if err != nil {
			err = fmt.Errorf("exsrapi.GetEventidOfEventBox(): %w", err)
			return err
		}

		for _, block := range blocklist {
			var eventinf exsrapi.Event_Inf
			blockid := fmt.Sprintf("%d",block.Block_id)
			eidb := eid + "?block_id=" + blockid
			err = exsrapi.GetEventinf(eidb, &eventinf)
			if err != nil {
				log.Printf("GetEventinf(): %v", err)
				//	return fmt.Errorf("GetEventinf(): %v", status)
			} else {
				eventinf.Event_ID = eidb
				eventinf.Event_name += "[" + block.Label + "](" + blockid + ")"
				eventinflist = append(eventinflist, eventinf)
			}
		}
	}

	isntins, err := srdblib.InsertEventinflistToEvent(eventinflist)
	if err != nil {
		err = fmt.Errorf("srdblib.InsertEventinflistToEvent(): %w", err)
		return
	}
	log.Printf("InserBlockevneToEvent(): %v", isntins)
	return

}
