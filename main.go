package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

type PostType []struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func fetchPosts() PostType {
	postSampleAPI := "https://jsonplaceholder.typicode.com/posts?_limit=3"
	// Perform the GET request
	resp, err := http.Get(postSampleAPI)
	if err != nil {
		log.Fatalf("Error occurred making a request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Declare a map to store the parsed JSON
	result := make(PostType, 100)

	// Unmarshal the JSON into the map
	abcError := json.Unmarshal(body, &result)
	if abcError != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}
	return result
}

func main() {
	component := 
	dashboardHandler := func(w http.ResponseWriter, r *http.Request) {
		dashTemplate := template.Must(template.ParseFiles("dashboard.html"))
		postsFetched := fetchPosts()
		dashTemplate.Execute(w, map[string]PostType{"postsFetched": postsFetched})
	}
	receiveNewPostHandler := func(w http.ResponseWriter, r *http.Request) {
		postTitle := r.PostFormValue("post-title")
		postBody := r.PostFormValue("post-body")
		newPostCard := `
		<h2>%s</h2>
		<p>%s</p>
		`
		finalNewPostRender := fmt.Sprintf(newPostCard, postTitle, postBody)
		w.Write([]byte(finalNewPostRender))

	}
	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/add-new-post", receiveNewPostHandler)
	fmt.Println("Server is listening to port 3003 ...")
	log.Fatal(http.ListenAndServe(":3003", nil))
}
