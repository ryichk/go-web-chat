package main

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
	"go-web-chat/trace"
	"log"
	"net/http"
)

type room struct {
	// forwardは他のクライアントに転送するメッセージを保持するチャネル
	forward chan *message

	// joinとleaveはマップclientsへの追加・削除に使用される
	// joinはチャットルームに参加しようとしているクライアントのチャネル
	join chan *client
	// leaveはチャットルームから退出しようとしているクライアントのチャネル
	leave chan *client
	// clientsには在室している全てのクライアントが保持される
	// チャネルを使わずにマップclientsを直接操作することは望ましくない
	// 複数のgoroutineがマップを同時に変更する可能性があり、メモリ破壊など予期せぬ状態になりうる
	clients map[*client]bool
	// tracerはチャットルーム上で行われた操作のログを受け取る
	tracer trace.Tracer
}

func (r *room) run() {
	for {
		// 共有されているメモリに対して同期化や変更がいくつか必要な場合、select文を利用する
		select {
		//case節のコードが同時に実行されることはない
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("A new client has joined.")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("The client has left.")
		case msg := <-r.forward:
			r.tracer.Trace("Received a message: ", msg.Message)
			// 全てのクライアントにメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					// メッセージ送信
					r.tracer.Trace(" -- Sent to client")
				default:
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- Failed to send. Clean up the client.")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// WebSocketを利用するためにwebsocket.Upgrader型を使ってHTTP接続をアップグレードする必要がある
var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

// *room型をhttp.Handler型に適合
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// WebSocketコネクションを取得
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("クッキーの取得に失敗しました:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan *message, messageBufferSize),
		room:   r,
		// エンコードされたクッキーの値をマップのオブジェクトへ復元
		userData: objx.MustFromBase64(authCookie.Value),
	}
	// クライアント入室
	r.join <- client
	// クライアントの終了時に退室処理
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}
