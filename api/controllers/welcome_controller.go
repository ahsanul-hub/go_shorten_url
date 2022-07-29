package controllers

import (
	"net/http"

	"shorterer-link/api/responses"
)

func (server *Server) Wellcome(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To my API")
}
