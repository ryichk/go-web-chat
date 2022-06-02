package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"go-web-chat/trace"
	"log"
	"net/http"
	"os"
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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

func main() {
	// flag.Stringは、*string型の値を返す。つまり、フラグの値が保持されているアドレスを返す。
	addr := flag.String("addr", ":8080", "アプリのアドレス")
	// フラグを解釈する
	// コマンドラインで指定された文字列から*addrに情報をセット
	flag.Parse()

	if err := godotenv.Load("../.env"); err != nil {
		log.Println(".envを読み込めませんでした: %v", err)
	}
	// Gomniauthのセットアップ
	gomniauth.SetSecurityKey(os.Getenv("GOMNIAUTH_SECURITY_KEY"))
	gomniauth.WithProviders(
		google.New(os.Getenv("GOOGLE_OAUTH2_CLIENT_ID"), os.Getenv("GOOGLE_OAUTH2_API_KEY"), "http://localhost:3000/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	// route
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	// http.Handlerインタフェースを実装していないハンドラもパスの関連付けを行える
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/room", r)
	go r.run()
	// start web server
	log.Println("webサーバを開始します。ポート: ", *addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServer:", err)
	}
}
