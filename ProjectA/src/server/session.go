package server

import "github.com/gorilla/sessions"

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func init() {
store.Options = &sessions.Options{
    Path:     "/",
    MaxAge:   60 * 60, // 1 hour in seconds
	HttpOnly: true,

}}