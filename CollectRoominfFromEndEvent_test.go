package main

import (
	"io"
	"log"
	"os"

	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-gorp/gorp"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srdblib/v3"
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
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))

	//      cookiejarがセットされたHTTPクライアントを作る
	client, jar, err := exsrapi.CreateNewClient("ShowroomCGI")
	if err != nil {
		log.Printf("CreateNewClient: %s\n", err.Error())
		return
	}
	//      すべての処理が終了したらcookiejarを保存する。
	defer jar.Save()

	//      データベースとの接続をオープンする。
	var dbconfig *srdblib.DBConfig
	Db, dbconfig, err = srdblib.OpenDb("DBConfig.enc.yml")
	if err != nil {
		log.Printf("srdblib.OpenDb() error. err=%s.\n", err.Error())
		return
	}
	if dbconfig.UseSSH {
		defer srdblib.Dialer.Close()
	}
	defer Db.Close()

	log.Printf("dbconfig=%v\n", dbconfig)

	//	srdblib.Tevent = "wevent"
	//	srdblib.Teventuser = "weventuser"
	//	srdblib.Tuser = "wuser"
	//	srdblib.Tuserhistory = "wuserhistory"

	dial := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}
	Dbmap = &gorp.DbMap{Db: Db, Dialect: dial, ExpandSliceArgs: true}
	Dbmap.AddTableWithName(srdblib.Wuser{}, "wuser").SetKeys(false, "Userno")
	Dbmap.AddTableWithName(srdblib.Userhistory{}, "wuserhistory").SetKeys(false, "Userno", "Ts")
	Dbmap.AddTableWithName(srdblib.Event{}, "wevent").SetKeys(false, "Eventid")
	Dbmap.AddTableWithName(srdblib.Eventuser{}, "weventuser").SetKeys(false, "Eventid", "Userno")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CollectRoominfFromEndEvent(client, "wevent", "weventuser", "wuser", "wuserhistory"); (err != nil) != tt.wantErr {
				t.Errorf("GetRoominfAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
