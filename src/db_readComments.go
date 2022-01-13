package src

import "database/sql"

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
		err = rows.Scan(&comment.ID, &comment.PostId, &comment.Title, &comment.Content, &comment.CreatorUsrName, &comment.Like, &comment.DisLike)
		if err != nil {
			panic(err)
		}
		result = append(result, comment)
	}
	return result
}
