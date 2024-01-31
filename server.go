package main

import (
	"encoding/json"
	// "fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
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
	unmarshalizingError := json.Unmarshal(body, &result)
	if unmarshalizingError != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}
	return result
}

func generateSlug(s string) string {
    // Make the string URL-friendly
    slug := strings.ToLower(s)
    slug = strings.ReplaceAll(slug, " ", "-")
    // Remove special characters
    reg, _ := regexp.Compile("[^a-zA-Z0-9-]+")
    slug = reg.ReplaceAllString(slug, "")
    return slug
}

func main() {
    // Initialize standard Go html template engine
    engine := html.New("./views", ".html")
    // If you want other engine, just replace with following
    // Create a new engine with django
    // engine := django.New("./views", ".django")

    app := fiber.New(fiber.Config{
        Views: engine,
    })
    app.Get("/", func(c *fiber.Ctx) error {
		postsFetched := fetchPosts()
        // Render index template
        return c.Render("all-posts", fiber.Map{
            "postsFetched": postsFetched,
        })
    })
	app.Get("/new-post", func(c *fiber.Ctx) error {
        // Render index template
        return c.Render("new-post", fiber.Map{
            "Username": "Guest User",
        })
    })
	app.Post("/create-new-post", func(c *fiber.Ctx) error {
		postTitle := c.FormValue("post-title","no title entered!")
		// postBody := c.FormValue("post-body","no body entered!")
		// Generate slug from postTitle
        slug := generateSlug(postTitle)
		return c.Redirect("/post/" + slug)
		// TODO: save the new post in DB
        // Render index template
        // return c.Render("post", fiber.Map{
        //     "Username": "Guest User",
		// 	"Title": postTitle,
		// 	"Body": postBody,
        // })
    })
	app.Get("/post/:slug", func(c *fiber.Ctx) error {
        slug := c.Params("slug")
        // TODO: Retrieve the post data based on slug from the DB
        return c.Render("post", fiber.Map{
            "Username": "Guest User",
			"Title": slug,
            // Add retrieved post data here
        })
    })

    log.Fatal(app.Listen(":3003"))
}
// func main() {
// 	dashboardHandler := func(w http.ResponseWriter, r *http.Request) {
// 		dashTemplate := template.Must(template.ParseFiles("dashboard.html"))
// 		postsFetched := fetchPosts()
// 		dashTemplate.Execute(w, map[string]PostType{"postsFetched": postsFetched})
// 	}
// 	receiveNewPostHandler := func(w http.ResponseWriter, r *http.Request) {
// 		postTitle := r.PostFormValue("post-title")
// 		postBody := r.PostFormValue("post-body")
// 		newPostSection := `
// 		<h2>%s</h2>
// 		<p>%s</p>
// 		`
// 		finalNewPostRender := fmt.Sprintf(newPostSection, postTitle, postBody)
// 		w.Write([]byte(finalNewPostRender))

// 	}
// 	http.HandleFunc("/", dashboardHandler)
// 	http.HandleFunc("/add-new-post", receiveNewPostHandler)
// 	http.HandleFunc("/login", nil)
// 	http.HandleFunc("/logout", nil)
// 	http.HandleFunc("/refresh-token", nil)

// 	fmt.Println("Server is listening to port 3003 ...")
// 	log.Fatal(http.ListenAndServe(":3003", nil))
// }
