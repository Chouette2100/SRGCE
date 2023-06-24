package main
import (
	"log"
	"time"
	"strings"
	"bytes"
	"strconv"
	"fmt"
	"os"
	

	//	"database/sql"
	//	_ "github.com/go-sql-driver/mysql"

	"net/http"

	"github.com/PuerkitoBio/goquery"

	//	"github.com/Chouette2100/srdblib"
	"github.com/Chouette2100/srapi"
	"github.com/Chouette2100/exsrapi"
)
type RoomInfo struct {
	Name      string //     ルーム名のリスト
	Longname  string
	Shortname string
	Account   string //     アカウントのリスト、アカウントは配信のURLの最後の部分の英数字です。
	ID        string //     IDのリスト、IDはプロフィールのURLの最後の部分で5～6桁の数字です。
	Userno    int
	//      APIで取得できるデータ(1)
	Genre      string
	Rank       string
	Irank      int
	Nrank      string
	Prank      string
	Followers  int
	Sfollowers string
	Fans       int
	Fans_lst   int
	Level      int
	Slevel     string
	//      APIで取得できるデータ(2)
	Order        int
	Point        int //     イベント終了後12時間〜36時間はイベントページから取得できることもある
	Spoint       string
	Istarget     string
	Graph        string
	Iscntrbpoint string
	Color        string
	Colorvalue   string
	//	Colorinflist ColorInfList
	Formid       string
	Eventid      string
	Status       string
	Statuscolor  string
}

type RoomInfoList []RoomInfo

// sort.Sort()のための関数三つ
func (r RoomInfoList) Len() int {
	return len(r)
}

func (r RoomInfoList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RoomInfoList) Choose(from, to int) (s RoomInfoList) {
	s = r[from:to]
	return
}

var SortByFollowers bool

// 降順に並べる
func (r RoomInfoList) Less(i, j int) bool {
	//      return e[i].point < e[j].point
	if SortByFollowers {
			return r[i].Followers > r[j].Followers
	} else {
			return r[i].Point > r[j].Point
	}
}


