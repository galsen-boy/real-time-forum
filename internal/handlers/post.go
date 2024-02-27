package handlers

import (
	"encoding/json"
	"net/http"

	"real-time-forum/internal/config"
	"real-time-forum/internal/database"
	"real-time-forum/internal/structure"
)

// PostHandler handles the /post endpoint
func PostHandler(w http.ResponseWriter, r *http.Request) {
	//Prevents the endpoint being called by other url paths
	if r.URL.Path != "/post" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	//Checks whether it is a POST or GET request
	switch r.Method {
	case "GET":
		var posts []structure.Post
		var err error

		//Checks for a passed search parameter
		param := r.URL.Query().Get("param")
		if param == "" {
			//If not found, returns all users
			posts, err = database.FindAllPosts(config.Path)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			//If found, checks for the search data
			data := r.URL.Query().Get("data")
			//If no search data is found, returns an error
			if data == "" {
				http.Error(w, "400 bad request", http.StatusBadRequest)
				return
			}

			//Finds the posts based on the parameter and data
			posts, err = database.FindPostByParam(config.Path, param, data)
			if err != nil {
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}
		}

		//Marshals the array of post structs to a json object
		resp, err := json.Marshal(posts)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		//Writes the json object to the frontend
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	case "POST":
		//Stores the unmarshalled register data
		var newPost structure.Post

		//Decodes the request body into the post struct
		//Returns a bad request if there's an error
		err := json.NewDecoder(r.Body).Decode(&newPost)
		if err != nil {
			http.Error(w, "400 bad request.", http.StatusBadRequest)
			return
		}
		//Checks whether the post is empty
		cookie, err := r.Cookie("session")
		if err != nil {
			return
		}

		foundVal := cookie.Value
		//Checks whether the user is logged in
		curr, err := database.CurrentUser(config.Path, foundVal)
		if err != nil {
			return
		}

		//Attemps to add the new post to the database
		err = database.NewPost(config.Path, newPost, curr)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		//Sends a message back if successfully posted
		var msg = structure.Resp{Msg: "New post added"}
		//Marshals the message to a json object
		resp, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		//Writes the json object to the frontend
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	default:
		//Prevents the use of other request types
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
