package handlers

import (
	"database/sql"
	"forum/database"
	"log"
	"net/http"
	"strconv"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		userData := database.GetAuthenticatedUserData(db, r)
		if !userData.IsAuthenticated {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		content := r.FormValue("content")
		userID, ok := GetAuthenticatedUserID(r)
		if !ok {
			http.Error(w, "Unauthorized to use", http.StatusUnauthorized)
			return
		}
		err = database.InsertComment(db, postID, userID, content)
		if err != nil {
			http.Error(w, "Internal Server 500 Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