func GetEventInfAndRoomList(
	eventid string,
	breg int,
	ereg int,
	eventinfo *exsrapi.Event_Inf,
	roominfolist *RoomInfoList,
) (
	status int,
) {

	//	画面からのデータ取得部分は次を参考にしました。
	//		はじめてのGo言語：Golangでスクレイピングをしてみた
	//		https://qiita.com/ryo_naka/items/a08d70f003fac7fb0808

	//	_url := "https://www.showroom-live.com/event/" + EventID
	//	_url = "file:///C:/Users/kohei47/Go/src/EventRoomList03/20210128-1143.html"
	//	_url = "file:20210128-1143.html"

	var doc *goquery.Document
	var err error

	inputmode := "url"
	eventidorfilename := eventid
	maxroom := ereg

	status = 0

	if inputmode == "file" {

		//	ファイルからドキュメントを作成します
		f, e := os.Open(eventidorfilename)
		if e != nil {
			//	log.Fatal(e)
			log.Printf("err=[%s]\n", e.Error())
			status = -1
			return
		}
		defer f.Close()
		doc, err = goquery.NewDocumentFromReader(f)
		if err != nil {
			//	log.Fatal(err)
			log.Printf("err=[%s]\n", err.Error())
			status = -1
			return
		}

		content, _ := doc.Find("head > meta:nth-child(6)").Attr("content")
		content_div := strings.Split(content, "/")
		(*eventinfo).Event_ID = content_div[len(content_div)-1]

	} else {
		//	URLからドキュメントを作成します
		_url := "https://www.showroom-live.com/event/" + eventidorfilename
		/*
			doc, err = goquery.NewDocument(_url)
		*/
		resp, error := http.Get(_url)
		if error != nil {
			log.Printf("GetEventInfAndRoomList() http.Get() err=%s\n", error.Error())
			status = 1
			return
		}
		defer resp.Body.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)

		//	bufstr := buf.String()
		//	log.Printf("%s\n", bufstr)

		//	doc, error = goquery.NewDocumentFromReader(resp.Body)
		doc, error = goquery.NewDocumentFromReader(buf)
		if error != nil {
			log.Printf("GetEventInfAndRoomList() goquery.NewDocumentFromReader() err=<%s>.\n", error.Error())
			status = 1
			return
		}

		(*eventinfo).Event_ID = eventidorfilename
	}
	//	log.Printf(" eventid=%s\n", (*eventinfo).Event_ID)

	cevent_id, exists := doc.Find("#eventDetail").Attr("data-event-id")
	if !exists {
		log.Printf("data-event-id not found. Event_ID=%s\n", (*eventinfo).Event_ID)
		status = -1
		return
	}
	eventinfo.I_Event_ID, _ = strconv.Atoi(cevent_id)

	selector := doc.Find(".detail")
	(*eventinfo).Event_name = selector.Find(".tx-title").Text()
	if (*eventinfo).Event_name == "" {
		log.Printf("Event not found. Event_ID=%s\n", (*eventinfo).Event_ID)
		status = -1
		return
	}
	(*eventinfo).Period = selector.Find(".info").Text()
	eventinfo.Period = strings.Replace(eventinfo.Period, "\u202f", " ", -1)
	period := strings.Split((*eventinfo).Period, " - ")
	if inputmode == "url" {
		(*eventinfo).Start_time, _ = time.Parse("Jan 2, 2006 3:04 PM MST", period[0]+" JST")
		(*eventinfo).End_time, _ = time.Parse("Jan 2, 2006 3:04 PM MST", period[1]+" JST")
	} else {
		(*eventinfo).Start_time, _ = time.Parse("2006/01/02 15:04 MST", period[0]+" JST")
		(*eventinfo).End_time, _ = time.Parse("2006/01/02 15:04 MST", period[1]+" JST")
	}

	(*eventinfo).EventStatus = "BeingHeld"
	if (*eventinfo).Start_time.After(time.Now()) {
		(*eventinfo).EventStatus = "NotHeldYet"
	} else if (*eventinfo).End_time.Before(time.Now()) {
		(*eventinfo).EventStatus = "Over"
	}

	//	イベントに参加しているルームの数を求めます。
	//	参加ルーム数と表示されているルームの数は違うので、ここで取得したルームの数を以下の処理で使うわけではありません。
	SNoEntry := doc.Find("p.ta-r").Text()
	fmt.Sscanf(SNoEntry, "%d", &((*eventinfo).NoEntry))
	log.Printf("[%s]\n[%s] [%s] (*event).EventStatus=%s NoEntry=%d\n",
		(*eventinfo).Event_name,
		(*eventinfo).Start_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).End_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).EventStatus, (*eventinfo).NoEntry)
	log.Printf("breg=%d ereg=%d\n", breg, ereg)

	//	eventno, _, _ := SelectEventNoAndName(eventidorfilename)
	//	log.Printf(" eventno=%d\n", eventno)
	//	(*eventinfo).Event_no = eventno

	//	抽出したルームすべてに対して処理を繰り返す(が、イベント開始後の場合の処理はルーム数がbreg、eregの範囲に限定）
	//	イベント開始前のときはすべて取得し、ソートしたあてで範囲を限定する）
	doc.Find(".listcardinfo").EachWithBreak(func(i int, s *goquery.Selection) bool {
		//	log.Printf("i=%d\n", i)
		if (*eventinfo).Start_time.Before(time.Now()) {
			if i < breg-1 {
				return true
			}
			if i == maxroom {
				return false
			}
		}

		var roominfo RoomInfo

		roominfo.Name = s.Find(".listcardinfo-main-text").Text()

		spoint1 := strings.Split(s.Find(".listcardinfo-sub-single-right-text").Text(), ": ")

		var point int64
		if spoint1[0] != "" {
			spoint2 := strings.Split(spoint1[1], "pt")
			fmt.Sscanf(spoint2[0], "%d", &point)

		} else {
			point = -1
		}
		roominfo.Point = int(point)

		ReplaceString := ""

		selection_c := s.Find(".listcardinfo-menu")

		account, _ := selection_c.Find(".room-url").Attr("href")
		if inputmode == "file" {
			ReplaceString = "https://www.showroom-live.com/"
		} else {
			ReplaceString = "/r/"
		}
		roominfo.Account = strings.Replace(account, ReplaceString, "", -1)
		roominfo.Account = strings.Replace(roominfo.Account, ReplaceString, "/", -1)

		roominfo.ID, _ = selection_c.Find(".js-follow-btn").Attr("data-room-id")
		roominfo.Userno, _ = strconv.Atoi(roominfo.ID)

		roominfo.Irank = i + 1

		*roominfolist = append(*roominfolist, roominfo)

		//	log.Printf("%11s %-20s %-10s %s\n",
		//		humanize.Comma(int64(roominfo.Point)), roominfo.Account, roominfo.ID, roominfo.Name)
		return true

	})

	(*eventinfo).NoRoom = len(*roominfolist)

	log.Printf(" GetEventInfAndRoomList() len(*roominfolist)=%d\n", len(*roominfolist))

	return
}

