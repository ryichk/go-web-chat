package main

import (
	"fmt"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// authというCookieの有無をチェック
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 未認証の場合、ログインページにリダイレクト
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// エラー発生
		panic(err.Error())
	} else {
		// 成功。ラップされたハンドラを呼び出す
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// サードパーティへのログイン処理を受け持つ
// パスの形式: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		// ログイン処理
		fmt.Fprintf(w, "TODO: %s login", provider)
	default:
		// 404
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}
