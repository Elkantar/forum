package main

import (
	"fmt"
	"log"
	"net/http"

	gohtml "forum/gohtml"

	_ "github.com/mattn/go-sqlite3"
)

/*var db *sql.DB*/
func main() {
	fmt.Println("Server launch at : http://localhost:8050/HomePage ")
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	// image client
	img := http.FileServer(http.Dir("uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", img))

	// image serveur
	imgs := http.FileServer(http.Dir("images"))
	http.Handle("/images/", http.StripPrefix("/images/", imgs))

	jas := http.FileServer(http.Dir("js"))
	http.Handle("/js/", http.StripPrefix("/js/", jas))
	http.HandleFunc("/main", gohtml.Server)
	http.HandleFunc("/CreatePost", gohtml.CreatePost)
	http.HandleFunc("/detail", gohtml.Detail)
	http.HandleFunc("/Post", gohtml.Posts)
	http.HandleFunc("/HomePage", gohtml.Home)
	http.HandleFunc("/Likepost", gohtml.LikePage)
	http.HandleFunc("/Dislikepost", gohtml.DislikePage)
	http.HandleFunc("/Like", gohtml.LikePost)
	http.HandleFunc("/Dislike", gohtml.DislikePost)

	//
	err := http.ListenAndServe(":8050", nil)
	if err != nil {
		log.Fatal(err)
	}
}
