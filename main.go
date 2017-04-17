package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"

	"io/ioutil"
)

type apelidoRow struct {
	Id int `json:"id"`
	Apelido string `json:"apelido"`
	Owner string `json:"owner"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", list).Methods("GET")
	router.HandleFunc("/add", post).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}

func list(w http.ResponseWriter, r *http.Request) {
	db := connect()
	rows, err := db.Query("select * from apelidos")
	errorCheck(err, "Error during query ")

	var apelidoList []apelidoRow
	for rows.Next() {
		var id int
		var apelido string
		var owner string

		err = rows.Scan(&id, &apelido, &owner)
		errorCheck(err, "error scanning rows: ")

		apelidoList = append(apelidoList, apelidoRow{id, apelido, owner})
	}

	jsonify,err := json.Marshal(apelidoList)
	errorCheck(err)
	w.Write(jsonify)
}

func post(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token != "barazinho" {
		w.Write([]byte("end game"))

		return
	}

	db := connect()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var apelido apelidoRow
	json.Unmarshal(body, &apelido)

	_, err = db.Exec("insert into apelidos (apelido, owner) values(?, ?)", apelido.Apelido, apelido.Owner)
	errorCheck(err, "error inserting apelido: ")

	w.Write([]byte("ok"))


}

func connect() *sql.DB{
	db, err := sql.Open("mysql", "root:barazinho@tcp(34.203.214.226:3306)/apelidos")
	if err != nil {
		log.Fatal("error getting connection: ", err)
	}

	return db
}

func errorCheck(err error, message ...string) {
	if err != nil {
		log.Fatal(message, err)
	}
}
