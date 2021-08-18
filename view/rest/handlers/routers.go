package handlers

import "github.com/gorilla/mux"

func (h *Handlers) ApiRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/auth", h.Login).Methods("POST")
	router.HandleFunc("/users", h.CreateUser).Methods("POST")

	router.HandleFunc("/users/{id}", h.Auth(h.GetUser)).Methods("GET")
	router.HandleFunc("/users/", h.Auth(h.UpdateUser)).Methods("PUT")
	router.HandleFunc("/users/", h.Auth(h.DeleteUser)).Methods("DELETE")

	router.HandleFunc("/posts", h.Auth(h.CreatePost)).Methods("POST")
	router.HandleFunc("/posts/{id}", h.Auth(h.GetPost)).Methods("GET")
	router.HandleFunc("/posts/{id}", h.Auth(h.UpdatePost)).Methods("PUT")
	router.HandleFunc("/posts/{id}", h.Auth(h.DeletePost)).Methods("DELETE")
	router.HandleFunc("/posts", h.Auth(h.GetListOfPosts)).Methods("GET")

	router.HandleFunc("/comments", h.Auth(h.CreateComment)).Methods("POST")
	router.HandleFunc("/comments/{id}", h.Auth(h.GetComment)).Methods("GET")
	router.HandleFunc("/comments/{id}", h.Auth(h.UpdateComment)).Methods("PUT")
	router.HandleFunc("/comments/{id}", h.Auth(h.DeleteComment)).Methods("DELETE")
	router.HandleFunc("/comments", h.Auth(h.GetListOfComments)).Methods("GET")

	return router
}
