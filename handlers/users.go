package handlers

import "github.com/hiosi123/udongmat/storage"

type UserHandler struct {
	Storage *storage.UserStorage
	// SessionManager *auth.SessionManager
}
