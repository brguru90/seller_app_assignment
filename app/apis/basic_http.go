package apis

import (
	app_db "app/db"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func test(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	fmt.Fprint(w, "ok\n")
}

func sign_in(w http.ResponseWriter, req *http.Request) {
	name := req.FormValue("name")
	email := req.FormValue("email")
	is_supplier := req.FormValue("is_supplier")
	if name == "" || email == "" {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	res, err := app_db.SQLExecTimeout(req.Context(), `INSERT INTO users(name,email) values(?,?)`, name, email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error in sign in", http.StatusBadGateway)
		return
	}
	user_id, _ := res.LastInsertId()
	expiry := time.Now().Add(time.Minute * 10)
	http.SetCookie(w, &http.Cookie{Name: "email", Value: email, HttpOnly: false, Expires: expiry})
	http.SetCookie(w, &http.Cookie{Name: "user_id", Value: fmt.Sprintf("%d", user_id), HttpOnly: false, Expires: expiry})
	http.SetCookie(w, &http.Cookie{Name: "name", Value: name, HttpOnly: false, Expires: expiry})
	http.SetCookie(w, &http.Cookie{Name: "is_supplier", Value: is_supplier, HttpOnly: false, Expires: expiry})
	fmt.Fprint(w, "ok\n")
}

func sign_out(w http.ResponseWriter, req *http.Request) {
	DeleteCookie(w, "user_id")
	DeleteCookie(w, "email")
	DeleteCookie(w, "name")
	DeleteCookie(w, "is_supplier")
	fmt.Fprint(w, "ok\n")
}

func RouteStatic(handler http.Handler, next_handler http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, ".html") {
			handler.ServeHTTP(w, req)
		} else {
			next_handler.ServeHTTP(w, req)
		}
	}
}

func BasicHttpServer() {
	mux := http.NewServeMux()
	{
		fs := http.FileServer(http.Dir("./static"))
		fs2 := http.FileServer(http.Dir("./static/dir"))
		mux.Handle("/", fs)
		mux.Handle("/static", fs2)
	}
	{
		mux.HandleFunc("/test", test)
		mux.Handle("/sign_in", M_POST(sign_in))
		mux.HandleFunc("/sign_out", sign_out)
		mux.HandleFunc("/list_supply", list_supply)
		mux.HandleFunc("/get_service_price", get_service_price)
		mux.Handle("/bid_service", M_POST_HANDLER(WithAuth(bid_service)))
		mux.Handle("/publish_auction", M_POST_HANDLER(WithAuth(publish_auction)))
	}

	host := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	fmt.Printf("running server at %s ...\n", host)
	// if err := http.ListenAndServe(host, http.HandlerFunc(RouteStatic(fs, mux))); err != nil {
	// 	fmt.Println("some error")
	// 	fmt.Println(err)
	// }
	if err := http.ListenAndServe(host, mux); err != nil {
		fmt.Println("some error")
		fmt.Println(err)
	}
}
