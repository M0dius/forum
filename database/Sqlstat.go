package database

import (
	"database/sql"
	"fmt"
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
	SessionSessionID int            `json:"current_session"`
	RoleID           int            `json:"role_id"`
	Avatar           sql.NullString `json:"avatar"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Comment struct {
	ID        int
	PostID    int
	UserID    int
	FirstName string
	LastName  string
	Username  string
	Content   string
	CreatedAt time.Time
	Avatar    sql.NullString // Add this field
	Likes     int
	Dislikes  int
}

type Post struct {
	PostID     int
	Image      sql.NullString
	Content    string
	PostAt     time.Time
	UserUserID int
	Username   string
	FirstName  string
	LastName   string
	Avatar     sql.NullString
	Likes      int
	Dislikes   int
	Comments   int
	Categories []Category
}

type Notification struct {
	ID        int
	UserID    int
	PostID    int
	Message   string
	CreatedAt time.Time
	UserImage string
	UserName  string
}

type Report struct {
	ID           int
	PostID       int
	ReportedBy   int
	ReportReason string
	CreatedAt    time.Time
}

type UserLog struct {
	ID        int
	UserID    int
	Action    string
	Timestamp time.Time
}

type UserSession struct {
	ID     int
	UserID int
	Start  time.Time
	End    time.Time
}

func GetAllCategories(db *sql.DB) ([]Category, error) {
	rows, err := db.Query("SELECT idcategories, name, description FROM categories")
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func GetComments(db *sql.DB) ([]Comment, error) {
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
		if err := rows.Scan(&comment.ID, &comment.Content, &commentAt, &comment.PostID, &comment.UserID); err != nil {
			return nil, err
		}

		comment.CreatedAt = commentAt

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// func GetUserPosts1(db *sql.DB, username string) ([]Post, error) {
// 	rows, err := db.Query(`
//         SELECT postid, image, content, post_at
//         FROM post
//         JOIN user ON user_userid = userid
//         WHERE Username = ?
//         ORDER BY post.post_at DESC
//     `, username)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var posts []Post
// 	for rows.Next() {
// 		var post Post
// 		if err := rows.Scan(&post.PostID, &post.Image, &post.Content, &post.PostAt); err != nil {
// 			return nil, err
// 		}
// 		posts = append(posts, post)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return posts, nil
// }

func GetUserReaction(db *sql.DB, userid int, filter string) ([]Post, error) {

	Lpost, err := GetUserLikedPosts(db, userid)
	if err != nil {
		log.Println("Error fetching liked posts:", err)
	}
	Dpost, err := GetUserDislikedPosts(db, userid)
	if err != nil {
		log.Println("Error fetching disliked posts:", err)
	}

	if filter == "likes" {
		// Lpost = append(Lpost, Dpost...)
		return Lpost, nil
	} else {
		// Dpost = append(Dpost, Lpost...)
		return Dpost, nil
	}
}
func GetUserLikedPosts(db *sql.DB, userID int) ([]Post, error) {
	rows, err := db.Query(`
        SELECT post.postid, post.image, post.content, post.post_at,
		        		user.avatar, user.F_name, user.L_name, user.Username,
		 (SELECT COUNT(*) FROM likes WHERE likes.post_postid = post.postid) AS Likes,
               (SELECT COUNT(*) FROM dislikes WHERE dislikes.post_postid = post.postid) AS Dislikes,
               (SELECT COUNT(*) FROM comment WHERE comment.post_postid = post.postid) AS Comments
        FROM post
        JOIN likes ON post.postid = likes.post_postid
			JOIN user ON post.user_userid = user.userid 
        WHERE likes.user_userid = ?
        ORDER BY post.post_at DESC`, userID)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostID, &post.Image, &post.Content, &post.PostAt, &post.Avatar, &post.FirstName, &post.LastName, &post.Username, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error in rows:", err)
		return nil, err
	}

	return posts, nil
}

