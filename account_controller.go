package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	// "strconv"
)

var DBClient IBoltClient

////////////////////////// Replicator related endpoints starts here /////////////////////////////

func (a *App) GetAccount(w http.ResponseWriter, req *http.Request) {

	reqLogger := req.Context().Value(keyLoggerMiddleware).(*log.Entry)

	// Read the 'accountId' path parameter from the mux map
	var accountId = mux.Vars(req)["id"]

	if accountId != "" {

		// Read the account struct BoltDB
		account, err := DBClient.QueryAccount(accountId)

		// If err, return a 404
		if err != nil {
			reqLogger.Error("No Data found " + accountId + ": " + err.Error())
			reqLogger.Error("Error in GetAccount method", "")
			reqLogger.Info(a.Module, "_Failure")
			respondWithError(w, req, http.StatusNotFound, "No Data found")
			return
		}

		account_list := make([]Account, 0, 1)
		account_list = append(account_list, account)
		reqLogger.Info("Response data - ", account)
		reqLogger.Info(a.Module, "_Success")
		respondWithJSON(w, req, http.StatusOK, account_list)
	} else {
		reqLogger.Error("Please pass all parameters")
		reqLogger.Error("Error in GetAccount method", "")
		reqLogger.Info(a.Module, "_Failure")
		respondWithError(w, req, http.StatusBadRequest, "Please pass all parameters")
		return
	}
}

func (a *App) NewAccount(w http.ResponseWriter, req *http.Request) {
	reqLogger := req.Context().Value(keyLoggerMiddleware).(*log.Entry)
	reqLogger.Info("Raw request data - ", req.Body)
	var acc Account
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&acc)
	reqLogger.Info("Converted Request data - ", acc)
	if err != nil {
		reqLogger.Error("Unable to parse json from request body")
		reqLogger.Error("Error in NewAccount method - ", err.Error())
		reqLogger.Info(a.Module, "_Failure")
		respondWithError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	if acc.Id != "" && acc.Name != "" {

		// Create account in BoltDB
		account, err := DBClient.CreateAccount(acc)
		// If err, return a 404
		if err != nil {
			reqLogger.Error("Error Creating new account : " + err.Error())
			reqLogger.Error("Error in NewAccount method", "")
			reqLogger.Info(a.Module, "_Failure")
			respondWithError(w, req, http.StatusInternalServerError, "Error Creating new account")
			return
		}

		reqLogger.Info("Response data - ", account)
		reqLogger.Info(a.Module, "_Success")
		respondWithJSON(w, req, http.StatusCreated, account)
		//respondWithJSON(w, req, http.StatusCreated, map[string]string{"result": "Key added/updated"})
	} else {
		reqLogger.Error("Please pass all parameters")
		reqLogger.Error("Error in NewAccount method", "")
		reqLogger.Info(a.Module, "_Failure")
		respondWithError(w, req, http.StatusBadRequest, "Please pass all parameters")
		return
	}
}

func TestEndpoint(w http.ResponseWriter, req *http.Request) {
	decoded := context.Get(req, "decoded")
	// var user User
	// mapstructure.Decode(decoded.(jwt.MapClaims), &user)
	json.NewEncoder(w).Encode(decoded)
}

func (a *App) CreateToken(w http.ResponseWriter, req *http.Request) {
	reqLogger := req.Context().Value(keyLoggerMiddleware).(*log.Entry)
	reqLogger.Info("Raw request data - ", req.Body)
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)
	if user.Username != "" && user.Password != "" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"password": user.Password,
		})
		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			reqLogger.Error("Error signing string : " + err.Error())
			reqLogger.Error("Error in CreateToken method", "")
			reqLogger.Info(a.Module, "_Failure")
			respondWithError(w, req, http.StatusInternalServerError, "Error signing string")
			return
		}
		respondWithJSON(w, req, http.StatusOK, map[string]string{"Token": tokenString})
		//json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
	} else {
		reqLogger.Error("Please pass all parameters")
		reqLogger.Error("Error in CreateToken method", "")
		reqLogger.Info(a.Module, "_Failure")
		respondWithError(w, req, http.StatusBadRequest, "Please pass all parameters")
		return
	}

}

func (a *App) ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte("secret"), nil
				})
				if error != nil {
					json.NewEncoder(w).Encode(Exception{Message: error.Error()})
					return
				}
				if token.Valid {
					//context.Set(req, "decoded", token.Claims)
					context.Set(req, "decoded", true)
					next(w, req)
				} else {
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
}
