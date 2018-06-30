package chatproject

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// RequestStruct テスト
type RequestStruct struct {
	RequestString string `json:"RequestString"`
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
	// リクエストパラメータのjsonを構造体にパース
	var req RequestStruct
	json.Unmarshal(body, &req)
	// log.Infof(ctx, string(req.RequestString))

	// 構造体をjsonにパース
	result, err := json.Marshal(req)
	// log.Infof(ctx, string(result))

	// エラーがあるか
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスはjsonで返却
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
