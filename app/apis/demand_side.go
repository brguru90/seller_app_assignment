package apis

import (
	app_db "app/db"
	app_utils "app/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Auction struct {
	AuctionID   string `json:"auction_id"`
	Title       string `json:"title"`
	PublishedBy string `json:"published_by,omitempty"`
	BidPrice    string `json:"bid_highest_price"`
	CloseAt     string `json:"close_at"`
	BidedBy     string `json:"bided_by,omitempty"`
	BidedById   string `json:"bided_by_id,omitempty"`
}

func list_supply(w http.ResponseWriter, req *http.Request) {
	close_at := time.Now()
	var auctions_row *sql.Rows
	var err error
	if req.URL.Query().Get("all") != "true" {
		auctions_row, err = app_db.SQLQueryTimeout(req.Context(), `SELECT auction_id,title,bid_highest_price,close_at,"" AS bidding_user,"" AS name,published_by from auction_services where close_at > ?;`, close_at)
	} else {
		// user_id, err := req.Cookie("user_id")
		// if err != nil {
		// 	fmt.Println(err)
		// 	http.Error(w, "please login as supplier", http.StatusBadGateway)
		// 	return
		// }
		auctions_row, err = app_db.SQLQueryTimeout(req.Context(), `SELECT auction_id,title,bid_highest_price,close_at,bidding_user,name AS bidder_name,"" AS published_by from auction_services LEFT JOIN users ON bidding_user=user_id;`)
	}
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error", http.StatusBadGateway)
		return
	}
	auctions := []Auction{}
	for auctions_row.Next() {
		var auction Auction
		auctions_row.Scan(&auction.AuctionID, &auction.Title, &auction.BidPrice, &auction.CloseAt, &auction.BidedById, &auction.BidedBy, &auction.PublishedBy)
		auctions = append(auctions, auction)
	}
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(auctions)
	result, _ := json.MarshalIndent(auctions, "", "  ")
	fmt.Fprint(w, string(result))
}

func get_service_price(w http.ResponseWriter, req *http.Request) {
	var bid_highest_price string
	err := app_db.DATABASE_CONN.QueryRow("SELECT bid_highest_price from auction_services where auction_id = ?;", req.URL.Query().Get("auction_id")).Scan(&bid_highest_price)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error", http.StatusBadGateway)
		return
	}
	fmt.Fprint(w, bid_highest_price)
}

func bid_service(w http.ResponseWriter, req *http.Request) {
	auction_id := req.FormValue("auction_id")
	bid_price := req.FormValue("bid_price")

	if auction_id == "" || bid_price == "" {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	close_at := time.Now()

	var bid_highest_price string
	err := app_db.DATABASE_CONN.QueryRow("SELECT bid_highest_price from auction_services where auction_id = ?;", auction_id).Scan(&bid_highest_price)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error", http.StatusBadGateway)
		return
	}

	if app_utils.StrToInt64(bid_price) <= app_utils.StrToInt64(bid_highest_price) {
		fmt.Println(err)
		http.Error(w, "Lower price not allowed", http.StatusForbidden)
		return
	}

	user_id, _ := req.Cookie("user_id")

	res, err := app_db.SQLExecTimeout(req.Context(), `UPDATE auction_services SET bid_highest_price = ?,bidding_user = ? WHERE auction_id = ? AND published_by!=? AND bid_highest_price < ? AND close_at > ?;`, bid_price, user_id.Value, auction_id, user_id.Value, bid_price, close_at)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error in bidding for auction", http.StatusBadGateway)
		return
	}
	count, err := res.RowsAffected()
	if err != nil || count == 0 {
		fmt.Println(err)
		http.Error(w, "Bidding freezed", http.StatusBadGateway)
		return
	}

	fmt.Fprint(w, "success\n")
}
