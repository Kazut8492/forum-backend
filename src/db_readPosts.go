package src

import (
	"database/sql"
	"strings"
)

func ReadPosts(db *sql.DB) []Post {
	db_readPosts := `
		SELECT * FROM post
		ORDER BY post_id
	`

	rows, err := db.Query(db_readPosts)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var result []Post
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.CategoryStr, &post.CreatorUsrName, &post.Like, &post.DisLike)
		if err != nil {
			panic(err.Error())
		}
		result = append(result, post)
	}
	for index, post := range result {
		result[index].CategoryArr = strings.Split(post.CategoryStr, ",")
	}

	return result
}
