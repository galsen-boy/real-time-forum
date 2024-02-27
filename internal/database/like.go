package database

import (
	"database/sql"
	"errors"
	"strconv"

	"real-time-forum/internal/structure"
)

// Updates liked_posts or disliked_posts table
func UpdateLikedPosts(col string, pid, uid, i int, db *sql.DB) error {
	//Checks which column to update
	switch col {
	case "likes":
		//Checks whether it needs to be added or removed and executes the sql statment
		if i == 1 {
			_, err := db.Exec(AddLike, pid, uid)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Exec(RemoveLike, pid, uid)
			if err != nil {
				return err
			}
		}
	case "dislikes":
		//Checks whether it needs to be added or removed and executes the sql statment
		if i == 1 {
			_, err := db.Exec(AddDislike, pid, uid)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Exec(RemoveDislike, pid, uid)
			if err != nil {
				return err
			}
		}
	default:
		//Returns an error if a different column is passed
		return errors.New("cannot edit that column")
	}

	return nil
}

// Adds like to database
func UpdateLikeDislike(path, post_id, user_id, col string, i int) error {
	//Opens the database
	db, err := OpenDB(path)
	if err != nil {
		return err
	}

	defer db.Close()

	//Checks whether i is 1 or -1
	if i != 1 && i != -1 {
		return errors.New("can only change by 1 or -1")
	}

	//Finds the post
	temp, err := FindPostByParam(path, "id", post_id)
	if err != nil {
		return err
	}

	p := temp[0]

	//Converts post_id and user_id to integers
	pid, err := strconv.Atoi(post_id)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(user_id)
	if err != nil {
		return err
	}

	//Checks which column to update
	switch col {
	case "likes":
		//Updates the like count by +1 or -1 and updates the post table
		count := p.Likes + i

		_, err := db.Exec(UpdateLike, count, pid)
		if err != nil {
			return err
		}
	case "dislikes":
		//Increases the dislike count by 1 and updates the post table
		count := p.Dislikes + i

		_, err := db.Exec(UpdateDislike, count, pid)
		if err != nil {
			return err
		}
	default:
		//Returns an error when trying to update a different column
		return errors.New("can only update likes and dislikes")
	}

	//Updates the liked_posts table
	UpdateLikedPosts(col, pid, uid, i, db)

	return nil
}

// Finds all users who liked or disliked a post
func PostLikedBy(path, post_id, col string) ([]structure.User, error) {
	var q *sql.Rows

	//Opens the database
	db, err := OpenDB(path)
	if err != nil {
		return []structure.User{}, err
	}

	defer db.Close()

	//Converts post_id to an integer
	pid, err := strconv.Atoi(post_id)
	if err != nil {
		return []structure.User{}, err
	}

	switch col {
	case "likes":
		//Finds all users that liked a post through liked_posts table
		q, err = db.Query(GetPostLikes, pid)
		if err != nil {
			return []structure.User{}, err
		}
	case "dislikes":
		//Finds all users that disliked a post through disliked_posts table
		q, err = db.Query(GetPostDislikes, pid)
		if err != nil {
			return []structure.User{}, err
		}
	default:
		return []structure.User{}, errors.New("incorrect column")
	}

	//Converts the query results to an array of user structs
	users, err := ConvertRowToUser(q)
	if err != nil {
		return []structure.User{}, err
	}

	return users, nil
}

// Finds all posts liked or disliked by a user
func UserLiked(path, user_id, col string) ([]structure.Post, error) {
	var q *sql.Rows

	//Opens the database
	db, err := OpenDB(path)
	if err != nil {
		return []structure.Post{}, err
	}

	defer db.Close()

	//Converts post_id to an integer
	uid, err := strconv.Atoi(user_id)
	if err != nil {
		return []structure.Post{}, err
	}

	switch col {
	case "likes":
		//Finds all users that liked a post through liked_posts table
		q, err = db.Query(GetUserLikes, uid)
		if err != nil {
			return []structure.Post{}, err
		}
	case "dislikes":
		//Finds all users that disliked a post through disliked_posts table
		q, err = db.Query(GetUserDislikes, uid)
		if err != nil {
			return []structure.Post{}, err
		}
	default:
		return []structure.Post{}, errors.New("incorrect column")
	}

	//Converts the query results to an array of user structs
	posts, err := ConvertRowToPost(q)
	if err != nil {
		return []structure.Post{}, err
	}

	return posts, nil
}
