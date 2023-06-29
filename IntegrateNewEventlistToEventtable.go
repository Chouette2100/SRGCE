package main

import (
	"log"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srapi"
)

func IntegrateNewEventlistToEventtable(eventlist []srapi.Event) (
	err error,
) {

	fn := exsrapi.PrtHdr()
	defer exsrapi.PrintExf("", fn)()

	eventinflist := make([]exsrapi.Event_Inf, 0)
	for _, event := range eventlist {
		eventinflist = append(eventinflist, *exsrapi.ConvertEventToEventinf(&event))
	}

	err = StoreEventinflistInEvent(eventinflist)
	if err != nil {
		log.Printf("InsertEventinflistToEvent(): %s", err.Error())
		return
	}

	return
}
