package chatproject

import (
	"chatProject/src/message"
	"chatProject/src/room"
	"chatProject/src/user"
	"net/http"
)

// 最初に呼ばれる GAEではmain()は呼ばれない
func init() {
	// ルーム
	http.HandleFunc("/room/list", room.ListHandler)
	http.HandleFunc("/room/create", room.CreateHandler)

	//メッセージ
	http.HandleFunc("/message/send", message.SendHandler)
	http.HandleFunc("/message/load", message.LoadHandler)

	// ユーザ
	http.HandleFunc("/user/regist", user.RegistHandler)
}
