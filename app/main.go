package main

import (
	"app/apis"
	app_db "app/db"
)

func main() {
	app_db.ConnectToDB()
	app_db.DBInit()
	apis.BasicHttpServer()
}