func GetUserDislikedPosts(db *sql.DB, userID int) ([]Post, error) {
	rows, err := db.Query(`
 SELECT post.postid, post.image, post.content, post.post_at,
        		user.avatar, user.F_name, user.L_name, user.Username,
		               (SELECT COUNT(*) FROM likes WHERE likes.post_postid = post.postid) AS Likes,
               (SELECT COUNT(*) FROM dislikes WHERE dislikes.post_postid = post.postid) AS Dislikes,
               (SELECT COUNT(*) FROM comment WHERE comment.post_postid = post.postid) AS Comments
        FROM post
        JOIN dislikes ON post.postid = dislikes.post_postid
			JOIN user ON post.user_userid = user.userid 
        WHERE dislikes.user_userid = ?
        ORDER BY post.post_at DESC
    `, userID)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostID, &post.Image, &post.Content, &post.PostAt, &post.Avatar, &post.FirstName, &post.LastName, &post.Username, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error in rows:", err)
		return nil, err
	}

	return posts, nil
}

func GetUserCommentedPosts(db *sql.DB, userid int, filter string) ([]Post, error) {
	if filter == "newest" {
		filter = "ASC"
	} else {
		filter = "DESC"
	}
	rows, err := db.Query(`
        SELECT DISTINCT post.postid, post.image, post.content, post.post_at
        FROM post
        JOIN comment ON post.postid = comment.post_postid
        JOIN user ON comment.user_userid = user.userid
        WHERE user.userid = ?
        ORDER BY post.post_at ?
    `, userid, filter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostID, &post.Image, &post.Content, &post.PostAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetAllPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query(`
        SELECT post.postid, post.image, post.content, post.post_at, post.user_userid, user.Username, user.F_name, user.L_name, user.Avatar,
               (SELECT COUNT(*) FROM likes WHERE likes.post_postid = post.postid) AS Likes,
               (SELECT COUNT(*) FROM dislikes WHERE dislikes.post_postid = post.postid) AS Dislikes,
               (SELECT COUNT(*) FROM comment WHERE comment.post_postid = post.postid) AS Comments
        FROM post
        JOIN user ON post.user_userid = user.userid
        ORDER BY post.post_at DESC
    `)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var postAt string
		if err := rows.Scan(&post.PostID, &post.Image, &post.Content, &postAt, &post.UserUserID, &post.Username, &post.FirstName, &post.LastName, &post.Avatar, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		post.PostAt, err = time.Parse(time.RFC3339, postAt)
		if err != nil {
			log.Println("Error parsing post_at:", err)
			return nil, err
		}

		// Fetch categories for the post
		categories, err := GetCategoriesForPost(db, post.PostID)
		if err != nil {
			log.Println("Error fetching categories for post:", err)
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error in rows:", err)
		return nil, err
	}

	return posts, nil
}

func GetCategoriesForPost(db *sql.DB, postID int) ([]Category, error) {
	rows, err := db.Query(`
        SELECT c.idcategories, c.name, c.description
        FROM categories c
        JOIN post_has_categories phc ON c.idcategories = phc.categories_idcategories
        WHERE phc.post_postid = ?
    `, postID)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
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

func GetFilteredPosts(db *sql.DB, filter string) ([]Post, error) {
	var rows *sql.Rows
	var err error

	switch filter {
	case "following":
		rows, err = db.Query(`
            SELECT post.postid, post.image, post.content, post.post_at, post.user_userid, user.Username, user.F_name, user.L_name, user.Avatar,
                   (SELECT COUNT(*) FROM likes WHERE likes.post_postid = post.postid) AS Likes,
                   (SELECT COUNT(*) FROM dislikes WHERE dislikes.post_postid = post.postid) AS Dislikes,
                   (SELECT COUNT(*) FROM comment WHERE comment.post_postid = post.postid) AS Comments
            FROM post
            JOIN user ON post.user_userid = user.userid
            ORDER BY post.post_at DESC
        `)
	case "friends":
		rows, err = db.Query(`
            SELECT post.postid, post.image, post.content, post.post_at, post.user_userid, user.Username, user.F_name, user.L_name, user.Avatar,
                   (SELECT COUNT(*) FROM likes WHERE likes.post_postid = post.postid) AS Likes,
                   (SELECT COUNT(*) FROM dislikes WHERE dislikes.post_postid = post.postid) AS Dislikes,
                   (SELECT COUNT(*) FROM comment WHERE comment.post_postid = post.postid) AS Comments
            FROM post
            JOIN user ON post.user_userid = user.userid
            ORDER BY post.post_at DESC
        `)
	case "top-rated":
		rows, err = db.Query(`
            SELECT post.postid, post.image, post.content, post.post_at, post.user_userid, user.Username, user.F_name, user.L_name, user.Avatar,
                   (SELECT COUNT(*) FROM likes WHERE likes.post_postid = post.postid) AS Likes,
                   (SELECT COUNT(*) FROM dislikes WHERE dislikes.post_postid = post.postid) AS Dislikes,
                   (SELECT COUNT(*) FROM comment WHERE comment.post_postid = post.postid) AS Comments
            FROM post
            JOIN user ON post.user_userid = user.userid
            ORDER BY Likes DESC, post.post_at DESC
        `)
	case "oldest":
		rows, err = db.Query(`
            SELECT post.postid, post.image, post.content, post.post_at, post.user_userid, user.Username, user.F_name, user.L_name, user.Avatar,
                   (SELECT COUNT(*) FROM likes WHERE likes.post_postid = post.postid) AS Likes,
                   (SELECT COUNT(*) FROM dislikes WHERE dislikes.post_postid = post.postid) AS Dislikes,
                   (SELECT COUNT(*) FROM comment WHERE comment.post_postid = post.postid) AS Comments
            FROM post
            JOIN user ON post.user_userid = user.userid
            ORDER BY post.post_at ASC
        `)
	default:
		rows, err = db.Query(`
            SELECT post.postid, post.image, post.content, post.post_at, post.user_userid, user.Username, user.F_name, user.L_name, user.Avatar,
                   (SELECT COUNT(*) FROM likes WHERE likes.post_postid = post.postid) AS Likes,
                   (SELECT COUNT(*) FROM dislikes WHERE dislikes.post_postid = post.postid) AS Dislikes,
                   (SELECT COUNT(*) FROM comment WHERE comment.post_postid = post.postid) AS Comments
            FROM post
            JOIN user ON post.user_userid = user.userid
			ORDER BY post.post_at DESC
        `, filter)
	}

	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var postAt string
		if err := rows.Scan(&post.PostID, &post.Image, &post.Content, &postAt, &post.UserUserID, &post.Username, &post.FirstName, &post.LastName, &post.Avatar, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		post.PostAt, err = time.Parse(time.RFC3339, postAt)
		if err != nil {
			log.Println("Error parsing post_at:", err)
			return nil, err
		}

		// Fetch categories for the post
		categories, err := GetCategoriesForPost(db, post.PostID)
		if err != nil {
			log.Println("Error fetching categories for post:", err)
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error in rows:", err)
		return nil, err
	}

	return posts, nil
}

func GetPostsByMultiCategory(db *sql.DB, categoryName string) ([]Post, error) {
	rows, err := db.Query(`
        SELECT post.postid, post.image, post.content, post.post_at, post.user_userid, user.Username, user.F_name, user.L_name, user.Avatar,
               (SELECT COUNT(*) FROM likes WHERE likes.post_postid = post.postid) AS Likes,
               (SELECT COUNT(*) FROM dislikes WHERE dislikes.post_postid = post.postid) AS Dislikes,
               (SELECT COUNT(*) FROM comment WHERE comment.post_postid = post.postid) AS Comments
        FROM post
        JOIN user ON post.user_userid = user.userid
        JOIN post_has_categories phc ON post.postid = phc.post_postid
        JOIN categories c ON phc.categories_idcategories = c.idcategories
   		WHERE c.name = ?
        ORDER BY post.post_at DESC
    `, categoryName)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var postAt string
		if err := rows.Scan(&post.PostID, &post.Image, &post.Content, &postAt, &post.UserUserID, &post.Username, &post.FirstName, &post.LastName, &post.Avatar, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		// Parse the postAt string into a time.Time object
		post.PostAt, err = time.Parse(time.RFC3339, postAt)
		if err != nil {
			log.Println("Error parsing post_at:", err)
			return nil, err
		}

		// Fetch categories for the post
		categories, err := GetCategoriesForPost(db, post.PostID)
		if err != nil {
			log.Println("Error fetching categories for post:", err)
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error in rows:", err)
		return nil, err
	}

	return posts, nil
}

func GetPostsByCategory(db *sql.DB, categoryName string) ([]Post, error) {
	rows, err := db.Query(`
        SELECT post.postid, post.image, post.content, post.post_at, post.user_userid, user.Username, user.F_name, user.L_name, user.Avatar,
               (SELECT COUNT(*) FROM likes WHERE likes.post_postid = post.postid) AS Likes,
               (SELECT COUNT(*) FROM dislikes WHERE dislikes.post_postid = post.postid) AS Dislikes,
               (SELECT COUNT(*) FROM comment WHERE comment.post_postid = post.postid) AS Comments
        FROM post
        JOIN user ON post.user_userid = user.userid
        JOIN post_has_categories phc ON post.postid = phc.post_postid
        JOIN categories c ON phc.categories_idcategories = c.idcategories
        WHERE c.name = ?
        ORDER BY post.post_at DESC
    `, categoryName)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var postAt string
		if err := rows.Scan(&post.PostID, &post.Image, &post.Content, &postAt, &post.UserUserID, &post.Username, &post.FirstName, &post.LastName, &post.Avatar, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		post.PostAt, err = time.Parse(time.RFC3339, postAt)
		if err != nil {
			log.Println("Error parsing post_at:", err)
			return nil, err
		}

		// Fetch categories for the post
		categories, err := GetCategoriesForPost(db, post.PostID)
		if err != nil {
			log.Println("Error fetching categories for post:", err)
			return nil, err
		}
		post.Categories = categories

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error in rows:", err)
		return nil, err
	}

	return posts, nil
}

func GetLastNotifications(db *sql.DB, userID int) ([]Notification, error) {
	rows, err := db.Query(`
        SELECT n.notificationid, n.user_userid, n.post_id, n.message, n.created_at, u.Avatar, u.Username
        FROM notifications n
        JOIN user u ON n.user_userid = u.userid
        WHERE n.user_userid = ?
        ORDER BY n.created_at DESC
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notification Notification
		var avatar sql.NullString

		err := rows.Scan(&notification.ID, &notification.UserID, &notification.PostID, &notification.Message, &notification.CreatedAt, &avatar, &notification.UserName)
		if err != nil {
			return nil, err
		}

		notification.UserImage = avatar.String
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

func InsertPost(db *sql.DB, content string, image sql.NullString, userID string) (int, error) {
	stmt, err := db.Prepare("INSERT INTO post (image, content, post_at, user_userid) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(image, content, time.Now(), userID)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastID), nil
}

func InsertPostCategory(db *sql.DB, postID int, categoryID int) error {
	stmt, err := db.Prepare("INSERT INTO post_has_categories (post_postid, categories_idcategories) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(postID, categoryID)
	return err
}

func GetUserPosts(db *sql.DB, userID int, filter string) ([]Post, error) {
	var x string
	if filter == "oldest" {
		x = "post.post_at ASC"
	} else {
		x = "post.post_at DESC"
	}

	rows, err := db.Query(`SELECT 
		post.postid, post.content, post.post_at, post.user_userid, 
		user.avatar, user.F_name, user.L_name, user.Username,
				 (SELECT COUNT(*) FROM likes WHERE likes.post_postid = post.postid) AS Likes,
               (SELECT COUNT(*) FROM dislikes WHERE dislikes.post_postid = post.postid) AS Dislikes,
               (SELECT COUNT(*) FROM comment WHERE comment.post_postid = post.postid) AS Comments 
	FROM post 
	JOIN user ON post.user_userid = user.userid 
	WHERE post.user_userid = ? ORDER BY `+x, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostID, &post.Content, &post.PostAt, &post.UserUserID, &post.Avatar, &post.FirstName, &post.LastName, &post.Username, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetFollowersCount(db *sql.DB, userID int) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM followers WHERE user_userid = ?", userID).Scan(&count)
	return count, err
}

func GetFollowingCount(db *sql.DB, userID int) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM following WHERE user_userid = ?", userID).Scan(&count)
	return count, err
}

func GetFriendsCount(db *sql.DB, userID int) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM friends WHERE user_userid = ?", userID).Scan(&count)
	return count, err
}

func IsFollowing(db *sql.DB, userID int, profileUserID int) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM followers WHERE user_userid = ? AND follower_userid = ?", profileUserID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetTotalUsersCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
	return count, err
}

func GetTotalPostsCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM post").Scan(&count)
	return count, err
}

func GetTotalCategoriesCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM categories").Scan(&count)
	return count, err
}

func GetAllReports(db *sql.DB) ([]Report, error) {
	rows, err := db.Query("SELECT id, post_id, reported_by, report_reason, created_at FROM reports ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []Report
	for rows.Next() {
		var report Report
		if err := rows.Scan(&report.ID, &report.PostID, &report.ReportedBy, &report.ReportReason, &report.CreatedAt); err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

func GetCommentsForPost(db *sql.DB, postID int) ([]Comment, error) {
	var comments []Comment

	query := `SELECT comment.commentid, comment.post_postid, comment.user_userid, user.F_name, user.L_name, user.Username, comment.content, comment.comment_at, user.Avatar,
	 		  (SELECT COUNT(*) FROM comment_dislikes WHERE comment_dislikes.commentid = comment.commentid) AS Dislikes,
			  (SELECT COUNT(*) FROM comment_likes WHERE comment_likes.commentid = comment.commentid) AS Likes
              FROM comment
              JOIN user ON comment.user_userid = user.userid
			  WHERE comment.post_postid = ?`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("GetCommentsForPost: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		var commentAt time.Time

		// Scan each row into the Comment struct
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.FirstName, &comment.LastName, &comment.Username, &comment.Content, &commentAt, &comment.Avatar, &comment.Dislikes, &comment.Likes); err != nil {
			return nil, fmt.Errorf("GetCommentsForPost: %v", err)
		}

		comment.CreatedAt = commentAt

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetCommentsForPost: %v", err)
	}

	return comments, nil
}

func ToggleLike(db *sql.DB, postID int, userID int) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM likes WHERE post_postid = ? AND user_userid = ?)", postID, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ToggleLike: %v", err)
	}

	if exists {
		_, err = db.Exec("DELETE FROM likes WHERE post_postid = ? AND user_userid = ?", postID, userID)
	} else {
		_, err = db.Exec("DELETE FROM dislikes WHERE post_postid = ? AND user_userid = ?", postID, userID)
		if err != nil {
			return fmt.Errorf("ToggleLike: %v", err)
		}
		_, err = db.Exec("INSERT INTO likes (post_postid, like_at, user_userid) VALUES (?, ?, ?)", postID, time.DateTime, userID)
	}
	return err
}

func ToggleDislike(db *sql.DB, postID int, userID int) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM dislikes WHERE post_postid = ? AND user_userid = ?)", postID, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ToggleDislike: %v", err)
	}

	if exists {
		_, err = db.Exec("DELETE FROM dislikes WHERE post_postid = ? AND user_userid = ?", postID, userID)
	} else {
		_, err = db.Exec("DELETE FROM likes WHERE post_postid = ? AND user_userid = ?", postID, userID)
		if err != nil {
			return fmt.Errorf("ToggleDislike: %v", err)
		}
		_, err = db.Exec("INSERT INTO dislikes (post_postid, dislike_at, user_userid) VALUES (?, ?, ?)", postID, time.DateTime, userID)
	}
	return err
}

func ToggleCommentLike(db *sql.DB, commentID int, userID int) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM comment_likes WHERE commentid = ? AND userid = ?)", commentID, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ToggleCommentLike: %v", err)
	}

	if exists {
		_, err = db.Exec("DELETE FROM comment_likes WHERE commentid = ? AND userid = ?", commentID, userID)
	} else {
		_, err = db.Exec("DELETE FROM comment_dislikes WHERE commentid = ? AND userid = ?", commentID, userID)
		if err != nil {
			return fmt.Errorf("ToggleCommentLike: %v", err)
		}
		_, err = db.Exec("INSERT INTO comment_likes (commentid, like_at, userid) VALUES (?, ?, ?)", commentID, time.DateTime, userID)
	}
	return err
}

func ToggleCommentDislike(db *sql.DB, commentID int, userID int) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM comment_dislikes WHERE commentid = ? AND userid = ?)", commentID, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ToggleCommentDislike: %v", err)
	}

	if exists {
		_, err = db.Exec("DELETE FROM comment_dislikes WHERE commentid = ? AND userid = ?", commentID, userID)
	} else {
		_, err = db.Exec("DELETE FROM comment_likes WHERE commentid = ? AND userid = ?", commentID, userID)
		if err != nil {
			return fmt.Errorf("ToggleCommentDislike: %v", err)
		}
		_, err = db.Exec("INSERT INTO comment_dislikes (commentid, dislike_at, userid) VALUES (?, ?, ?)", commentID, time.DateTime, userID)
	}
	return err
}

func GetUserLogs(db *sql.DB, userID int) ([]UserLog, error) {
	rows, err := db.Query("SELECT id, user_id, action, timestamp FROM user_logs WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []UserLog
	for rows.Next() {
		var log UserLog
		if err := rows.Scan(&log.ID, &log.UserID, &log.Action, &log.Timestamp); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func GetUserSessions(db *sql.DB, userID int) ([]UserSession, error) {
	rows, err := db.Query("SELECT sessionid, userid, start, end FROM sessions WHERE userid = ?", userID)
	if err != nil {
		log.Println("Failed to fetch user sessions:", err)
		return nil, err
	}
	defer rows.Close()

	var sessions []UserSession
	for rows.Next() {
		var session UserSession
		err := rows.Scan(&session.ID, &session.UserID, &session.Start, &session.End)
		if err != nil {
			log.Println("Failed to scan user session:", err)
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func GetFollowers(db *sql.DB, userID int) ([]User, error) {
	rows, err := db.Query(`
        SELECT user.userid, user.F_name, user.L_name, user.Username, user.Avatar
        FROM followers
        JOIN user ON followers.follower_userid = user.userid
        WHERE followers.user_userid = ?
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Avatar); err != nil {
			return nil, err
		}
		followers = append(followers, user)
	}
	return followers, nil
}

func GetFollowing(db *sql.DB, userID int) ([]User, error) {
	rows, err := db.Query(`
        SELECT user.userid, user.F_name, user.L_name, user.Username, user.Avatar
        FROM following
        JOIN user ON following.following_userid = user.userid
        WHERE following.user_userid = ?
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Avatar); err != nil {
			return nil, err
		}
		following = append(following, user)
	}
	return following, nil
}

func GetFriends(db *sql.DB, userID int) ([]User, error) {
	rows, err := db.Query(`
        SELECT user.userid, user.F_name, user.L_name, user.Username, user.Avatar
        FROM friends
        JOIN user ON friends.friend_userid = user.userid
        WHERE friends.user_userid = ?
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Avatar); err != nil {
			return nil, err
		}
		friends = append(friends, user)
	}
	return friends, nil
}

func GetTotalLikes(db *sql.DB, userID int) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM likes WHERE user_userid = ?", userID).Scan(&count)
	return count, err
}

func GetTotalPosts(db *sql.DB, userID int) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM post WHERE user_userid = ?", userID).Scan(&count)
	return count, err
}

func GetUserByID(db *sql.DB, userID int) (User, error) {
	var user User
	err := db.QueryRow("SELECT userid, F_name, L_name, Username, Email, Avatar, role_id FROM user WHERE userid = ?", userID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Avatar, &user.RoleID)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetRoleNameByID(db *sql.DB, roleID int) (string, error) {
	var roleName string
	err := db.QueryRow("SELECT role_name FROM user_roles WHERE roleid = ?", roleID).Scan(&roleName)
	if err != nil {
		log.Printf("Error fetching role name for roleID %d: %v\n", roleID, err) // Add this line for debugging
		return "", err
	}
	return roleName, nil
}
