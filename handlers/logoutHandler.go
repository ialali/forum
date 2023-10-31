package handlers

import "net/http"

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the session token by setting the cookie to expire.
	http.SetCookie(w, &http.Cookie{
		Name:   "sessionToken",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Expiration time in seconds (0 or negative value)
	})

	// Redirect the user to the login page or any other desired page.
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
