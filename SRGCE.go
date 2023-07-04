/*!
Copyright © 2022 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php

*/

/*
	ソースのダウンロード、ビルドについて以下簡単に説明します。詳細は以下の記事を参考にしてください。
	WindowsもLinuxも特記した部分以外は同じです。

		【Windows】かんたんなWebサーバーの作り方
			https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/c5cab5

		---------------------

		【Windows】Githubにあるサンプルプログラムの実行方法
			https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/e27fc9

		【Unix/Linux】Githubにあるサンプルプログラムの実行方法
			https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/220e38

			ロードモジュールさえできればいいということでしたらコマンド一つでできます。

【Unix/Linux】

	$ mysql -p -u .....
	Enter password:
	mysql> show databases;
	CREATE DATABASE IF NOT EXISTS `showroom` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
	mysql> use `showroom`;
	Database changed
	mysql> source SQL/wevent.sql
	mysql> source SQL/weventuser.sql
	mysql> source SQL/user.sql
	mysql> source SQL/userhistory.sql
	mysql> quit

	$ cd ~/go/src
	$ curl -OL https://github.com/Chouette2100/t008srapi/archive/refs/tags/vl.m.n.tar.gz
	$ tar xvf vl.m.n.tar.gz
	$ mv SRGCE-l.m.n SRGCE
	$ cd SRGCE
	$ go mod init
	$ go mod tidy
	$ go build SRGCE.go
	$ vi DBConfig.yml
	$ ./SRGCE

【Windows】

Microsoft Windows [Version 10.0.22000.856]
(c) Microsoft Corporation. All rights reserved.

C:\Users\chouette>cd go

C:\Users\chouette\go>cd src

作業はかならず %HOMEPATH%\go\src の下で行います。

以下、要するに https://github.com/Chouette2100/t008srapi/releases にあるv0.2.0のZIPファイルSource code (zip) からソースをとりだしてくださいということなので、ブラウザでダウンロードしてエクスプローラで解凍というこでもけっこうです。なんならこの記事の最後にあるgithubのソースをエディターにコピペで作るということでもかまいません（この場合文字コードはかならずUTF-8にしてください 改行はLFになっています。というようなことを考えるとやっぱりダウンロードして解凍が安全かも）

C:\Users\chouette\go\src>mkdir t008srapi

C:\Users\chouette\go\src>cd t008srapi

C:\Users\chouette\go\src\t008srapi>curl -OL https://github.com/Chouette2100/t008srapi/archive/refs/tags/v0.2.0.zip

	  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
	                                 Dload  Upload 	client, cookiejar, err := exsrapi.CreateNewClient("")
		if err != nil {
			log.Printf("exsrapi.CeateNewClient(): %s", err.Error())
			return //	エラーがあれば、ここで終了
		}
		defer cookiejar.Save()

		//	テンプレートで使用する関数を定義する
		funcMap := template.FuncMap{
			"Comma":         func(i int) string { return humanize.Comma(int64(i)) },                       //	3桁ごとに","を入れる関数。
			"UnixTimeToStr": func(i int64) string { return time.Unix(int64(i), 0).Format("01-02 15:04") }, //	UnixTimeを年月日時分に変換する関数。
		}

		// テンプレートをパースする
		tpl := template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/t008top.gtpl"))

		// テンプレートに埋め込むデータ（ポイントやランク）を作成する
		top := new(T008top)
		top.TimeNow = time.Now().Unix()

		top.Eventlist, err = srapi.MakeEventListByApi(client)
		if err != nil {
			err = fmt.Errorf("MakeListOfPoints(): %w", err)
			log.Printf("MakeListOfPoints() returned error %s\n", err.Error())
			top.ErrMsg = err.Error()
		}
		top.Totalcount = len(top.Eventlist)

		//	ソートが必要ないときは次の行とimportの"sort"をコメントアウトする。
		//	無名関数のリターン値でソート条件を変更できます。
		sort.Slice(top.Eventlist, func(i, j int) bool { return top.Eventlist[i].Ended_at > top.Eventlist[j].Ended_at })

	  Total   Spent    Left  Speed
	  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0

100  6265    0  6265    0     0   6777      0 --:--:-- --:--:-- --:--:-- 16400

C:\Users\chouette\go\src\t008srapi>call powershell -command "Expand-Archive v0.2.0.zip"

C:\Users\chouette\go\src\t008srapi>tree
フォルダー パスの一覧
ボリューム シリアル番号は E2CD-BDF1 です
C:.
└─v0.2.0

	└─t008srapi-0.2.0
	    ├─public
	    └─templates

C:\Users\chouette\go\src\t008srapi>xcopy /e v0.2.0\t008srapi-0.2.0\*.*
v0.1.0\t007srapi-0.1.0\freebsd.bat
v0.1.0\t007srapi-0.1.0\freebsd.sh
v0.1.0\t007srapi-0.1.0\LICENSE
v0.1.0\t007srapi-0.1.0\README.md
v0.1.0\t007srapi-0.1.0\t007srapi.go
v0.1.0\t007srapi-0.1.0\public\index.html
v0.1.0\t007srapi-0.1.0\templates\top.gtpl
7 File(s) copied

C:\Users\chouette\go\src\t008srapi>rmdir /s /q v0.2.0

C:\Users\chouette\go\src\t008srapi>del v0.2.0.zip

ここで次のような構成になっていればOKです。top.gtpl と index.html が所定の場所にあることをかならず確かめてください。

C:%HOMEPATH%\go\src\t008srapi --+-- t008srapi.go

	|
	+-- \templates --- t008top.gtpl
	|
	+-- \public    --- index.html

ここからはコマンド三つでビルドが完了します。

C:\Users\chouette\go\src\t008srapi>go mod init
go: creating new go.mod: module t008srapi
go: to add module requirements and sums:

	go mod tidy

C:\Users\chouette\go\src\t008srapi>go mod tidy
go: finding module for package github.com/dustin/go-humanize
go: downloading github.com/dustin/go-humanize v1.0.0
go: found github.com/dustin/go-humanize in github.com/dustin/go-humanize v1.0.0

C:\Users\chouette\go\src\t008srapi>go build t008srapi.go

あとは

C:\Users\chouette\go\src\t008srapi>t008srapi

でWebサーバが起動します。ここでセキュリティー上の警告が出ると思いますが、説明をよく読んで問題ないと思ったらアクセスを許可してください（もちろん許可しなければWebサーバは使えなくなります）

# Webサーバを起動したままにしておいてブラウザを開き

http://localhost:8080/t008top

で、実行時点でのイベントの一覧が表示されます。

	Ver. 01AA00	基本的動作を確認する。
	Ver. 01AB00	進捗状況の表示をuserの更新基準からeventuserの更新基準に変更する。
	Ver. 01AB01	取得したデータがどのように処理されたかを表示する。
	Ver. 01AB02	最終結果取得済みのイベントをスキップする。
	Ver. 01AC00	イベント終了日範囲の取得方法を変更する（tnow.Truncate(24 * time.Hour)が9時以前は1日前の日になるから）
*/
package main

