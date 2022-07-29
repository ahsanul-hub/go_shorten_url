package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"shorterer-link/api/auth"
	"shorterer-link/api/models"
	"shorterer-link/api/responses"
	"shorterer-link/api/utils"
)

func newRedisClient(host string, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})
	return client
}

var rdb = newRedisClient("localhost:6379", "")
var ctx = context.Background()

func (server *Server) CreateUrl(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	url := models.Url{}
	err = json.Unmarshal(body, &url)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//url.Prepare()
	err = url.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	url.UserID = uid
	key := url.CustomUrl
	data, _ := json.Marshal(url)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	add := rdb.HSet(ctx, "url", key, data)
	if err := add.Err(); err != nil {
		//fmt.Printf("unable to set data. error", err)
		return
	}
	if uid != url.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	urlCreated, err := url.SaveUrl(server.DB)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, urlCreated.ID))
	responses.JSON(w, http.StatusCreated, urlCreated)
}

func (server *Server) GetAllUrl(w http.ResponseWriter, r *http.Request) {

	url := models.Url{}
	//var a interface{}
	res := rdb.HGetAll(ctx, "url").Scan(url)
	if res == nil {
		//fmt.Printf("unable to set data. error", err)
		json.NewEncoder(w).Encode(res)
		return
	}

	urlList, err := url.FindAllUrl(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, urlList)
}

func (server *Server) GetUrl(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["customUrl"]

	url := models.Url{}

	res, err := rdb.HGet(ctx, "url", key).Result()
	if err != nil {
		log.Println("unable to Get data. error", err)
	}

	err = json.Unmarshal([]byte(res), &url)
	if err == nil {
		json.NewEncoder(w).Encode(url)
		return
	}
	urlReceived, err := url.FindUrlByLink(server.DB, key)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	data, _ := json.Marshal(urlReceived)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	add := rdb.HSet(ctx, "url", key, data)
	if err := add.Err(); err != nil {
		//fmt.Printf("unable to set data. error", err)
		return
	}
	responses.JSON(w, http.StatusOK, urlReceived)
}

func (server *Server) RedirectUrl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["customUrl"]
	url := models.Url{}

	res, err := rdb.HGet(ctx, "url", key).Result()
	if err != nil {
		log.Println("unable to Get data. error", err)
	}

	err = json.Unmarshal([]byte(res), &url)
	if err == nil {
		http.Redirect(w, r, url.OriginalUrl, 301)
		return
	}

	urlReceived, err := url.FindUrlByLink(server.DB, key)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, urlReceived.OriginalUrl, 301)

}

func (server *Server) UpdateUrl(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["customUrl"]

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	urlUpdate := models.Url{}
	// Check if the url exist
	url := models.Url{}

	err = server.DB.Debug().Model(models.Url{}).Where("custom_url = ?", key).Take(&url).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Url not found"))
		return
	}

	if uid != url.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	err = json.Unmarshal(body, &urlUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	urlUpdate.UserID = uid
	if uid != urlUpdate.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	//postUpdate.Prepare()
	err = urlUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	urlUpdate.ID = url.ID //this is important to tell the model the post id to update, the other update field are set above

	urlUpdated, err := urlUpdate.UpdateAUrl(server.DB, urlUpdate.ID)

	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	data, _ := json.Marshal(urlUpdated)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	del := rdb.HDel(ctx, "url", key)
	if err := del.Err(); err != nil {
		//fmt.Printf("unable to delete previous data from redis. error", err)
		return
	}

	add := rdb.HSet(ctx, "url", urlUpdated.CustomUrl, data)
	if err := add.Err(); err != nil {
		//fmt.Printf("unable to set data to redis. error", err)
		return
	}

	responses.JSON(w, http.StatusOK, urlUpdated)
}

func (server *Server) DeleteUrl(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["customUrl"]

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	url := models.Url{}
	err = server.DB.Debug().Model(models.Url{}).Where("custom_url = ?", key).Take(&url).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this post?
	if uid != url.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = url.DeleteAUrl(server.DB, key, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	_, err = rdb.HDel(ctx, "url", key).Result()
	if err != nil {
		log.Println("unable to delete data on cache")
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusOK, "success delete url")
}
