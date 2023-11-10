package handlers

import (
	"database/sql"
	"forum/database"
	"log"
	"net/http"
	"text/template"
)

func FilterPosts(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Check if the user is authenticated
	userID, isAuthenticated := GetAuthenticatedUserID(r)
	if !isAuthenticated {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 1. Parse the filtering criteria from the request.
	category := r.FormValue("category")
	created := r.FormValue("created")
	liked := r.FormValue("liked")

	// Initialize an empty slice to hold filtered posts.
	var posts []database.Post

	// 2. Based on the criteria, call the corresponding functions to retrieve the filtered posts.
	switch {
	case category != "":
		// Filter by category
		posts, _ = database.GetPostsByCategory(db, category)
	case created == "true":
		// Filter created posts by the authenticated user
		posts, _ = database.GetOwnedPosts(db, userID)
	case liked == "true":
		// Filter liked posts by the authenticated user
		posts, _ = database.GetLikedPosts(db, userID, true)
	default:
		// No criteria selected, show all posts
		posts, _ = database.GetPosts(db)
	}

	// 3. For each post, retrieve like/dislike counts and usernames for comments.
	for i, post := range posts {
		comments, err := database.GetCommentsForPost(db, posts[i].ID)
		if err != nil {
			http.Error(w, "Error fetching comments", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		likeCount, dislikeCount, err := database.GetPostLikesCount(db, post.ID)
		if err != nil {
			http.Error(w, "Error fetching post likes/dislikes", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		// Fetch like/dislike counts and usernames for comments
		for j := range comments {
			likeCount, dislikeCount, err := database.GetCommentLikesCount(db, comments[j].ID)
			if err != nil {
				http.Error(w, "Error fetching comment likes/dislikes", http.StatusInternalServerError)
				log.Println(err)
				return
			}

			comments[j].LikeCount = likeCount
			comments[j].DislikeCount = dislikeCount

			// Assuming you have a function to get the username by user ID
			// username, err := database.GetUserByID(db, comments[j].UserID)
			// if err != nil {
			// 	http.Error(w, "Error fetching username", http.StatusInternalServerError)
			// 	log.Println(err)
			// 	return
			// }

			// Assign the username to the comment
			// comments[j].Username = username
		}

		// Assign comments to the post
		post.Comments = comments
		post.LikeCount = likeCount
		post.DislikeCount = dislikeCount
		posts[i] = post
	}

	// 4. Render the filtered posts to the page.
	userData := GetAuthenticatedUserData(db, r)

	data := database.PageData{
		IsAuthenticated: userData.IsAuthenticated,
		Username:        userData.Username,
		Posts:           posts,
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error Parsing index.html", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering the template", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
