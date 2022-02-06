package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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
var replies map[string]Reply

func getReplyDB(db *sql.DB) {
	results, err := db.Query("Select * FROM activities_groups.Replies")

	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		// map this type to the record in the table
		var reply Reply
		err = results.Scan(&reply.ReplyID, &reply.PosterName, &reply.PostContent,
			&reply.PostDate, &reply.PostTime, &reply.ConversationID)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println(reply.ReplyID, reply.PosterName, reply.PostContent,
			reply.PostDate, reply.PostTime, reply.ConversationID)
	}
}

func allReplies(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "List of all Replies")
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/activities_groups")

	if err != nil {
		panic(err.Error())
	}

	getReplyDB(db)

	defer db.Close()

	kv := r.URL.Query()

	for k, v := range kv {
		fmt.Println(k, v)
	}
	// returns all the replies in JSON
	json.NewEncoder(w).Encode(replies)
}

func reply(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Fprintf(w, "Detail for reply "+params["replyID"])
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, r.Method)

	if r.Method == "GET" {
		if _, ok := replies[params["replyID"]]; ok {
			json.NewEncoder(w).Encode(
				replies[params["replyID"]])
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No reply found"))
		}
	}

}

func createReplyDB(db *sql.DB, ID int, PN string, PC string, PD string, PT string, CID string) {
	query := fmt.Sprintf("INSERT INTO activities_groups.Replies VALUES (%d, '%s', '%s','%s', '%s', %s)",
		ID, PN, PC, PD, PT, CID)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func createReply(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Header.Get("Content-type") == "application/json" {

		db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/activities_groups")

		if err != nil {
			panic(err.Error())
		}

		// POST is for creating new reply

		// read the string sent to the service
		var newReply Reply
		reqBody, err := ioutil.ReadAll(r.Body)

		if err == nil {
			// convert JSON to object
			json.Unmarshal(reqBody, &newReply)

			if newReply.ReplyID == 0 {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte(
					"422 - Please supply reply " +
						"information " + "in JSON format"))
				return
			}

			// check if reply exists; add only if
			// reply does not exist
			if _, ok := replies[params["replyID"]]; !ok {
				createReplyDB(db, newReply.ReplyID, newReply.PosterName, newReply.PostContent, newReply.PostDate, newReply.PostTime, newReply.ConversationID)
				replies[params["replyID"]] = newReply
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("201 - Reply added: " +
					params["replyID"]))
			} else {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(
					"409 - Duplicate reply ID"))
			}
		} else {
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply reply information " +
				"in JSON format"))
		}

		defer db.Close()
	}

}

func updateReplyDB(db *sql.DB, ID int, PN string, PC string, PD string, PT string, CID string) {
	query := fmt.Sprintf(
		"UPDATE activities_groups.Replies SET PosterName='%s', PostContent='%s', PostDate='%s', PostTime='%s', ConversationID='%s' WHERE ReplyID=%d",
		PN, PC, PD, PT, CID, ID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

}

func updateReply(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Header.Get("Content-type") == "application/json" {
		//---PUT is for creating or updating
		// existing reply---

		db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/activities_groups")

		if err != nil {
			panic(err.Error())
		}

		var newReply Reply
		reqBody, err := ioutil.ReadAll(r.Body)

		if err == nil {
			json.Unmarshal(reqBody, &newReply)

			if newReply.ReplyID == 0 {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte(
					"422 - Please supply reply " +
						" information " +
						"in JSON format"))
				return
			}

			// check if reply exists; add only if
			// reply does not exist
			if _, ok := replies[params["replyID"]]; !ok {
				replies[params["replyID"]] =
					newReply
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("201 - Reply added: " +
					params["replyID"]))
				updateReplyDB(db, newReply.ReplyID, newReply.PosterName, newReply.PostContent, newReply.PostDate, newReply.PostTime, newReply.ConversationID)
			} else {
				// update reply
				replies[params["replyID"]] = newReply
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("202 - Reply updated: " +
					params["replyID"]))
			}
		} else {
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply " +
				"reply information " +
				"in JSON format"))

		}
		defer db.Close()
	}
}

func deleteReplyDB(db *sql.DB, ID int) {
	query := fmt.Sprintf(
		"DELETE FROM activities_groups.Replies WHERE ReplyID=%d", ID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

}

func deleteReply(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Header.Get("Content-type") == "application/json" {

		db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/activities_groups")

		if err != nil {
			panic(err.Error())
		}

		if _, ok := replies[params["replyID"]]; ok {
			delete(replies, params["replyID"])

			replyID, err := strconv.Atoi(params["replyID"])

			if err != nil {
				panic(err.Error())
			}

			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Reply deleted: " +
				params["replyID"]))
			deleteReplyDB(db, replyID)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No reply found"))
		}

		defer db.Close()
	}
}

func main() {
	// instantiate conversations
	replies = make(map[string]Reply)

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/groups/replies", allReplies)
	router.HandleFunc("/api/v1/groups/replies/{replyID}", reply).Methods("GET")
	router.HandleFunc("/api/v1/groups/addReply", createReply).Methods("POST")
	router.HandleFunc("/api/v1/groups/updateReply", updateReply).Methods("PUT")
	router.HandleFunc("/api/v1/groups/deleteReply", deleteReply).Methods("DELETE")

	fmt.Println("Listening at port 8191")
	log.Fatal(http.ListenAndServe(":8191", router))
}
