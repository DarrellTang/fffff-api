package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var db *sql.DB
var slogger *zap.SugaredLogger

func main() {
	InitLogger()
	defer slogger.Sync()

	var err error
	slogger.Infof("Connecting to postgres db")
	connStr := fmt.Sprintf("user=postgres password=%s dbname=postgres sslmode=disable host=postgres", os.Getenv("PGPASSWORD"))
	db, err = sql.Open("postgres", connStr)
	if err != nil {
    slogger.Infof("%s",err)
		return
	}
	defer db.Close()

	http.HandleFunc("/nq", nqlist)
	http.HandleFunc("/hq", hqlist)

	slogger.Infof("Serving and listening on port 80")
	err = http.ListenAndServe(":80", nil)
	if err != nil {
    slogger.Infof("%s",err)
		return
	}
}

func nqlist(w http.ResponseWriter, r *http.Request) {
	slogger.Infof("Retrieving normal quality shopping list from db")
	rows, err := db.Query("SELECT json_agg(NqShoppingList) FROM NqShoppingList")
	if err != nil {
    slogger.Infof("%s",err)
		return
	}
	defer rows.Close()

  var json []byte
	rows.Next()
  err = rows.Scan(&json)
  if err != nil {
    fmt.Println(err)
    return
  }

  slogger.Infof("Writing json to http response")
  w.Header().Set("Content-Type", "application/json")
  w.Write(json)
}

func hqlist(w http.ResponseWriter, r *http.Request) {
	slogger.Infof("Retrieving high quality shopping list from db")
	rows, err := db.Query("SELECT json_agg(HqShoppingList) FROM HqShoppingList")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

  var json []byte
  if rows.Next() {
    err = rows.Scan(&json)
    if err != nil {
      fmt.Println(err)
      return
    }
  } else {
    json = []byte("{}")
  }

  slogger.Infof("Writing json to http response")
  w.Header().Set("Content-Type", "application/json")
  w.Write(json)
}

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	slogger = logger.Sugar()
}
