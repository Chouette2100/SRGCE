package main

import (
	"io"
	"log"
	"os"

	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
)

func TestGetIDofEventbox(t *testing.T) {
	tests := []struct {
		name             string
		wantIdofeventbox []string
		wantErr          bool
	}{
		{
			name: "GetIDofEventbox",
			wantIdofeventbox: []string{},
			wantErr:          false,
		},
		// TODO: Add test cases.
	}

	logfile, err := exsrapi.CreateLogfile("TestGetIDofEventbox")
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

	srdblib.Tevent = "wevent"
	srdblib.Teventuser = "weventuser"
	srdblib.Tuser = "wuser"
	srdblib.Tuserhistory = "wuserhistory"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIdofeventbox, err := GetIDofEventbox()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAndInsertEventBox() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIdofeventbox, tt.wantIdofeventbox) {
				t.Errorf("GetAndInsertEventBox() = %v, want %v", gotIdofeventbox, tt.wantIdofeventbox)
			}
		})
	}
}
