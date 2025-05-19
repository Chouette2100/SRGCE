package main

import (
	"log"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srapi/v2"
)

func IntegrateNewEventlistToEventtable(
	tevent string,
	eventlist []srapi.Event,
	) (
	err error,
) {

	fn := exsrapi.PrtHdr()
	defer exsrapi.PrintExf("", fn)()

	eventinflist := make([]exsrapi.Event_Inf, 0)
	for _, event := range eventlist {
		eventinflist = append(eventinflist, *exsrapi.ConvertEventToEventinf(&event))
	}

	err = StoreEventinflistInEvent(tevent, eventinflist)
	if err != nil {
		log.Printf("StoreEventinflistInEvent(): %s", err.Error())
		return
	}

	return
}
