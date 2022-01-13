package src

// Commented out becasue there is one-liner way to do it. Pls refer to Login/Signup Submit Handlers
// func ReadUser(w http.ResponseWriter, db *sql.DB, loginUser User) User {
// 	statement, _ := db.Query("SELECT * FROM user WHERE username = ?", loginUser.Username)
// 	// if err != nil {
// 	// 	w.WriteHeader(500)
// 	// 	fmt.Println(err.Error())
// 	// 	log.Fatal(1)
// 	// }
// 	defer statement.Close()

// 	user := User{}
// 	for statement.Next() {
// 		err := statement.Scan(&user.ID, &user.Username, &user.Email, &user.Pass)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			return User{}
// 		}
// 	}
// 	return user
// }
