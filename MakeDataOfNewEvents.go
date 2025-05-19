package main

import (
	"fmt"
	"log"

	//	"database/sql"
	//	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/copier"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

func MakeDataOfNewEvents() (err error) {

	fn := exsrapi.PrtHdr()
	defer exsrapi.PrintExf("", fn)()


	var thdata *exsrapi.Thdata
	thdata, err = exsrapi.ReadThdata()
	if err != nil {
		err = fmt.Errorf("ReadThdata() error: %w", err)
		return
	}

	sqlst := "select * from wevent "
	sqlst += " where achk = 0 and  now() between SUBDATE(starttime,INTERVAL ? hour) and endtime  "
	sqlst += "   and eventid not in "
	sqlst += "  ( select eventid from event "
	sqlst += " where Achk = 0  and now() between SUBDATE(starttime,INTERVAL ? hour) and endtime ) "
	sqlst += "  order by starttime "

	var rows []interface{}
	rows, err = srdblib.Dbmap.Select(srdblib.Wevent{}, sqlst, thdata.Hh, thdata.Hh)

	//	srdblib.Dbmap.AddTableWithName(srdblib.Event{}, "event").SetKeys(false, "Eventid")
	for _, v := range rows {
		wevent := v.(*srdblib.Wevent)
		//	log.Printf("%24s%s\n", event.Eventid, event.Event_name)
		// event := srdblib.Event(*wevent)
		var event srdblib.Event
		copier.Copy(&event, wevent)
		err = MakeDataOfEvent(&event, thdata)
		if err != nil {
			err = fmt.Errorf("MakeDataOfEvent() error: %w", err)
			return
		}
		//	if i == 0 {
		//		break
		//	}
	}

	return
}

func MakeDataOfEvent(event *srdblib.Event, thdata *exsrapi.Thdata) (err error) {

	log.Printf(" MakeDataOfEvent() eventid=%s\n", event.Eventid)

	event.Intervalmin = 5
	event.Modmin, event.Modsec = exsrapi.MakeSampleTime(240, 40)

	event.Fromorder = thdata.From
	event.Toorder = thdata.To

	event.Resethh = 4
	event.Resetmm = 0
	//	event.Nobasis =
	//	event.Target =
	event.Maxdsp = 25
	event.Cmap = 2
	//	event.Target =
	event.Rstatus = ""
	//	event.Maxpoint =
	var eventinf exsrapi.Event_Inf = exsrapi.Event_Inf{
		Event_ID: event.Eventid,
		Event_name: event.Event_name,
	}
	err = exsrapi.SetThdata(&eventinf, thdata)
	if err != nil {
		err = fmt.Errorf("SetThdata() error: %w", err)
		return
	}
	event.Thinit = eventinf.Thinit
	event.Thdelta = eventinf.Thdelta
	//	event.achk =
	//	event.Aclr =

	log.Printf("Thinit=%d, Thdelta=%d\n", event.Thinit, event.Thdelta)

	err = srdblib.Dbmap.Insert(event)
	if err != nil {
		err = fmt.Errorf("Dbmap.Insert() error: %w", err)
	}

	return
}