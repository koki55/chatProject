package message

import (
	"chatProject/src/user"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// Messages メッセージ
type Messages struct {
	ID      int64  `datastore:"-"`
	Message string `json:"Message"`
	RoomID  int64  `datastore:"RoomID"`
	UserID  int64  `datastore:"UserID"`
	IconID  int32  `datastore:"IconID"`
	Name    string
	Date    time.Time
}

// LoadHandler メッセージを取得
func LoadHandler(w http.ResponseWriter, r *http.Request) {
	// コンテキスト作成
	ctx := appengine.NewContext(r)

	// リクエストパラメータ
	body, _ := ioutil.ReadAll(r.Body)
	var message Messages
	json.Unmarshal(body, &message)

	messages, err := getMessagesByRoomID(ctx, message.RoomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 構造体をjsonにパース
	result, err := json.Marshal(messages)
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

// SendHandler メッセージ送信
func SendHandler(w http.ResponseWriter, r *http.Request) {
	// ログを残す等に必要
	ctx := appengine.NewContext(r)

	// リクエストパラメータ
	body, _ := ioutil.ReadAll(r.Body)
	log.Infof(ctx, string(body))

	// インスタンス化
	message := Messages{
		Date: time.Now(),
	}
	// リクエストを構造体に
	json.Unmarshal(body, &message)

	// リクエストのユーザIDからユーザ情報を取得
	userKey := datastore.NewKey(ctx, "Users", "", message.UserID, nil)
	var user user.Users
	err := datastore.Get(ctx, userKey, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// メッセージにユーザ情報を追加
	message.IconID = user.IconID
	message.Name = user.Name

	// 作成するmessageのIDを取得
	newkey := datastore.NewIncompleteKey(ctx, "Messages", nil)
	messageKey, err := datastore.Put(ctx, newkey, &message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	message.ID = messageKey.IntID()

	// 構造体をjsonにパース
	result, err := json.Marshal(message)
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

// ルームIDに紐づくメッセージを100件まで取得
func getMessagesByRoomID(ctx context.Context, RoomID int64) ([]Messages, error) {
	query := datastore.NewQuery("Messages").Filter("RoomID =", RoomID).Limit(100)

	messages := make([]Messages, 0)
	iter := query.Run(ctx)
	for {
		var message Messages
		key, err := iter.Next(&message)
		if err == datastore.Done {
			break
		} else if err != nil {
			return nil, err
		}
		message.ID = key.IntID()
		messages = append(messages, message)
	}
	return messages, nil
}
