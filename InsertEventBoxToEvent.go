package main

import (
	"fmt"
	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
	"log"
)

func InsertEventBoxToWevent() (
	err error,
) {

	idofeventbox, err := GetIDofEventbox()
	if err != nil {
		err = fmt.Errorf("GetIDofEventbox(): %w", err)
		return
	}

	eventinflist := make([]exsrapi.Event_Inf, 0)
	for _, eid := range idofeventbox {

		namelist, err := exsrapi.GetEventidOfEventBox(eid)
		if err != nil {
			err = fmt.Errorf("exsrapi.GetEventidOfEventBox(): %w", err)
			return err
		}

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
	}

	isntins, err := srdblib.InsertEventinflistToEvent(eventinflist)
	if err != nil {
		err = fmt.Errorf("srdblib.InsertEventinflistToEvent(): %w", err)
		return
	}
	log.Printf("InsertEventBoxToWevent(): %v", isntins)
	return

}
