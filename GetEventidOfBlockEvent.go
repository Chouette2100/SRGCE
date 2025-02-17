/*
	!

Copyright © 2022 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php
*/
package main

import (
	// "fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	// "encoding/json"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Block struct {
	Label   string
	Block_id int
}

type BlockInf struct {
	Show_rank_label string
	Block_list      []Block
}

type BlockInfList struct {
	Blockinf []BlockInf
}

// ブロックイベントの子のイベントのeventidを取得する。
func GetEventidOfBlockEvent(
	eventid string, //	ブロックイベントの親イベントのeventid
) (
	blockinflist BlockInfList, //	このブロックイベントのラベルとブロック番号のペア

	err error,
) {

	var doc *goquery.Document

	url := "https://www.showroom-live.com/event/" + eventid
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println("Error parsing document:", err)
		return
	}

	var scriptContent string
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(text, "window.__NUXT__") {
			scriptContent = text
			return
		}
	})

	if scriptContent == "" {
		log.Println("Script with window.__NUXT__ not found")
		return
	}

	// 1. まずshow_rank_labelとblock_listの塊を見つける
	blockGroupPattern := `{show_rank_label:([^,]+),block_list:\[(.*?)\]}`
	blockGroupPatternWithQuotes := `{show_rank_label:"([^"]+)",block_list:\[(.*?)\]}`

	reBlockGroup := regexp.MustCompile(blockGroupPattern)
	reBlockGroupWithQuotes := regexp.MustCompile(blockGroupPatternWithQuotes)

	var blockGroups []BlockInf

	atoi := func(s string) int {
		v, e := strconv.Atoi(s)
		if e != nil {
			return 0
		}
		return v
	}

	// 変数パターンでのマッチング
	blockGroupMatches := reBlockGroup.FindAllStringSubmatch(scriptContent, -1)
	for _, match := range blockGroupMatches {
		if len(match) >= 3 {
			group := BlockInf{
				Show_rank_label: match[1],
				Block_list:      []Block{},
			}

			// 2. block_list内の各アイテムを見つける
			blockItemPattern := `{block_id:([^,]+),label:([^}]+)}`
			reBlockItem := regexp.MustCompile(blockItemPattern)

			blockItemMatches := reBlockItem.FindAllStringSubmatch(match[2], -1)
			for _, itemMatch := range blockItemMatches {
				if len(itemMatch) >= 3 {
					group.Block_list = append(group.Block_list, Block{
						Block_id: atoi(itemMatch[1]),
						Label:   strings.TrimSpace(itemMatch[2]),
					})
				}
			}

			blockGroups = append(blockGroups, group)
		}
	}

	// 文字列リテラルパターンでのマッチング
	blockGroupMatchesWithQuotes := reBlockGroupWithQuotes.FindAllStringSubmatch(scriptContent, -1)
	for _, match := range blockGroupMatchesWithQuotes {
		if len(match) >= 3 {
			group := BlockInf{
				Show_rank_label: match[1],
				Block_list:      []Block{},
			}

			// block_list内の各アイテムを見つける（ラベルが変数と文字列リテラルの両方のパターンに対応）
			blockItemPattern := `{block_id:([^,]+),label:([^}]+)}`
			reBlockItem := regexp.MustCompile(blockItemPattern)

			blockItemMatches := reBlockItem.FindAllStringSubmatch(match[2], -1)
			for _, itemMatch := range blockItemMatches {
				if len(itemMatch) >= 3 {
					group.Block_list = append(group.Block_list, Block{
						Block_id: atoi(itemMatch[1]),
						Label:   strings.TrimSpace(itemMatch[2]),
					})
				}
			}

			blockGroups = append(blockGroups, group)
		}
	}

	// 結果の表示
	for _, group := range blockGroups {
		log.Printf("ShowRankLabel: %s\n", group.Show_rank_label)
		for _, item := range group.Block_list {
			log.Printf("  BlockID: %d, Label: %s\n", item.Block_id, item.Label)
		}
		log.Println()
	}
	blockinflist.Blockinf = blockGroups

	/*
		//	画面からのデータ取得部分は次を参考にしました。
		//		はじめてのGo言語：Golangでスクレイピングをしてみた
		//		https://qiita.com/ryo_naka/items/a08d70f003fac7fb0808

		var doc *goquery.Document

		//	URLからドキュメントを作成します
		_url := "https://www.showroom-live.com/event/" + eventid
		resp, error := http.Get(_url)
		if error != nil {
			log.Printf("GetEventInf() http.Get() err=%s\n", error.Error())
			err = fmt.Errorf("http.Get(): %w", error)
			return
		}
		defer resp.Body.Close()

		doc, error = goquery.NewDocumentFromReader(resp.Body)
		if error != nil {
			log.Printf("GetEventInf() goquery.NewDocumentFromReader() err=<%s>.\n", error.Error())
			err = fmt.Errorf("goquery.NewDocumentFromReader(): %w", error)
			return
		}

		//	ブロック情報がJSONとして得られる
		//	tjson, bl := doc.Find(".js-event-lower-cate-section div div event-block").Attr("data-list")
		tjson, bl := doc.Find("#js-event-block > event-block").Attr("data-list")
		if !bl {
			err = fmt.Errorf("doc.Find().Attr(): %t", bl)
			return
		}
		//	tjson = tjson[ 1: len(tjson)-1]

		err = json.NewDecoder(strings.NewReader(tjson)).Decode(&blockinflist.Blockinf)
		if err != nil {
			err = fmt.Errorf("json.NewDecoder().Decode(): %w", err)
			return
		}
	*/

	return
}
