package chatproject

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Ping テスト
// 名前がintなのにstringなのはご愛嬌
type Ping struct {
	Testint string `json:"Testint"`
}

// 初期化　最初に呼ばれる GAEではmain()は呼ばれない
func init() {
	http.HandleFunc("/ping", pingHandler)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	// ログを残す等に必要
	// ctx := appengine.NewContext(r)

	// リクエストパラメータ
	body, _ := ioutil.ReadAll(r.Body)
	// log.Infof(ctx, string(body))

	// レスポンスのテスト
	// リクエストのjsonを構造体に挿入してパース
	var ping Ping
	json.Unmarshal(body, &ping)
	// log.Infof(ctx, string(ping.Testint))

	// 構造体をjsonにパース
	res, err := json.Marshal(ping)
	// log.Infof(ctx, string(res))

	// エラーがあるか
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスはjsonで返却
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
