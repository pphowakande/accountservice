package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func (a *App) InitTestDB() {
	DBClient = &BoltClient{}
	DBClient.OpenBoltDb()
	DBClient.Seed()
}

func TestCreateToken(t *testing.T) {
	AppCtx.InitTestDB()
	r := mux.NewRouter()
	r.HandleFunc("/authenticate", headerSetter(AppCtx.requestLogger(AppCtx.CreateToken))).Methods("POST")

	Convey("Given I have valid request", t, func() {
		input, _ := json.Marshal(User{Username: "poonam", Password: "poonam"})
		body := bytes.NewBuffer(input)
		req, _ := http.NewRequest("POST", "/authenticate", body)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		Convey("status code should be 200", func() {
			So(res.Code, ShouldEqual, http.StatusOK)
		})
	})

	Convey("Given I have invalid request", t, func() {
		input, _ := json.Marshal(User{Username: "", Password: ""})
		body := bytes.NewBuffer(input)
		req, _ := http.NewRequest("POST", "/authenticate", body)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		Convey("status code should be 400", func() {
			So(res.Code, ShouldEqual, http.StatusBadRequest)
		})
	})
}

func TestGetAccount(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/account/{id}", headerSetter(AppCtx.requestLogger(AppCtx.ValidateMiddleware(AppCtx.GetAccount)))).Methods("GET")
	Convey("Given I have valid get request", t, func() {
		req, _ := http.NewRequest("GET", "/account/11", nil)
		req.Header.Set("authorization", "JWT eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXNzd29yZCI6InBvb25hbSIsInVzZXJuYW1lIjoicG9vbmFtIn0.WV6MnDQif1niCNiJb9vZAc8ECwIOj37rItfl8hknf58")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		Convey("status code should be 200", func() {
			So(res.Code, ShouldEqual, http.StatusOK)

			Convey("accounts list should be 1", func() {
				var accounts []Account
				json.NewDecoder(res.Body).Decode(&accounts)
				So(len(accounts), ShouldEqual, 1)
			})
		})
	})

	Convey("Given I have invalid get request (wrong route)", t, func() {
		req, _ := http.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		Convey("status code should be 404", func() {
			So(recorder.Code, ShouldEqual, http.StatusNotFound)
		})
	})
}

func TestNewAccount(t *testing.T) {
	r := mux.NewRouter()

	r.HandleFunc("/account", headerSetter(AppCtx.requestLogger(AppCtx.NewAccount))).Methods("POST")
	Convey("Given I have valid request", t, func() {
		input, _ := json.Marshal(Account{Id: "2", Name: "person2"})
		body := bytes.NewBuffer(input)
		req, _ := http.NewRequest("POST", "/account", body)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		Convey("status code should be 201", func() {
			So(res.Code, ShouldEqual, http.StatusCreated)
		})
	})

	Convey("Given I have invalid request", t, func() {
		input, _ := json.Marshal(Account{Id: "", Name: ""})
		body := bytes.NewBuffer(input)
		req, _ := http.NewRequest("POST", "/account", body)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		Convey("status code should be 400", func() {
			So(res.Code, ShouldEqual, http.StatusBadRequest)
		})
	})
}
