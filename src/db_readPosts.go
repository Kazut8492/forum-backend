package src

import (
	"database/sql"
	"log"
	"strings"
)

func ReadPosts(db *sql.DB) []Post {

	postRows, err := db.Query(`
		SELECT * FROM post
		ORDER BY post_id
	`)
	if err != nil {
		panic(err.Error())
	}
	defer postRows.Close()

	var result []Post
	for postRows.Next() {
		post := Post{}
		err = postRows.Scan(&post.ID, &post.Title, &post.Content, &post.CategoryStr, &post.CreatorUsrName)
		if err != nil {
			panic(err.Error())
		}

		// Read likes
		likeRows, err := db.Query(`
			SELECT * FROM like WHERE post_id = ?
		`, post.ID)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer likeRows.Close()
		var likes []Like
		for likeRows.Next() {
			var like Like
			err = likeRows.Scan(&like.ID, &like.PostId, &like.CommentId, &like.CreatorUsrName)
			if err != nil {
				panic(err.Error())
			}
			likes = append(likes, like)
		}

		// Read dislikes
		dislikeRows, err := db.Query(`
			SELECT * FROM dislike WHERE post_id = ?
		`, post.ID)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer dislikeRows.Close()
		var dislikes []Dislike
		for dislikeRows.Next() {
			var dislike Dislike
			err = dislikeRows.Scan(&dislike.ID, &dislike.PostId, &dislike.CommentId, &dislike.CreatorUsrName)
			if err != nil {
				panic(err.Error())
			}
			dislikes = append(dislikes, dislike)
		}

		// Read comments
		comments := ReadComments(db, post.ID)

		post.Likes = likes
		post.Dislikes = dislikes
		post.Comments = comments

		result = append(result, post)
	}
	for index, post := range result {
		result[index].CategoryArr = strings.Split(post.CategoryStr, ",")
	}

	return result
}
