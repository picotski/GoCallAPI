package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/picotski/api/models/call"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

type Health struct {
	Status int
	Time   string
}

type PageResponce struct {
	TotalCount int    `json:"totalCount"`
	PrevPage   int    `json:"prevPage"`
	NextPage   int    `json:"nextPage"`
	Calls      []call.Call `json:"calls"`
}

func (a *App) Initialize(user, password, dbName, hostAddr string) {
	connectionString := fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		hostAddr,
		user,
		password,
		dbName,
	)

	var err error

	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Print("Err")
		log.Fatal(err)
	}

	// Delete table on start
	if err := call.DeleteCallTable(a.DB); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Table deleted")
	}

	// Init table on start
	if err := call.CreateCallTable(a.DB); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Table created")
	}

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) Run(address string) {
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/calls", a.getCalls).Methods("GET")
	a.Router.HandleFunc("/call", a.createCall).Methods("POST")
	a.Router.HandleFunc("/call/{id:[0-9]+}", a.getCall).Methods("GET")
	a.Router.HandleFunc("/call/{id:[0-9]+}", a.updateCall).Methods("PUT")
	a.Router.HandleFunc("/call/{id:[0-9]+}", a.deleteCall).Methods("DELETE")
	a.Router.HandleFunc("/stop/{id:[0-9]+}", a.endCall).Methods("GET")
	a.Router.HandleFunc("/health", a.healthCheck).Methods("GET")
}

// Get by page
func (a *App) getCalls(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	count, _ := strconv.Atoi(r.URL.Query().Get("count"))

	if count > 10 || count < 1 {
		count = 10
	}
	if page < 1 {
		page = 1
	}

	totalCount, err := call.CountCalls(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	calls, err := call.GetCalls(a.DB, page-1, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}
	nextPage := page + 1
	if totalCount < nextPage*count-(count-1) {
		nextPage = page
	}

	res := PageResponce{
		TotalCount: totalCount,
		PrevPage:   prevPage,
		NextPage:   nextPage,
		Calls:      calls,
	}

	respondWithJSON(w, http.StatusOK, res)
}

// Get call by id
func (a *App) getCall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid call ID")
		return
	}

	c := call.Call{ID: id}
	if err := c.GetCall(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Call not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

// Create a new call
func (a *App) createCall(w http.ResponseWriter, r *http.Request) {
	var c call.Call
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	c.StartCall()

	if err := c.CreateCall(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, c)
}

// End call
func (a *App) endCall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid call ID")
		return
	}

	c := call.Call{ID: id}
	if err := c.GetCall(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Call not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if c.Status == "Ended" {
		respondWithError(w, http.StatusBadRequest, "Call already ended")
		return
	}

	if err := c.StopCall(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

// Update call by id
func (a *App) updateCall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid call ID")
		return
	}

	var c call.Call
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}
	defer r.Body.Close()

	c.ID = id

	if err := c.UpdateCall(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

// Delete call by id
func (a *App) deleteCall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Call ID")
		return
	}

	c := call.Call{ID: id}
	if err := c.DeleteCall(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	responce := Health{
		Time:   time.Now().String(),
		Status: http.StatusOK,
	}

	respondWithJSON(w, http.StatusOK, responce)
}

// Helper
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
