package room

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

// Rooms ルーム
type Rooms struct {
	RoomID int64 `datastore:"-"`
	Title  string
	Date   time.Time
}

// ListHandler ルーム一覧を取得
func ListHandler(w http.ResponseWriter, r *http.Request) {

	// ルーム一覧
	ctx := appengine.NewContext(r)
	rooms, err := getAllRooms(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 構造体をjsonにパース
	result, err := json.Marshal(rooms)
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

// CreateHandler ルーム作成
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	key := datastore.NewIncompleteKey(ctx, "Rooms", nil)

	// リクエストパラメータ
	body, _ := ioutil.ReadAll(r.Body)
	log.Infof(ctx, string(body))

	room := Rooms{
		Date: time.Now(),
	}

	// リクエストを構造体に
	json.Unmarshal(body, &room)

	key, err := datastore.Put(ctx, key, &room)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	room.RoomID = key.IntID()

	// 構造体をjsonにパース
	result, err := json.Marshal(room)
	log.Infof(ctx, string(result))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// レスポンスはjsonで返却
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

// 100件までルームを取得
func getAllRooms(ctx context.Context) ([]Rooms, error) {
	query := datastore.NewQuery("Rooms").Order("-Date").Limit(100)

	rooms := make([]Rooms, 0)
	iter := query.Run(ctx)
	for {
		var room Rooms
		key, err := iter.Next(&room)
		if err == datastore.Done {
			break
		} else if err != nil {
			return nil, err
		}
		room.RoomID = key.IntID()
		rooms = append(rooms, room)
	}
	return rooms, nil
}
