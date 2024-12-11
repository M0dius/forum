package database

import (
	"database/sql"
	"log"
	"time"
)

type User struct {
	ID               int            `json:"id"`
	FirstName        string         `json:"first_name"`
	LastName         string         `json:"last_name"`
	Username         string         `json:"username"`
	Email            string         `json:"email"`
	Password         string         `json:"password"`
	SessionSessionID int            `json:"session_sessionid"`
	RoleID           int            `json:"role_id"`
	Avatar           sql.NullString `json:"avatar"` // Use sql.NullString to handle NULL values
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Comment struct {
	ID         int       // Comment ID
	Content    string    // Comment text
	CommentAt  time.Time // Timestamp of the comment
	PostPostID int       // Post ID the comment belongs to
	UserUserID int       // User ID who made the comment
}

type Post struct {
	PostID     int
	Image      string
	Content    string
	PostAt     string
	UserUserID int
	Username   string
	Avatar     sql.NullString // Use sql.NullString to handle NULL values
}

// GetAllCategories retrieves all categories from the database
func GetAllCategories(db *sql.DB) ([]Category, error) {
	rows, err := db.Query("SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func GetComment(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT *  FROM comment")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetComments(db *sql.DB) ([]Comment, error) {
	// Query to retrieve all comments
	rows, err := db.Query("SELECT commentid, content, comment_at, post_postid, user_userid FROM comment")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		var commentAt time.Time // SQLite DATETIME is fetched as a string

		// Scan each row into the Comment struct
		if err := rows.Scan(&comment.ID, &comment.Content, &commentAt, &comment.PostPostID, &comment.UserUserID); err != nil {
			return nil, err
		}

		// Parse comment_at into a time.Time object
		comment.CommentAt = commentAt

		comments = append(comments, comment)
	}

	// Check for errors after iterating
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func GetAllPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query(`
        SELECT post.postid, post.image, post.content, post.post_at, post.user_userid, user.Username, user.Avatar
        FROM post
        JOIN user ON post.user_userid = user.userid
    `)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var avatar sql.NullString
		if err := rows.Scan(&post.PostID, &post.Image, &post.Content, &post.PostAt, &post.UserUserID, &post.Username, &avatar); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		post.Avatar = avatar
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error in rows:", err)
		return nil, err
	}

	return posts, nil
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT userid, F_name, L_name, Username, Email, Avatar FROM user")
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var avatar sql.NullString
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &avatar); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		user.Avatar = avatar
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error in rows:", err)
		return nil, err
	}

	return users, nil
}
