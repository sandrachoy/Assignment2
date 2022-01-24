package main

import (
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
	PostDate       string `json: "PostDate"`
	PostTime       string `json: "PostTime"`
}

type Reply struct {
	ReplyID        string `json:"ReplyID"`
	PosterName     string `json: "PosterName"`
	PostDate       string `json: "PostDate"`
	PostTime       string `json: "PostTime"`
	ConversationID string `json:"ConversationID"`
}

// used for storing conversations on the REST API
var conversations map[string]Conversation

func allConversations(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "List of all Conversations")
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

func createConversation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Header.Get("Content-type") == "application/json" {

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
	}

}

func updateConversation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Header.Get("Content-type") == "application/json" {
		//---PUT is for creating or updating
		// existing conversation---

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
	}
}

func deleteConversation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Header.Get("Content-type") == "application/json" {

		if _, ok := conversations[params["conversationID"]]; ok {
			delete(conversations, params["conversationID"])
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Conversation deleted: " +
				params["conversationID"]))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No conversation found"))
		}

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

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
