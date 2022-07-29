package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
	"shorterer-link/api/models"
)

func TestCreateUrl(t *testing.T) {

	err := refreshUserAndUrlTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	token, err := server.SignIn(user.Email, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)
	samples := []struct {
		inputJSON    string
		statusCode   int
		originalUrl  string
		customUrl    string
		user_id      uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"originalUrl":"https://www.linkedin.com/in/ahsanulwalad", "customUrl": "the content"}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			originalUrl:  "https://www.linkedin.com/in/ahsanulwalad",
			customUrl:    "the content",
			user_id:      user.ID,
			errorMessage: "",
		},
		{
			inputJSON:    `{"originalUrl":"The title", "customUrl": "the content"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Url Already Taken",
		},
		{
			// When no token is passed
			inputJSON:    `{"originalUrl":"When no token is passed", "customUrl": "the content", "user_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"originalUrl":"When incorrect token is passed", "customUrl": "the content", "user_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"originalUrl": "", "customUrl": "The content", "user_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			inputJSON:    `{"originalUrl": "This is a title", "customUrl": "", "user_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Content",
		},
		{
			inputJSON:    `{"originalUrl": "This is an awesome title", "customUrl": "the content"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required User",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"originalUrl": "This is an awesome title", "customUrl": "the content", "user_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/url", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateUrl)
		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["originalUrl"], v.originalUrl)
			assert.Equal(t, responseMap["customUrl"], v.customUrl)
			assert.Equal(t, responseMap["user_id"], float64(v.user_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetUrl(t *testing.T) {

	err := refreshUserAndUrlTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUsersAndUrl()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/url-list", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	log.Println(req)

	rr := httptest.NewRecorder()
	log.Println(rr)
	handler := http.HandlerFunc(server.GetAllUrl)
	handler.ServeHTTP(rr, req)

	var url []models.Url
	err = json.Unmarshal([]byte(rr.Body.String()), &url)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(url), 2)
}
func TestGetUrlByID(t *testing.T) {

	err := refreshUserAndUrlTable()
	if err != nil {
		log.Fatal(err)
	}
	url, err := seedOneUserAndUrl()
	if err != nil {
		log.Fatal(err)
	}
	urlSample := []struct {
		id           string
		statusCode   int
		originalUrl  string
		customUrl    string
		user_id      uint32
		errorMessage string
	}{
		{
			id:          strconv.Itoa(int(url.ID)),
			statusCode:  200,
			originalUrl: url.OriginalUrl,
			customUrl:   url.CustomUrl,
			user_id:     url.UserID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range urlSample {

		req, err := http.NewRequest("GET", "/url", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetUrl)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, url.OriginalUrl, responseMap["originaUrl"])
			assert.Equal(t, url.CustomUrl, responseMap["customUrl"])
			assert.Equal(t, float64(url.UserID), responseMap["user_id"]) //the response author id is float64
		}
	}
}

func TestDeleteUrl(t *testing.T) {

	var UrlUserEmail, UrlUserPassword string
	var UrlUserID uint32
	var AuthUrlID uint64

	err := refreshUserAndUrlTable()
	if err != nil {
		log.Fatal(err)
	}
	users, url, err := seedUsersAndUrl()
	if err != nil {
		log.Fatal(err)
	}
	//Let's get only the Second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		UrlUserEmail = user.Email
		UrlUserPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}
	//Login the user and get the authentication token
	token, err := server.SignIn(UrlUserEmail, UrlUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the second post
	for _, url := range url {
		if url.ID == 1 {
			continue
		}
		AuthUrlID = url.ID
		UrlUserID = url.UserID
	}
	urlSample := []struct {
		id           string
		user_id      uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(int(AuthUrlID)),
			user_id:      UrlUserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(int(AuthUrlID)),
			user_id:      UrlUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(int(AuthUrlID)),
			user_id:      UrlUserID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(1)),
			user_id:      1,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range urlSample {

		req, _ := http.NewRequest("GET", "/url", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteUrl)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
