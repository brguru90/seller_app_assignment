package apis

import (
	"net/http"
	"time"
)

func M_POST(handler func(w http.ResponseWriter, req *http.Request)) http.Handler {
	next := http.HandlerFunc(handler)
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, req)
	})
}

func M_POST_HANDLER(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, req)
	})
}

func WithAuth(handler func(w http.ResponseWriter, req *http.Request)) http.Handler {
	next := http.HandlerFunc(handler)
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		email, err := req.Cookie("email")
		if err != nil || email.Value == "" {
			http.Error(w, "Unauthorized !!!", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, req)
	})
}

func WithAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		email, err := req.Cookie("email")
		if err != nil || email.Value == "" {
			http.Error(w, "Unauthorized !!!", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, req)
	})
}

func DeleteCookie(w http.ResponseWriter, key string) {
	c := &http.Cookie{
		Name:     key,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}
