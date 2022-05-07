package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// templ is a template
type templateHandler struct {
	once     sync.Once
	filename string
	// ポインタ変数(アドレスが入る)
	templ *template.Template
}

// ServeHTTP handles HTTP requests
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	// flag.Stringは、*string型の値を返す。つまり、フラグの値が保持されているアドレスを返す。
	addr := flag.String("addr", ":8080", "アプリのアドレス")
	// フラグを解釈する
	// コマンドラインで指定された文字列から*addrに情報をセット
	flag.Parse()

	r := newRoom()
	// route
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.run()
	// start web server
	log.Println("webサーバを開始します。ポート: ", *addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
