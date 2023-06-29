package main

import (
	"os"
	"io"
	"log"

	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/srdblib"
	"github.com/Chouette2100/exsrapi"
)

func TestCollectRoominfFromEndEvent(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test1",
			wantErr: false,
		},

		// TODO: Add test cases.
	}
	logfile, err := exsrapi.CreateLogfile("TestGetRoominfAll", "log")
	if err != nil {
		log.Printf("exsrapi.CreateLogfile() error. err=%s.\n", err.Error())
		return
	}
	defer logfile.Close()
	log.SetOutput(io.MultiWriter(logfile,os.Stdout))


	//      データベースとの接続をオープンする。
	dbconfig, err := srdblib.OpenDb("DBConfig.yml")
	if err != nil {
		log.Printf("srdblib.OpenDb() error. err=%s.\n", err.Error())
		return
	}
	if dbconfig.UseSSH {
		defer srdblib.Dialer.Close()
	}
	defer srdblib.Db.Close()

	log.Printf("dbconfig=%v\n",	dbconfig)

	srdblib.Tevent = "wevent"
	srdblib.Teventuser = "weventuser"
	srdblib.Tuser = "wuser"
	srdblib.Tuserhistory = "wuserhistory"


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CollectRoominfFromEndEvent(); (err != nil) != tt.wantErr {
				t.Errorf("GetRoominfAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
