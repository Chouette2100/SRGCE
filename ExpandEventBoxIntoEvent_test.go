package main

import(
	"log"
	"io"
	"os"

"testing"

"github.com/Chouette2100/exsrapi"
"github.com/Chouette2100/srdblib"
)

func TestExpandEventBoxIntoEvent(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "TestInsertEventBoxToWevent",
			wantErr:		false,
		},
		// TODO: Add test cases.
	}
	
	logfile, err := exsrapi.CreateLogfile("TestGetAndInsertEventBox")
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


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExpandEventBoxIntoEvent("wevent"); (err != nil) != tt.wantErr {
				t.Errorf("InsertEventBoxToWevent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
