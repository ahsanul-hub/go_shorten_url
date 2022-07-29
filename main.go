package main

import "shorterer-link/api"

//type Link struct {
//	Id           int    `json:"id"`
//	OriginalLink string `json:"originalLink"`
//	CustomLink   string `json:"customLink"`
//}
//
//var Data []Link

//var IndexData = map[string]string{}
//
//func newRedisClient(host string, password string) *redis.Client {
//	client := redis.NewClient(&redis.Options{
//		Addr:     host,
//		Password: password,
//		DB:       0,
//	})
//	return client
//}
//
//var rdb = newRedisClient("localhost:6379", "")
//var ctx = context.Background()
//
//func showLinks(w http.ResponseWriter, r *http.Request) {
//	getLinks := rdb.HGetAll(ctx, "links")
//	if err := getLinks.Err(); err != nil {
//		fmt.Printf("unable to Get data. error", err)
//		return
//	}
//	res, err := getLinks.Result()
//	if err != nil {
//		fmt.Printf("unable to Get data", err)
//		return
//	}
//	json.NewEncoder(w).Encode(res)
//
//}
//
//func showLink(w http.ResponseWriter, r *http.Request) {
//	//vars := mux.Vars(r)
//	//key := vars["CustomLink"]
//
//	//log.Print("Key: " + key)
//	//for _, link := range Data {
//	//	if link.OriginalLink == key {
//	//		json.NewEncoder(w).Encode(link)
//	//		log.Print(link)
//	//	}
//	//}
//	json.NewEncoder(w).Encode(IndexData)
//}
//
//func handleAdd(w http.ResponseWriter, r *http.Request) {
//
//	reqBody, _ := ioutil.ReadAll(r.Body)
//	//log.Println(reqBody)
//	var link map[string]string
//	json.Unmarshal(reqBody, &link)
//
//	key := link["customLink"]
//	value := link["originalLink"]
//	//data, _ := json.Marshal(link)
//	add := rdb.HSet(ctx, "links", key, value)
//	if err := add.Err(); err != nil {
//		fmt.Printf("unable to set data. error", err)
//		return
//	}
//
//	//log.Println(IndexData)
//	//link["tesguys"] = "foo"
//	//log.Print(link)
//	//
//	//link.Id = Data[len(Data)-1].Id + 1
//	//Data = append(Data, link)
//	//
//	json.NewEncoder(w).Encode(link)
//
//}
//func redirect(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	key := vars["CustomLink"]
//
//	var value = IndexData[key]
//	http.Redirect(w, r, value, 301)
//	log.Print(value)
//
//}
//
//func deleteLink(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	key := vars["CustomLink"]
//	_, err := rdb.HDel(ctx, "links", key).Result()
//	if err == redis.Nil {
//		json.NewEncoder(w).Encode("Link doesnt exist")
//	} else {
//		json.NewEncoder(w).Encode(` Success delete link `)
//	}
//
//}
//
//func updateLink(w http.ResponseWriter, r *http.Request) {
//	reqBody, _ := ioutil.ReadAll(r.Body)
//	//var newlink Link
//	//json.Unmarshal(reqBody, &newlink)
//	//log.Println(newlink)
//	//str := vars["id"]
//	//var id, _ = strconv.Atoi(str)
//	//for index, link := range Data {
//	//	if link.Id == id {
//	//		Data[index].CustomLink = newlink.CustomLink
//	//		Data[index].OriginalLink = newlink.OriginalLink
//	//		log.Println(Data)
//	//	}
//	//}
//	vars := mux.Vars(r)
//
//	var link map[string]string
//	json.Unmarshal(reqBody, &link)
//
//	key := vars["CustomLink"]
//	value := link["originalLink"]
//
//	_, err := rdb.HGet(ctx, "links", key).Result()
//	if err == redis.Nil {
//		json.NewEncoder(w).Encode("Link doesnt exist")
//	} else {
//		add := rdb.HSet(ctx, "links", key, value)
//		if err := add.Err(); err != nil {
//			fmt.Printf("unable to set data. error", err)
//			return
//		}
//		json.NewEncoder(w).Encode(link)
//	}
//
//}
//func main() {
//	//Data = []Link{
//	//	{Id: 1, OriginalLink: "www.hello.com", CustomLink: "www.custom.com"},
//	//	{Id: 2, OriginalLink: "https://github.com/ahsanul-hub", CustomLink: "github-aw"},
//	//}
//	//IndexData = make(map[string]Link)
//	//IndexData["1"] = Link{Id: 1, OriginalLink: "www.hello.com", CustomLink: "www.custom.com"}
//	//{
//	//	{"Id": "1", "OriginalLink": "www.hello.com", "CustomLink": "www.custom.com"},
//	//	{"Id": "2", "OriginalLink": "https://github.com/ahsanul-hub", "CustomLink": "github-aw"},
//	//}
//
//	myRouter := mux.NewRouter().StrictSlash(true)
//	myRouter.HandleFunc("/links", showLinks)
//	myRouter.HandleFunc("/add-link", handleAdd).Methods("POST")
//	myRouter.HandleFunc("/link/{CustomLink}", deleteLink).Methods("DELETE")
//	myRouter.HandleFunc("/link/{CustomLink}", updateLink).Methods("PATCH")
//	myRouter.HandleFunc("/link/{OriginalLink}", showLink)
//	myRouter.HandleFunc("/{CustomLink}", redirect)
//
//	fmt.Println("server started at localhost:9000")
//	log.Fatal(http.ListenAndServe(":9000", myRouter))
//}
//func tes() {
//	var data = map[string]map[string]string{}
//	data["Date_2"] = make(map[string]string)
//
//	data["Date_2"]["Sistem_A"] = "violet"
//
//	fmt.Println(data)
//}

func main() {
	api.Run()
}
