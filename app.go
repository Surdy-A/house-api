package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialise() error {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	// user := os.Getenv("DBUSER")
	// password := os.Getenv("DBPASSWORD")
	// dbname := os.Getenv("DBNAME")

	connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", "root", "Goodman8349**", "realhouse")

	app.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	error_message := map[string]string{"error": err}
	sendResponse(w, statusCode, error_message)
}

func (app *App) getHouses(w http.ResponseWriter, r *http.Request) {
	houses, err := getHouses(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, houses)
}

func (app *App) getHouse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid house ID")
		return
	}

	h := house{ID: key}
	err = h.getHouse(app.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			sendError(w, http.StatusNotFound, "House not found")
		default:
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	sendResponse(w, http.StatusOK, h)
}

func (app *App) createHouse(w http.ResponseWriter, r *http.Request) {
	var h house

	err := json.NewDecoder(r.Body).Decode(&h)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	err = h.createHouse(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusCreated, h)

}

func (app *App) updateHouse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid house ID")
		return
	}
	var h house
	err = json.NewDecoder(r.Body).Decode(&h)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	h.ID = key
	err = h.updateHouse(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, h)

}

func (app *App) deleteHouse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid house ID")
		return
	}
	h := house{ID: key}
	err = h.deleteHouse(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, map[string]string{"result": "successful deletion"})
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/house", app.createHouse).Methods("POST")
	app.Router.HandleFunc("/houses", app.getHouses).Methods("GET")
	app.Router.HandleFunc("/house/{id}", app.getHouse).Methods("GET")
	app.Router.HandleFunc("/house/{id}", app.updateHouse).Methods("PUT")
	app.Router.HandleFunc("/house/{id}", app.deleteHouse).Methods("DELETE")
}
