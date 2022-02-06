package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Conversation struct {
	ConversationID string `json:"ConversationID"`
	PosterName     string `json: "PosterName"`
	PostContent    string `json: "PostContent"`
	PostDate       string `json: "PostDate"`
	PostTime       string `json: "PostTime"`
}

type Reply struct {
	ReplyID        int    `json:"ReplyID"`
	PosterName     string `json: "PosterName"`
	PostContent    string `json: "PostContent"`
	PostDate       string `json: "PostDate"`
	PostTime       string `json: "PostTime"`
	ConversationID string `json:"ConversationID"`
}

// used for storing conversations on the REST API
var conversations map[string]Conversation

func getConversationDB(db *sql.DB) {
	results, err := db.Query("Select * FROM activities_groups.Conversations")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		// map this type to the record in the table
		var conversation Conversation
		err = results.Scan(&conversation.ConversationID, &conversation.PosterName, &conversation.PostContent,
			&conversation.PostDate, &conversation.PostTime)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println(conversation.ConversationID, conversation.PosterName, conversation.PostContent,
			conversation.PostDate, conversation.PostTime)
	}
}

func allConversations(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "List of all Conversations")
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/activities_groups")

	if err != nil {
		panic(err.Error())
	}

	getConversationDB(db)

	defer db.Close()

	kv := r.URL.Query()

	for k, v := range kv {
		fmt.Println(k, v)
	}
	// returns all the conversations in JSON
	json.NewEncoder(w).Encode(conversations)
}

func conversation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Fprintf(w, "Detail for conversation "+params["conversationID"])
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, r.Method)

	if r.Method == "GET" {
		if _, ok := conversations[params["conversationID"]]; ok {
			json.NewEncoder(w).Encode(
				conversations[params["conversationID"]])
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No conversation found"))
		}
	}

}

func createConversationDB(db *sql.DB, ID string, PN string, PC string, PD string, PT string) {
	query := fmt.Sprintf("INSERT INTO activities_groups.Conversations VALUES ('%s', '%s', '%s','%s', '%s')",
		ID, PN, PC, PD, PT)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func createConversation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Header.Get("Content-type") == "application/json" {

		db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/activities_groups")

		if err != nil {
			panic(err.Error())
		}

		// POST is for creating new conversation

		// read the string sent to the service
		var newConversation Conversation
		reqBody, err := ioutil.ReadAll(r.Body)

		if err == nil {
			// convert JSON to object
			json.Unmarshal(reqBody, &newConversation)

			if newConversation.ConversationID == "" {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte(
					"422 - Please supply conversation " +
						"information " + "in JSON format"))
				return
			}

			// check if conversation exists; add only if
			// conversation does not exist
			if _, ok := conversations[params["conversationID"]]; !ok {
				createConversationDB(db, newConversation.ConversationID, newConversation.PosterName, newConversation.PostContent, newConversation.PostDate, newConversation.PostTime)
				conversations[params["conversationID"]] = newConversation
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("201 - Conversation added: " +
					params["conversationID"]))
			} else {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(
					"409 - Duplicate conversation ID"))
			}
		} else {
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply conversation information " +
				"in JSON format"))
		}

		defer db.Close()
	}

}

func updateConversationDB(db *sql.DB, ID string, PN string, PC string, PD string, PT string) {
	query := fmt.Sprintf(
		"UPDATE activities_groups.Conversations SET PosterName='%s', PostContent='%s' PostDate='%s', PostTime='%s' WHERE ConversationID='%s'",
		PN, PC, PD, PT, ID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

}

func updateConversation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Header.Get("Content-type") == "application/json" {
		//---PUT is for creating or updating
		// existing conversation---

		db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/activities_groups")

		if err != nil {
			panic(err.Error())
		}

		var newConversation Conversation
		reqBody, err := ioutil.ReadAll(r.Body)

		if err == nil {
			json.Unmarshal(reqBody, &newConversation)

			if newConversation.ConversationID == "" {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte(
					"422 - Please supply conversation " +
						" information " +
						"in JSON format"))
				return
			}

			// check if conversation exists; add only if
			// conversation does not exist
			if _, ok := conversations[params["conversationID"]]; !ok {
				conversations[params["conversationID"]] =
					newConversation
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("201 - Conversation added: " +
					params["conversationID"]))
				updateConversationDB(db, newConversation.ConversationID, newConversation.PosterName, newConversation.PostContent, newConversation.PostDate, newConversation.PostTime)
			} else {
				// update conversation
				conversations[params["conversationID"]] = newConversation
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("202 - Conversation updated: " +
					params["conversationID"]))
			}
		} else {
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply " +
				"conversation information " +
				"in JSON format"))

		}
		defer db.Close()
	}
}

func deleteConversationDB(db *sql.DB, ID string) {
	query := fmt.Sprintf(
		"DELETE FROM activities_groups.Conversations WHERE ConversationID='%s'", ID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

}

func deleteConversation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Header.Get("Content-type") == "application/json" {

		db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/activities_groups")

		if err != nil {
			panic(err.Error())
		}

		if _, ok := conversations[params["conversationID"]]; ok {
			delete(conversations, params["conversationID"])
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Conversation deleted: " +
				params["conversationID"]))
			deleteConversationDB(db, params["conversationID"])
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No conversation found"))
		}

		defer db.Close()
	}
}

func main() {
	// instantiate conversations
	conversations = make(map[string]Conversation)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/groups/conversations", allConversations)
	router.HandleFunc("/api/v1/groups/conversations/{conversationID}", conversation).Methods("GET")
	router.HandleFunc("/api/v1/groups/addConversation", createConversation).Methods("POST")
	router.HandleFunc("/api/v1/groups/updateConversation", updateConversation).Methods("PUT")
	router.HandleFunc("/api/v1/groups/deleteConversation", deleteConversation).Methods("DELETE")

	fmt.Println("Listening at port 8191")
	log.Fatal(http.ListenAndServe(":8191", router))
}
