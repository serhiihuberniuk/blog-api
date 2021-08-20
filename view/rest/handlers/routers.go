package handlers

import "github.com/gorilla/mux"

func (h *Handlers) ApiRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/auth", h.Login).Methods("POST")
	router.HandleFunc("/users", h.CreateUser).Methods("POST")

	router.HandleFunc("/users/{id}", h.authMiddleware.Auth(h.GetUser)).Methods("GET")
	router.HandleFunc("/users", h.authMiddleware.Auth(h.UpdateUser)).Methods("PUT")
	router.HandleFunc("/users", h.authMiddleware.Auth(h.DeleteUser)).Methods("DELETE")

	router.HandleFunc("/posts", h.authMiddleware.Auth(h.CreatePost)).Methods("POST")
	router.HandleFunc("/posts/{id}", h.authMiddleware.Auth(h.GetPost)).Methods("GET")
	router.HandleFunc("/posts/{id}", h.authMiddleware.Auth(h.UpdatePost)).Methods("PUT")
	router.HandleFunc("/posts/{id}", h.authMiddleware.Auth(h.DeletePost)).Methods("DELETE")
	router.HandleFunc("/posts", h.authMiddleware.Auth(h.GetListOfPosts)).Methods("GET")

	router.HandleFunc("/comments", h.authMiddleware.Auth(h.CreateComment)).Methods("POST")
	router.HandleFunc("/comments/{id}", h.authMiddleware.Auth(h.GetComment)).Methods("GET")
	router.HandleFunc("/comments/{id}", h.authMiddleware.Auth(h.UpdateComment)).Methods("PUT")
	router.HandleFunc("/comments/{id}", h.authMiddleware.Auth(h.DeleteComment)).Methods("DELETE")
	router.HandleFunc("/comments", h.authMiddleware.Auth(h.GetListOfComments)).Methods("GET")

	return router
}
