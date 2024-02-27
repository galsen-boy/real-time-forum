package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"real-time-forum/internal/config"
	"real-time-forum/internal/database"
	"real-time-forum/internal/structure"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	//Prevents the endpoint being called by other url paths
	if r.URL.Path != "/comment" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	//Checks whether it is a POST or GET request
	switch r.Method {
	case "GET":
		//Checks for a passed search parameter and data
		param := r.URL.Query().Get("param")
		data := r.URL.Query().Get("data")
		if param == "" || data == "" {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}

		//Finds the comments based on the search parameter and data
		comments, err := database.FindCommentByParam(config.Path, param, data)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		//Marshals the array of comment structs to a json object
		resp, err := json.Marshal(comments)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		//Writes the json object to the frontend
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	case "POST":
		//Stores the unmarshalled register data
		var newComment structure.Comment

		//Decodes the request body into the post struct
		//Returns a bad request if there's an error
		err := json.NewDecoder(r.Body).Decode(&newComment)
		if err != nil {
			http.Error(w, "400 bad request.", http.StatusBadRequest)
			return
		}

		fmt.Println(newComment)

		//Attemps to add the new post to the database
		err = database.NewComment(config.Path, newComment)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		var msg = structure.Resp{Msg: "Sent comment"}

		resp, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	default:
		//Prevents the use of other request types
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
