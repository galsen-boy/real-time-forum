package database

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"real-time-forum/internal/structure"
)

// Attempts to insert a new post into the database
func NewPost(path string, p structure.Post, u structure.User) error {
	//Opens the database
	db, err := OpenDB(path)
	if err != nil {
		return err
	}

	defer db.Close()

	dt := time.Now().Format("01-02-2006 15:04:05")

	//Executes the insert statement
	_, err = db.Exec(AddPost, u.Id, p.Category, p.Title, p.Content, dt, p.Likes, p.Dislikes)
	if err != nil {
		return err
	}

	return nil
}

// Converts post table query results into an array of post structs
func ConvertRowToPost(rows *sql.Rows) ([]structure.Post, error) {
	var posts []structure.Post

	//Loops through the rows provided
	for rows.Next() {
		var p structure.Post

		//Stores the row data in a temporary post struct
		err := rows.Scan(&p.Id, &p.User_id, &p.Category, &p.Title, &p.Content, &p.Date, &p.Likes, &p.Dislikes)
		if err != nil {
			break
		}

		//Appends the temporary struct to the array
		posts = append(posts, p)
	}

	return posts, nil
}

// Gets all posts from the database
func FindAllPosts(path string) ([]structure.Post, error) {
	//Opens the database
	db, err := OpenDB(path)
	if err != nil {
		return []structure.Post{}, errors.New("failed to open database")
	}

	defer db.Close()

	//Finds all the users
	rows, err := db.Query(GetAllPost)
	if err != nil {
		return []structure.Post{}, errors.New("failed to find posts")
	}

	//Convert the rows to an array of users
	posts, err := ConvertRowToPost(rows)
	if err != nil {
		return []structure.Post{}, errors.New("failed to convert")
	}

	return posts, nil
}

// Gets posts from the database based on the passed parameter (id, user_id, category)
func FindPostByParam(path, parameter, data string) ([]structure.Post, error) {
	var q *sql.Rows

	//Opens the database
	db, err := OpenDB(path)
	if err != nil {
		return []structure.Post{}, errors.New("failed to open database")
	}

	defer db.Close()

	//Checks which parameter to search the database by
	switch parameter {
	case "id":
		//Converts the data to an integer
		i, err := strconv.Atoi(data)
		if err != nil {
			return []structure.Post{}, errors.New("id must be an integer")
		}

		//Searches the database by id
		q, err = db.Query(GetPostById, i)
		if err != nil {
			return []structure.Post{}, errors.New("could not find id")
		}
	case "user_id":
		//Searches the database by user_id
		q, err = db.Query(GetAllPostByUser, data)
		if err != nil {
			return []structure.Post{}, errors.New("could not find any posts by that user")
		}
	case "category":
		//Searches the database by category
		q, err = db.Query(GetAllPostByCategory, data)
		if err != nil {
			return []structure.Post{}, errors.New("could not find any posts with that category")
		}
	default:
		//Returns an error if searched by a different parameter
		return []structure.Post{}, errors.New("cannot search by that parameter")
	}

	//Converts the database rows to an array of post structs
	posts, err := ConvertRowToPost(q)
	if err != nil {
		return []structure.Post{}, errors.New("failed to convert")
	}

	return posts, nil
}