func GetEventInfAndRoomListBR(
	client *http.Client,
	eventid string,
	breg int,
	ereg int,
	eventinfo *exsrapi.Event_Inf,
	roominfolist *RoomInfoList,
) (
	status int,
) {

	status = 0

	var doc *goquery.Document
	var err error

	inputmode := "url"
	eventidorfilename := eventid

	status = 0

	//	URLからドキュメントを作成します
	_url := "https://www.showroom-live.com/event/" + eventidorfilename
	/*
		doc, err = goquery.NewDocument(_url)
	*/
	resp, error := http.Get(_url)
	if error != nil {
		log.Printf("GetEventInfAndRoomList() http.Get() err=%s\n", error.Error())
		status = 1
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bufstr := buf.String()

	log.Printf("%s\n", bufstr)

	//	doc, error = goquery.NewDocumentFromReader(resp.Body)
	doc, error = goquery.NewDocumentFromReader(buf)
	if error != nil {
		log.Printf("GetEventInfAndRoomList() goquery.NewDocumentFromReader() err=<%s>.\n", error.Error())
		status = 1
		return
	}

	(*eventinfo).Event_ID = eventidorfilename
	//	log.Printf(" eventid=%s\n", (*eventinfo).Event_ID)

	cevent_id, exists := doc.Find("#eventDetail").Attr("data-event-id")
	if !exists {
		log.Printf("data-event-id not found. Event_ID=%s\n", (*eventinfo).Event_ID)
		status = -1
		return
	}
	eventinfo.I_Event_ID, _ = strconv.Atoi(cevent_id)
	event_id := eventinfo.I_Event_ID

	selector := doc.Find(".detail")
	(*eventinfo).Event_name = selector.Find(".tx-title").Text()
	if (*eventinfo).Event_name == "" {
		log.Printf("Event not found. Event_ID=%s\n", (*eventinfo).Event_ID)
		status = -1
		return
	}
	(*eventinfo).Period = selector.Find(".info").Text()
	eventinfo.Period = strings.Replace(eventinfo.Period, "\u202f", " ", -1)
	period := strings.Split((*eventinfo).Period, " - ")
	if inputmode == "url" {
		(*eventinfo).Start_time, _ = time.Parse("Jan 2, 2006 3:04 PM MST", period[0]+" JST")
		(*eventinfo).End_time, _ = time.Parse("Jan 2, 2006 3:04 PM MST", period[1]+" JST")
	} else {
		(*eventinfo).Start_time, _ = time.Parse("2006/01/02 15:04 MST", period[0]+" JST")
		(*eventinfo).End_time, _ = time.Parse("2006/01/02 15:04 MST", period[1]+" JST")
	}

	(*eventinfo).EventStatus = "BeingHeld"
	if (*eventinfo).Start_time.After(time.Now()) {
		(*eventinfo).EventStatus = "NotHeldYet"
	} else if (*eventinfo).End_time.Before(time.Now()) {
		(*eventinfo).EventStatus = "Over"
	}

	//	イベントに参加しているルームの数を求めます。
	//	参加ルーム数と表示されているルームの数は違うので、ここで取得したルームの数を以下の処理で使うわけではありません。
	SNoEntry := doc.Find("p.ta-r").Text()
	fmt.Sscanf(SNoEntry, "%d", &((*eventinfo).NoEntry))
	log.Printf("[%s]\n[%s] [%s] (*event).EventStatus=%s NoEntry=%d\n",
		(*eventinfo).Event_name,
		(*eventinfo).Start_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).End_time.Format("2006/01/02 15:04 MST"),
		(*eventinfo).EventStatus, (*eventinfo).NoEntry)
	log.Printf("breg=%d ereg=%d\n", breg, ereg)

	//	eventno, _, _ := SelectEventNoAndName(eventidorfilename)
	//	log.Printf(" eventno=%d\n", eventno)
	//	(*eventinfo).Event_no = eventno

	eia := strings.Split(eventid, "?")
	bia := strings.Split(eia[1], "=")
	blockid, _ := strconv.Atoi(bia[1])

	/*
		event_id := 30030
		event_id := 31947
	*/

	ebr, err := srapi.GetEventBlockRanking(client, event_id, blockid, breg, ereg)
	if err != nil {
		log.Printf("GetEventBlockRanking() err=%s\n", err.Error())
		status = 1
		return
	}

	ReplaceString := "/r/"

	for i, br := range ebr.Block_ranking_list {

		var roominfo RoomInfo

		roominfo.ID = br.Room_id
		roominfo.Userno, _ = strconv.Atoi(roominfo.ID)

		roominfo.Account = strings.Replace(br.Room_url_key, ReplaceString, "", -1)
		roominfo.Account = strings.Replace(roominfo.Account, "/", "", -1)

		roominfo.Name = br.Room_name

		roominfo.Irank = i + 1

		*roominfolist = append(*roominfolist, roominfo)

	}

	(*eventinfo).NoRoom = len(*roominfolist)

	log.Printf(" GetEventInfAndRoomList() len(*roominfolist)=%d\n", len(*roominfolist))

	return
}
