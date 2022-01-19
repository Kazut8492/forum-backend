package src

import (
	"database/sql"
	"log"
)

func ReadComments(db *sql.DB, postId int) []Comment {
	db_readComments := `
		SELECT * FROM comment WHERE post_id = ?
		ORDER BY comment_id
	`

	rows, err := db.Query(db_readComments, postId)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var result []Comment
	for rows.Next() {
		comment := Comment{}
		err = rows.Scan(&comment.ID, &comment.PostId, &comment.Title, &comment.Content, &comment.CreatorUsrName)
		if err != nil {
			panic(err)
		}

		likeRows, err := db.Query(`
			SELECT * FROM like WHERE post_id = ? AND comment_id = ?
		`, postId, comment.ID)
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

		dislikeRows, err := db.Query(`
			SELECT * FROM dislike WHERE post_id = ? AND comment_id = ?
		`, postId, comment.ID)
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

		comment.Likes = likes
		comment.Dislikes = dislikes

		result = append(result, comment)
	}
	return result
}
