package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
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

const conversationURL = "http://localhost:5000/api/v1/groups/conversations"
const replyURL = "http://localhost:5000/api/v1/groups/replies"
const key = "2c78afaf-97da-4816-bbee-9ad239abb296"

var userinput = ""

func mainMenu() {
	fmt.Println("---Main Menu---")
	fmt.Println("1. Create Conversation")
	fmt.Println("2. Create Reply")
	fmt.Println("3. Your Conversations")
	fmt.Println("4. Your Replies")
	fmt.Println("5. View all Conversations")
	fmt.Println("6. View all Replies")
	fmt.Println("0. Exit Menu")
	fmt.Scanln(&userinput)

	switch userinput {
	case "1":
		createConversation()
	case "2":
		createReply()
	case "3":
		yourConversations()
	case "4":
		yourReplies()
	case "5":
		viewAllConversations()
	case "6":
		viewAllReplies()
	case "0":
		break
	default:
		fmt.Println("Please enter a valid number")
		mainMenu()
	}
}

func createConversation() {
	var conversationID, posterName, postContent string

	fmt.Println("---Create Conversation---")
	fmt.Println("Enter Conversation ID:")
	fmt.Println(&conversationID)
	fmt.Println("Enter Your Name:")
	fmt.Println(&posterName)
	fmt.Println("Enter Post Content:")
	fmt.Println(&postContent)

	addConversation(conversationID, posterName, postContent)
}

func createReply() {
	replyID := 1
	var posterName, conversationID, postContent string

	fmt.Println("---Create Reply---")
	fmt.Println("Enter Conversation ID:")
	fmt.Println(&conversationID)
	fmt.Println("Enter Your Name:")
	fmt.Println(&posterName)
	fmt.Println("Enter Reply Content:")
	fmt.Println(&postContent)

	addReply(replyID, conversationID, posterName, postContent)

	replyID += 1
}

func yourConversations() {
	fmt.Println("---Your Conversations---")
	fmt.Println("1. View Conversation")
	fmt.Println("2. Update Conversation")
	fmt.Println("3. Delete Conversation")
	fmt.Println("0. Back to Main Menu")
	fmt.Scanln(&userinput)

	switch userinput {
	case "1":
		viewConversation()
	case "2":
		updateConversation()
	case "3":
		deleteConversation()
	case "0":
		mainMenu()
	default:
		fmt.Println("Please enter a valid number")
		mainMenu()
	}
}

func viewConversation() {
	var posterName string

	fmt.Println("---View Conversation---")
	fmt.Println("Enter Your Conversation ID:")
	fmt.Println(&posterName)

	getConversations(posterName)
	yourConversations()
}

func updateConversation() {
	var conversationID, posterName, postContent string

	fmt.Println("---Update Conversation---")
	fmt.Println("Enter Conversation ID:")
	fmt.Println(&conversationID)
	fmt.Println("Enter Your Name:")
	fmt.Println(&posterName)
	fmt.Println("Enter Post Content (Will replace old content):")
	fmt.Println(&postContent)

	editConversation(conversationID, posterName, postContent)
	yourConversations()
}

func deleteConversation() {
	var conversationID string

	fmt.Println("---Delete Conversation---")
	fmt.Println("Enter Conversation ID:")
	fmt.Println(&conversationID)

	removeConversation(conversationID)
	yourConversations()
}

func yourReplies() {
	fmt.Println("---Your Replies---")
	fmt.Println("1. View Reply")
	fmt.Println("2. Update Reply")
	fmt.Println("3. Delete Reply")
	fmt.Println("0. Back to Main Menu")
	fmt.Scanln(&userinput)

	switch userinput {
	case "1":
		viewReply()
	case "2":
		updateReply()
	case "3":
		deleteReply()
	case "0":
		mainMenu()
	default:
		fmt.Println("Please enter a valid number")
		mainMenu()
	}
}

func viewReply() {
	var replyID int

	fmt.Println("---View Reply---")
	fmt.Println("Enter Reply ID:")
	fmt.Println(&replyID)

	getReplies(strconv.Itoa(replyID))
	yourReplies()
}

func updateReply() {
	var replyID int
	var posterName, postContent string

	fmt.Println("---Update Reply---")
	fmt.Println("Enter Reply ID:")
	fmt.Println(&replyID)
	fmt.Println("Enter Your Name:")
	fmt.Println(&posterName)
	fmt.Println("Enter Post Content (Will replace old content):")
	fmt.Println(&postContent)

	editReply(replyID, posterName, postContent)
	yourReplies()
}

func deleteReply() {
	var replyID int

	fmt.Println("---Delete Reply---")
	fmt.Println("Enter Reply ID:")
	fmt.Println(&replyID)

	removeReply(replyID)
	yourReplies()
}

