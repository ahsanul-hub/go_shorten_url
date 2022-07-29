package controllers

import "shorterer-link/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Wellcome)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/register", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Url routes
	s.Router.HandleFunc("/url", middlewares.SetMiddlewareJSON(s.CreateUrl)).Methods("POST")
	s.Router.HandleFunc("/url-list", middlewares.SetMiddlewareJSON(s.GetAllUrl)).Methods("GET")
	s.Router.HandleFunc("/url/{customUrl}", middlewares.SetMiddlewareJSON(s.GetUrl)).Methods("GET")
	s.Router.HandleFunc("/url/{customUrl}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUrl))).Methods("PUT")
	s.Router.HandleFunc("/url/{customUrl}", middlewares.SetMiddlewareAuthentication(s.DeleteUrl)).Methods("DELETE")
	s.Router.HandleFunc("/{customUrl}", middlewares.SetMiddlewareJSON(s.RedirectUrl)).Methods("GET")
}
