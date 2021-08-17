package handlers

import "github.com/gorilla/mux"

func (h *Handlers) ApiRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/auth", h.Login).Methods("POST")

	router.HandleFunc("/users", h.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")

	router.HandleFunc("/posts", h.CreatePost).Methods("POST")
	router.HandleFunc("/posts/{id}", h.GetPost).Methods("GET")
	router.HandleFunc("/posts/{id}", h.UpdatePost).Methods("PUT")
	router.HandleFunc("/posts/{id}", h.DeletePost).Methods("DELETE")
	router.HandleFunc("/posts", h.GetListOfPosts).Methods("GET")

	router.HandleFunc("/comments", h.CreateComment).Methods("POST")
	router.HandleFunc("/comments/{id}", h.GetComment).Methods("GET")
	router.HandleFunc("/comments/{id}", h.UpdateComment).Methods("PUT")
	router.HandleFunc("/comments/{id}", h.DeleteComment).Methods("DELETE")
	router.HandleFunc("/comments", h.GetListOfComments).Methods("GET")

	return router
}
