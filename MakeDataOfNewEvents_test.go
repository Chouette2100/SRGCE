package main

import "testing"
import (
	"io"
	"log"
	"os"
	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v2"

	"github.com/go-gorp/gorp"
)

func TestMakeDataOfNewEvents(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "TestMakeDataOfNewEvents",
			wantErr: false,
		},
		// TODO: Add test cases.
	}

	logfile, err := exsrapi.CreateLogfile("TestMakeDataOfNewEvents.txt")
	if err != nil {
		log.Printf("err=%s.\n", err.Error())
		return
	}
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))

	//      データベースとの接続をオープンする。
	dbconfig, err := srdblib.OpenDb("DBConfig.yml")
	if err != nil {
		log.Printf("Database error. err=%s.\n", err.Error())
		return
	}
	if dbconfig.UseSSH {
		defer srdblib.Dialer.Close()
	}
	defer srdblib.Db.Close()
	log.Printf("dbconfig=%+v.\n", dbconfig)

	//	srdblib.Tevent = "wevent"
	//	srdblib.Teventuser = "weventuser"
	//	srdblib.Tuser = "wuser"
	//	srdblib.Tuserhistory = "wuserhistory"

	dial := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}
	srdblib.Dbmap = &gorp.DbMap{Db: srdblib.Db, Dialect: dial, ExpandSliceArgs: true}

	srdblib.Dbmap.AddTableWithName(srdblib.Event{}, "event").SetKeys(false, "Eventid")
	srdblib.Dbmap.AddTableWithName(srdblib.Wevent{}, "wevent").SetKeys(false, "Eventid")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MakeDataOfNewEvents(); (err != nil) != tt.wantErr {
				t.Errorf("MakeDataOfNewEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
