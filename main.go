package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

////
type WrappedType map[string]interface{}

func WrapObject(name string, item interface{}) WrappedType {
	wrapped := map[string]interface{}{
		name: item,
	}
	return wrapped
}

//
// I used in memory implementation of key/value pairs as the first phase before transitioning over to SimpleDB
//
type Pair struct {
	User   string    `json:"User"`
	Key    string    `json:"Key"`
	Value  string    `json:"Value"`
	Last   string    `json:"Last,omitempty"`
	Time   time.Time `json:"Time,omitempty"`
	Server string    `json:"Server,omitempty"`
}

var pairs []Pair

func NewPair(user, key, value string) Pair {
	return Pair{User: user, Key: key, Value: value, Last: "added", Time: time.Now(), Server: serverName.ServerName}
}

func InitPairs() {

	// Pre populate few key/value pairs, for testing, used to seed the actual SimpleDB database
	pairs = append(pairs, NewPair("User1", "keyOne", "valueOne"))
	pairs = append(pairs, NewPair("User2", "keyTwo", "valueTwo"))
	pairs = append(pairs, NewPair("User1", "keyThree", "valueThree"))
	pairs = append(pairs, NewPair("User2", "keyFour", "valueFour"))
	log.Printf("%+v\n", pairs)

	/*
		// Create the real ones in SimpleDB
		for _, item := range pairs {
			PutPairSimpleDB(item.User, item.Key, item.Value, item.Time.String(), item.Server, nil)
		}
	*/

	/*
		// Try it out
		for _, item := range pairs {
			GetPairSimpleDB(item.User, item.Key, nil)
		}

		for _, item := range pairs {
			DeletePairSimpleDB(item.User, item.Key, nil)
		}

		for _, item := range pairs {
			GetPairSimpleDB(item.User, item.Key, nil)
		}
	*/
}

/*
// For in memory database
func FindPair(username string, key string) int {
	for index, item := range pairs {
		if item.User == username && item.Key == key {
			log.Printf("%+v\n", item)
			return index
		}
	}
	return -1
}
*/

func SendUserPairs(username string, w http.ResponseWriter) {
	/*
		// For in memory database
		for _, item := range pairs {
			if item.User == username {
				json.NewEncoder(w).Encode(WrapObject("Match", item))
			}
		}
	*/
	SelectPairsSimpleDB(username, w)
}

func GetPairs(username string, w http.ResponseWriter, r *http.Request) {
	EchoParseRequest(w, r)
	SendUserPairs(username, w)
}

func GetPair(username string, w http.ResponseWriter, r *http.Request) {
	EchoParseRequest(w, r)

	params := mux.Vars(r)

	/*
		// For in memory database
		index := FindPair(username, params["key"])
		if index != -1 {
			json.NewEncoder(w).Encode(WrapObject("Found", pairs[index]))
		} else {
			json.NewEncoder(w).Encode(WrapObject("Not Found", nil))
		}
	*/
	GetPairSimpleDB(username, params["key"], w)
}

func AddUpdatePair(username string, w http.ResponseWriter, r *http.Request) {

	pair := EchoParseRequest(w, r)

	params := mux.Vars(r)

	pair.User = username
	pair.Key = params["key"]
	pair.Last = "created"
	pair.Time = time.Now()
	pair.Server = serverName.ServerName
	json.NewEncoder(w).Encode(WrapObject("Incoming", pair))

	/*
		// For in memory database
		index := FindPair(username, params["key"])
		if index != -1 {
			json.NewEncoder(w).Encode(WrapObject("Found", pairs[index]))

			pair.Last = "updated"
			json.NewEncoder(w).Encode(WrapObject("Updated", pair))

			pairs[index] = pair
		} else {
			json.NewEncoder(w).Encode(WrapObject("Not Found", nil))

			pair.Last = "added"
			json.NewEncoder(w).Encode(WrapObject("Added", pair))

			pairs = append(pairs, pair)
		}
	*/
	PutPairSimpleDB(pair.User, pair.Key, pair.Value, pair.Time.String(), pair.Server, w)

	GetPairSimpleDB(pair.User, pair.Key, w)
}

func DeletePair(username string, w http.ResponseWriter, r *http.Request) {

	EchoParseRequest(w, r)

	params := mux.Vars(r)

	/*
		// For in memory database
		index := FindPair(username, params["key"])
		if index != -1 {
			json.NewEncoder(w).Encode(WrapObject("Found", pairs[index]))

			pairs[index].Last = "deleting"
			pairs[index].Time = time.Now()
			pairs[index].Server = serverName.ServerName
			json.NewEncoder(w).Encode(WrapObject("Deleting", pairs[index]))

			pairs = append(pairs[:index], pairs[index+1:]...)
			SendUserPairs(username, w)
		} else {
			json.NewEncoder(w).Encode(WrapObject("Not Found", nil))
		}
	*/
	DeletePairSimpleDB(username, params["key"], w)

	GetPairSimpleDB(username, params["key"], w)
}

////
type ServerName struct {
	ServerName string `json:"ServerName"`
}

var serverName ServerName

func InitServerName() {

	serverName.ServerName, _ = os.Hostname()
	log.Printf("%+v\n", serverName)
}

func EchoParseRequest(w http.ResponseWriter, r *http.Request) Pair {

	json.NewEncoder(w).Encode(serverName)

	params := mux.Vars(r)
	json.NewEncoder(w).Encode(WrapObject("Request params", params))

	var body Pair
	json.NewDecoder(r.Body).Decode(&body)
	json.NewEncoder(w).Encode(WrapObject("Request body", body))

	return body
}

func HandlePing(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(serverName)
}

////
type UserHttpHandlerFunc func(username string, w http.ResponseWriter, r *http.Request)

func Authenticate(h UserHttpHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		username, password, authOK := r.BasicAuth()
		if authOK == false {
			log.Println("BasicAuth failed")
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		if ValidateUser(username, password) != "" {
			h(username, w, r)
		} else {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
		}
	}
}

////
func Logger(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		log.Printf("Received %s\t%s", r.Method, r.RequestURI)
		h.ServeHTTP(w, r)
		log.Printf("Completed %s\t%s\t%s", r.Method, r.RequestURI, time.Since(start))
	}
}

////
func main() {

	InitServerName()
	InitSimpleDB()
	InitUsers()
	InitPairs()

	/*
		// Try it out
		for _, item := range users {
			SelectPairsSimpleDB(item.Username, nil)
		}
	*/

	router := mux.NewRouter()

	router.HandleFunc("/", Logger(HandlePing)).Methods("GET")
	router.HandleFunc("/pairs", Logger(Authenticate(GetPairs))).Methods("GET")
	router.HandleFunc("/pair/{key}", Logger(Authenticate(GetPair))).Methods("GET")
	router.HandleFunc("/pair/{key}", Logger(Authenticate(AddUpdatePair))).Methods("POST")
	router.HandleFunc("/pair/{key}", Logger(Authenticate(DeletePair))).Methods("DELETE")

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
