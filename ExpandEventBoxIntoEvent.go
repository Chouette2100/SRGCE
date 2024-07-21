package main

import (
	"fmt"
	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
	"log"
)

func ExpandEventBoxIntoEvent(
	tevent string,
) (
	err error,
) {

	fn := exsrapi.PrtHdr()
	defer exsrapi.PrintExf("", fn)()

	idofeventbox, err := ExtractIDofEventGroup(tevent, EventBox)
	if err != nil {
		err = fmt.Errorf("GetIDofEventbox(): %w", err)
		return
	}

	for _, eid := range idofeventbox {

		var namelist []string
		namelist, err = exsrapi.GetEventidOfEventBox(eid)
		if err != nil {
			err = fmt.Errorf("exsrapi.GetEventidOfEventBox(): %w", err)
			return
		}

		if len(namelist) == 0 {
			//      子のイベントが検出できていない。
			log.Printf("** イベントボックスの子のイベントが検出できません。 eventid=%s\n", eid)
			continue
		}

		eventinflist := make([]exsrapi.Event_Inf, 0)
		for _, name := range namelist {
			var eventinf exsrapi.Event_Inf
			err = exsrapi.GetEventinf(name, &eventinf)
			if err != nil {
				log.Printf("GetEventinf(): %v", err)
				//	return fmt.Errorf("GetEventinf(): %v", status)
			} else {
				eventinflist = append(eventinflist, eventinf)
			}
		}

		err = srdblib.InsertEventinflistToEvent(tevent, &eventinflist, true)
		if err != nil {
			err = fmt.Errorf("srdblib.InsertEventinflistToEvent(): %w", err)
			return
		}
		_, err = srdblib.Db.Exec("UPDATE " + tevent + " SET achk = ? where eventid = ?", EventBox % 4, eid)
		log.Printf("  %s is Event Box. Number of Child Event is %d\n", eid, len(eventinflist))

	}

	//	log.Printf("InsertEventBoxToWevent(): %v", isntins)
	return

}
