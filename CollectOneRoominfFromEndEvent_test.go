// Copyright © 2024 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package main

import (
	"log"
	"io"
	"os"

	"net/http"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-gorp/gorp"

	"github.com/Chouette2100/srdblib"
	"github.com/Chouette2100/exsrapi"

)

func TestCollectOneRoominfFromEndEvent(t *testing.T) {
	type args struct {
		client       *http.Client
		//	tevent       string
		//	teventuser   string
		//	tuser        string
		//	tuserhistory string
		tnow         time.Time
		eid          string
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
	
		//	srdblib.Tevent = "wevent"
		//	srdblib.Teventuser = "weventuser"
		//	srdblib.Tuser = "wuser"
		//	srdblib.Tuserhistory = "wuserhistory"
	
		dial := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}
		srdblib.Dbmap = &gorp.DbMap{Db: srdblib.Db, Dialect: dial, ExpandSliceArgs: true}
		srdblib.Dbmap.AddTableWithName(srdblib.Wuser{}, "wuser").SetKeys(false, "Userno")
		srdblib.Dbmap.AddTableWithName(srdblib.Wuserhistory{}, "wuserhistory").SetKeys(false, "Userno", "Ts")
		srdblib.Dbmap.AddTableWithName(srdblib.Wevent{}, "wevent").SetKeys(false, "Eventid")
		srdblib.Dbmap.AddTableWithName(srdblib.Weventuser{}, "weventuser").SetKeys(false, "Eventid", "Userno")


		//      cookiejarがセットされたHTTPクライアントを作る
		client, jar, err := exsrapi.CreateNewClient("ShowroomCGI")
		if err != nil {
			log.Printf("CreateNewClient: %s\n", err.Error())
			return
		}
		//      すべての処理が終了したらcookiejarを保存する。
		defer jar.Save()
	
		tnow := time.Now().Truncate((time.Second)) // 秒より下を切り捨てる。時刻の比較を整数の比較と同等にする。


	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestCollectOneRoominfFromEndEvent-1",
			args: args{
				client:       client,
				//	tevent:       "tevent",
				//	teventuser:   "teventuser",
				//	tuser:        "tuser",
				//	tuserhistory: "tuserhistory",
				tnow:         tnow,
				eid:          "kareai_newad_s?block_id=20101",
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}




	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//	if err := CollectOneRoominfFromEndEvent(tt.args.client, tt.args.tevent, tt.args.teventuser, tt.args.tuser, tt.args.tuserhistory, tt.args.tnow, tt.args.eid); (err != nil) != tt.wantErr {
			if err := CollectOneRoominfFromEndEvent(tt.args.client, tt.args.tnow, tt.args.eid); (err != nil) != tt.wantErr {
				t.Errorf("CollectOneRoominfFromEndEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
