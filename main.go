package main

import (
	"database/sql"
	"encoding/json"
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
		fmt.Println(err)
		return
	}
	defer db.Close()

	http.HandleFunc("/normal", nqlist)
	http.HandleFunc("/high", hqlist)

	slogger.Infof("Serving and listening on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func nqlist(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM NqShoppingList")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			fmt.Println(err)
			return
		}
		row := map[string]interface{}{
			"id":   id,
			"name": name,
		}
		result = append(result, row)
	}

	json.NewEncoder(w).Encode(result)
}

func hqlist(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM HqShoppingList")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id int
		var value string
		err = rows.Scan(&id, &value)
		if err != nil {
			fmt.Println(err)
			return
		}
		row := map[string]interface{}{
			"id":    id,
			"value": value,
		}
		result = append(result, row)
	}

	json.NewEncoder(w).Encode(result)
}

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	slogger = logger.Sugar()
}
