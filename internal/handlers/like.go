package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"real-time-forum/internal/config"
	"real-time-forum/internal/database"
	"real-time-forum/internal/structure"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	//Prevents the endpoint being called by other url paths
	if r.URL.Path != "/like" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	//Grabs the post_id and table column from the url
	pid := r.URL.Query().Get("post_id")
	col := r.URL.Query().Get("col")
	if pid == "" || col == "" {
		http.Error(w, "400 bad request", http.StatusBadRequest)
		return
	}

	//Checks whether it is a POST or GET request
	switch r.Method {
	case "GET":
		//Gets all users who liked the post from the database
		users, err := database.PostLikedBy(config.Path, pid, col)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		//Marshals the user structs to a json object
		resp, err := json.Marshal(users)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		//Writes the json object to the frontend
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	case "POST":
		//Grabs the session cookie
		c, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		//Finds the currenlty logged in user
		curr, err := database.CurrentUser(config.Path, c.Value)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		var users []structure.User

		//Finds all users who liked/disliked the post
		users, err = database.PostLikedBy(config.Path, pid, col)
		if err != nil {
			users = []structure.User{}
		}

		var other string
		var removecurr = false
		var removeother = false

		//Checks whether the user has already liked/disliked the post
		for _, u := range users {
			if u.Id == curr.Id {
				removecurr = true
				break
			}
		}

		//Finds the other column
		if col == "likes" {
			other = "dislikes"
		} else if col == "dislikes" {
			other = "likes"
		} else {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}

		var otherCol []structure.User

		//Finds the users who liked/disliked the post (other column)
		otherCol, err = database.PostLikedBy(config.Path, pid, other)
		if err != nil {
			otherCol = []structure.User{}
		}

		//Checks whether the user has already liked/disliked the post (other column)
		for _, u := range otherCol {
			if u.Id == curr.Id && !removecurr {
				removeother = true
				break
			}
		}

		//Converts the users id to a string
		uid := strconv.Itoa(curr.Id)

		//If the user liked/disliked (other column) the post, remove it from the database
		if removeother {
			err = database.UpdateLikeDislike(config.Path, pid, uid, other, -1)
			if err != nil {
				http.Error(w, "500 internal error", http.StatusInternalServerError)
				return
			}
		}

		//If the user already liked/disliked (selected column) the post, remove it from the database
		//Otherwise add it to the database
		if removecurr {
			err = database.UpdateLikeDislike(config.Path, pid, uid, col, -1)
			if err != nil {
				http.Error(w, "500 internal error", http.StatusInternalServerError)
				return
			}
		} else {
			err = database.UpdateLikeDislike(config.Path, pid, uid, col, 1)
			if err != nil {
				http.Error(w, "500 internal error", http.StatusInternalServerError)
				return
			}
		}

		currPost, err := database.FindPostByParam(config.Path, "id", pid)
		if err != nil {
			http.Error(w, "500 internal error", http.StatusInternalServerError)
			return
		}

		likes := strconv.Itoa(currPost[0].Likes)
		dislikes := strconv.Itoa(currPost[0].Dislikes)

		var msg = structure.Resp{Msg: likes + "|" + dislikes}

		resp, err := json.Marshal(msg)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	default:
		//Prevents other methods from being used
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}
}
