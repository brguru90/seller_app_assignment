package app_db

import (
	app_utils "app/utils"
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"os"
)

var DATABASE_CONN *sql.DB

func SQLExecTimeout(q_ctx context.Context, query string, args ...any) (sql.Result, error) {
	timeout_sec := time.Duration(app_utils.StrToInt64(os.Getenv("DB_TIMEOUT_SEC")))
	if q_ctx == nil {
		q_ctx = context.Background()
	}
	ctx, _ := context.WithTimeout(q_ctx, timeout_sec*time.Second)
	return DATABASE_CONN.ExecContext(ctx, query, args...)
}

func SQLQueryTimeout(q_ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	timeout_sec := time.Duration(app_utils.StrToInt64(os.Getenv("DB_TIMEOUT_SEC")))
	if q_ctx == nil {
		q_ctx = context.Background()
	}
	ctx, _ := context.WithTimeout(q_ctx, timeout_sec*time.Second)
	return DATABASE_CONN.QueryContext(ctx, query, args...)
}

func ConnectToDB() {

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/mysql", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	conn, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Printf("connection string: %s\n", connectionString)
		panic(err)
	}
	defer conn.Close()
	_, err = conn.Exec(
		fmt.Sprintf(`
			CREATE DATABASE IF NOT EXISTS %s;
			`,
			os.Getenv("DB_NAME"),
		),
	)
	if err != nil {
		fmt.Printf("connection string: %s\n", connectionString)
		panic(err)
	}

	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	DATABASE_CONN, err = sql.Open("mysql", dbConnectionString)
	if err != nil {
		fmt.Printf("connection string: %s\n", dbConnectionString)
		panic(err)
	}
	fmt.Println("db connected.")
}

func DBInit() {
	_, err := SQLExecTimeout(nil, `
	CREATE TABLE IF NOT EXISTS users (
		user_id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
		name varchar(255) NOT NULL,
		email varchar(255) NOT NULL UNIQUE
	)
	`)
	if err != nil {
		panic(err)
	}

	_, err = SQLExecTimeout(nil, `
	CREATE TABLE IF NOT EXISTS auction_services (
		auction_id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
		title varchar(255) NOT NULL,
		price FLOAT NOT NULL,
		close_at DATETIME NOT NULL,
		published_by int NOT NULL,
		bid_highest_price FLOAT NOT NULL,
		bidding_user int,
		CONSTRAINT fk_bidding  FOREIGN KEY (bidding_user) REFERENCES users(user_id),
		CONSTRAINT fk_publish  FOREIGN KEY (published_by) REFERENCES users(user_id) 
	)
	`)
	if err != nil {
		panic(err)
	}
	fmt.Println("db initialized.")
}