import (
	//	"html/template"
	"io" //　ログ出力設定用。必要に応じて。
	"log"
	//	"net/http"
	//	"net/http/cgi"
	"os"

	//	"github.com/dustin/go-humanize"

	//	"database/sql"
	//	_ "github.com/go-sql-driver/mysql"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srdblib"
)

const Version = "01AC00"

func main() {

	//      ログファイルを開く。
	logfile, err := exsrapi.CreateLogfile(Version)
	if err != nil {
		log.Printf("err=%s.\n", err.Error())
		return
	}
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))

	fn := exsrapi.PrtHdr()
	defer exsrapi.PrintExf("", fn)()

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

	//      テーブルは"w"で始まるものを操作の対象とする。
	srdblib.Tevent = "wevent"
	srdblib.Teventuser = "weventuser"
	srdblib.Tuser = "wuser"
	srdblib.Tuserhistory = "wuserhistory"

	//      現在開催中のイベントの一覧を求める。
	cel, err := CreateCurrentEventList()
	if err != nil {
		log.Printf("GetCurrentEventList(): %s", err.Error())
		return
	}

	//	取得したイベント情報をデータベースに格納する。
	err = IntegrateNewEventlistToEventtable(cel.Eventlist)
	if err != nil {
		log.Printf("InsertIntoEvent(): %s", err.Error())
		return
	}

	//	ブロックイベントを展開する。
	err = ExpandBlockEventIntoEvent()
	if err != nil {
		log.Printf("InsertBlockeventToEvent(): %s", err.Error())
	}

	//	イベントボックスを展開する。
	err = ExpandEventBoxIntoEvent()
	if err != nil {
		log.Printf("InsertEventBoxToWevent(): %s", err.Error())
	}

	//	結果が発表されたイベントの順位と獲得ポイントを取得する
	//	これはイベント終了日の翌日12時から翌々日12時までのあいだに行う必要がある)
	//	前日終了のイベントのデータを取得するか、前々日のものを取得するかは実行時刻に応じて判断される。
	CollectRoominfFromEndEvent()

}
