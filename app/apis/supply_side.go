package apis

import (
	app_db "app/db"
	app_utils "app/utils"
	"fmt"
	"net/http"
)

func publish_auction(w http.ResponseWriter, req *http.Request) {
	title := req.FormValue("title")
	price := req.FormValue("price")
	close_at_str := req.FormValue("close_at")
	close_at, close_at_err := app_utils.MsToTime(close_at_str)
	if title == "" || price == "" || close_at_err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	user_id, _ := req.Cookie("user_id")
	// t := close_at.Format(time.RFC3339)
	// "2023-11-12T00:00:00Z"
	// "2023-11-12T05:30:00+05:30"
	// t := fmt.Sprintf("%s+00:00", close_at.UTC().Format("2006-01-02T15:04:05")) // "2023-11-12T00:00:00+00:00"

	_, err := app_db.SQLExecTimeout(req.Context(), `INSERT INTO auction_services (title,price,bid_highest_price,close_at,published_by) values(?,?,?,?,?)`, title, price, price, app_utils.ToUTCDateString(close_at), user_id.Value)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error in publishing", http.StatusBadGateway)
		return
	}

	fmt.Fprint(w, "ok\n")
}