func viewAllConversations() {
	url := conversationURL

	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		var conversation Conversation

		data, _ := ioutil.ReadAll(response.Body)

		// parses the JSON-encoded data and stores the result in the value pointed to by the second argument
		err := json.Unmarshal([]byte(data), &conversation)

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		fmt.Println(response.StatusCode)
		fmt.Println("Conversation ID:", conversation.ConversationID)
		fmt.Println("Poster Name:", conversation.PosterName)
		fmt.Println("Post Content:", conversation.PostContent)
		fmt.Println("Post Date:", conversation.PostDate)
		fmt.Println("Post Time:", conversation.PostTime)

		response.Body.Close()
	}
}

func viewAllReplies() {
	url := replyURL

	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		var reply Reply

		data, _ := ioutil.ReadAll(response.Body)

		// parses the JSON-encoded data and stores the result in the value pointed to by the second argument
		err := json.Unmarshal([]byte(data), &reply)

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		fmt.Println(response.StatusCode)
		fmt.Println("Reply ID:", reply.ReplyID)
		fmt.Println("Poster Name:", reply.PosterName)
		fmt.Println("Post Content:", reply.PostContent)
		fmt.Println("Post Date:", reply.PostDate)
		fmt.Println("Post Time:", reply.PostTime)
		fmt.Println("Conversation ID:", reply.ConversationID)

		response.Body.Close()
	}
}

func addConversation(conversationID string, posterName string, postContent string) {
	conversation := Conversation{
		ConversationID: conversationID,
		PosterName:     posterName,
		PostContent:    postContent,
		PostDate:       "",
		PostTime:       "",
	}

	jsonValue, _ := json.Marshal(conversation)

	response, err := http.Post(conversationURL+"/"+conversationID,
		"application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}

}

func addReply(replyID int, conversationID string, posterName string, postContent string) {
	reply := Reply{
		ReplyID:        replyID,
		PosterName:     posterName,
		PostContent:    postContent,
		PostDate:       "",
		PostTime:       "",
		ConversationID: conversationID,
	}

	jsonValue, _ := json.Marshal(reply)

	response, err := http.Post(replyURL+"/"+strconv.Itoa(replyID),
		"application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

func getConversations(code string) {
	url := conversationURL
	if code != "" {
		url = conversationURL + "/" + code
	}

	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		var conversation Conversation

		data, _ := ioutil.ReadAll(response.Body)

		// parses the JSON-encoded data and stores the result in the value pointed to by the second argument
		err := json.Unmarshal([]byte(data), &conversation)

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		fmt.Println(response.StatusCode)
		fmt.Println("Conversation ID:", conversation.ConversationID)
		fmt.Println("Poster Name:", conversation.PosterName)
		fmt.Println("Post Content:", conversation.PostContent)
		fmt.Println("Post Date:", conversation.PostDate)
		fmt.Println("Post Time:", conversation.PostTime)

		response.Body.Close()
	}
}

func getReplies(code string) {
	url := replyURL
	if code != "" {
		url = replyURL + "/" + code
	}

	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		var reply Reply

		data, _ := ioutil.ReadAll(response.Body)

		// parses the JSON-encoded data and stores the result in the value pointed to by the second argument
		err := json.Unmarshal([]byte(data), &reply)

		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}

		fmt.Println(response.StatusCode)
		fmt.Println("Reply ID:", reply.ReplyID)
		fmt.Println("Poster Name:", reply.PosterName)
		fmt.Println("Post Content:", reply.PostContent)
		fmt.Println("Post Date:", reply.PostDate)
		fmt.Println("Post Time:", reply.PostTime)
		fmt.Println("Conversation ID:", reply.ConversationID)

		response.Body.Close()
	}
}

func editConversation(conversationID string, posterName string, postContent string) {
	conversation := Conversation{
		ConversationID: conversationID,
		PosterName:     posterName,
		PostContent:    postContent,
	}

	jsonValue, _ := json.Marshal(conversation)

	request, err := http.NewRequest(http.MethodPut,
		conversationURL+"/"+conversationID,
		bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}

}

func editReply(replyID int, posterName string, postContent string) {
	reply := Reply{
		ReplyID:     replyID,
		PosterName:  posterName,
		PostContent: postContent,
	}

	jsonValue, _ := json.Marshal(reply)

	request, err := http.NewRequest(http.MethodPut,
		replyURL+"/"+strconv.Itoa(replyID),
		bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
}

func removeReply(replyID int) {
	request, err := http.NewRequest(http.MethodDelete,
		replyURL+"/"+strconv.Itoa(replyID), nil)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}

}

func removeConversation(conversationID string) {
	request, err := http.NewRequest(http.MethodDelete,
		conversationURL+"/"+conversationID, nil)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}

}

func main() {
	mainMenu()
}
