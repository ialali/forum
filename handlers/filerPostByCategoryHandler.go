package handlers

import (
	"database/sql"
	"forum/database"
	"log"
	"net/http"
	"text/template"
)

func FilterPosts(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID, ok := GetAuthenticatedUserID(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 1. Parse the filtering criteria from the request.
	category := r.FormValue("category")
	created := r.FormValue("created")
	liked := r.FormValue("liked")

	// Initialize an empty slice to hold filtered posts.
	var filteredPosts []database.Post

	// 2. Based on the criteria, call the corresponding functions to retrieve the filtered posts.
	switch {
	case category != "":
		// Filter by category
		filteredPosts, _ = database.GetPostsByCategory(db, category)
	case created == "true":
		// Filter created posts by the authenticated user
		filteredPosts, _ = database.GetOwnedPosts(db, userID)
	case liked == "true":
		// Filter liked posts by the authenticated user
		filteredPosts, _ = database.GetLikedPosts(db, userID, true)
	default:
		// No criteria selected, show all posts
		filteredPosts, _ = database.GetPosts(db)
	}

	// 3. For each post, retrieve like/dislike counts and usernames for comments.
	for i, post := range filteredPosts {
		comments, err := database.GetCommentsForPost(db, post.ID)
		if err != nil {
			http.Error(w, "Error fetching comments", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		// Fetch like/dislike counts and usernames for comments
		for j, comment := range comments {
			likeCount, dislikeCount, err := database.GetCommentLikesCount(db, comment.ID)
			if err != nil {
				http.Error(w, "Error fetching comment likes/dislikes", http.StatusInternalServerError)
				log.Println(err)
				return
			}

			// username, err := database.GetUserByID(db, comment.UserID)
			// if err != nil {
			// 	http.Error(w, "Error fetching username", http.StatusInternalServerError)
			// 	log.Println(err)
			// 	return
			// }

			comments[j].LikeCount = likeCount
			comments[j].DislikeCount = dislikeCount
			// comments[j].Username = username
		}

		// Assign comments to the post
		filteredPosts[i].Comments = comments
	}

	// 4. Render the filtered posts to the page.
	userData := GetAuthenticatedUserData(db, r)

	data := struct {
		IsAuthenticated bool
		Username        string
		Posts           []database.Post
	}{
		IsAuthenticated: userData.IsAuthenticated,
		Username:        userData.Username,
		Posts:           filteredPosts,
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error Parsing index.html", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
