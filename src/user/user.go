package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// Users ユーザ
type Users struct {
	UserID int64 `datastore:"-"`
	IconID int32
	Name   string
}

// RegistHandler ユーザ登録
func RegistHandler(w http.ResponseWriter, r *http.Request) {
	// ログを残す等に必要
	ctx := appengine.NewContext(r)

	key := datastore.NewIncompleteKey(ctx, "Users", nil)

	// リクエストパラメータ
	body, _ := ioutil.ReadAll(r.Body)
	log.Infof(ctx, string(body))

	// インスタンス化
	var user Users

	json.Unmarshal(body, &user)

	key, err := datastore.Put(ctx, key, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.UserID = key.IntID()
	// 構造体をjsonにパース
	result, err := json.Marshal(user)
	log.Infof(ctx, string(result))

	// エラーがあるか
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// レスポンスはjsonで返却
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
