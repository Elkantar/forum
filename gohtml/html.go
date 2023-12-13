package forum

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var utilisateur string

type display struct {
	displayed bool
}

var (
	templates      = template.Must(template.ParseFiles("main.html"))
	templates1     = template.Must(template.ParseFiles("./html/HomePage.html"))
	templates2     = template.Must(template.ParseFiles("./html/detail.html"))
	templates3     = template.Must(template.ParseFiles("./html/CreatePost.html"))
	templatelike   = template.Must(template.ParseFiles("./html/like.html"))
	templatedislke = template.Must((template.ParseFiles("./html/dislike.html")))
)

var mailCookie string

func cookie(w http.ResponseWriter, r *http.Request, name string) {
	userLog = make(map[string]string)
	cookie, err := r.Cookie("logged")
	if err != nil {

		Token := GenerateToken()
		TokenString := Token.String()
		userLog[TokenString] = utilisateur

		user = append(user, userLog)
		// If the cookie doesn't exist, create a new one
		if utilisateur != "" {
			cookie = &http.Cookie{
				Name:     "logged",
				Value:    TokenString,
				Expires:  time.Now().Add(1 * time.Hour), // Expires in 2 hours
				HttpOnly: false,
			}
			// stocke le cookie dans la database

			db, err1 := sql.Open("sqlite3", "database/users.db")
			if err1 != nil {
				return
			}
			defer db.Close()

			_, err3 := db.Exec("UPDATE users SET status = ? WHERE name = ?", TokenString, utilisateur)
			if err3 != nil {
				fmt.Println(err3)
				return
			}

		} else {
			cookie = &http.Cookie{
				Name:     "unlogged",
				Value:    fmt.Sprintf("%d", time.Now().Unix()),
				Expires:  time.Now().Add(1 * time.Hour), // Expires in 2 hours
				HttpOnly: false,
			}
		}
		http.SetCookie(w, cookie)
	}
	fmt.Println("Cookie value: ", cookie.Value)
}

func Detail(w http.ResponseWriter, r *http.Request) {
	//----------Token/Cookie--------------//
	str, _ := ReadCookie(r, "logged")
	fmt.Println(str)
	// debut
	//------------Split l'url pour recuperer l'id du post et ainsi créé un url uniique pour chaque post-------------//
	var urlSplit2 string
	url := r.URL.String()
	urlSplit := strings.Split(url, "/detail?")
	for _, ch := range urlSplit {
		if ch != "" {
			urlSplit2 = ch
		}
	}
	atoiUrl, _ := strconv.Atoi(urlSplit2)
	datas := exportPostDetails(atoiUrl)
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	//--------------Traitement erreur 404 et 500--------------//
	if r.URL.Path != "/" && r.URL.Path != "/HomePage" && r.URL.Path != "/detail" {
		http.Error(w, "Error 404: file not found.", http.StatusNotFound)
		return
	}
	_, errt := template.ParseFiles("./html/detail.html")
	if errt != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error 500: Internal Server Error"))
		log.Println(http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case "GET":
		//--------------lancement du site--------------//
		templates2.ExecuteTemplate(w, "detail.html", &datas)
		// templates21.ExecuteTemplate(w, "detailComment.html", &datas2)
		// fmt.Println(RowPost)
	case "POST":
		templates2.ExecuteTemplate(w, "detail.html", &datas)
		// templates21.ExecuteTemplate(w, "detailComment.html", &datas2)
		var repUser string
		db, err1 := sql.Open("sqlite3", "database/users.db")
		// User = utilisateur
		if err1 != nil {
			return
		}
		str, _ = ReadCookie(r, "logged")
		err := db.QueryRow("SELECT name FROM users WHERE status = ?", str).Scan(&repUser)
		if err != nil {
			fmt.Println("reppost: ", err)
			return
		}
		r.ParseForm()
		if r.FormValue("repText") != "" {
			// Récupérer les valeurs de la première formulaire
			repText := r.FormValue("repText")
			importRep(repUser, repText, atoiUrl)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func Posts(w http.ResponseWriter, r *http.Request) {
	//----------Token/Cookie--------------//
	// debut
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	//--------------Traitement erreur 404 et 500--------------//
	if r.URL.Path != "/" && r.URL.Path != "/HomePage" && r.URL.Path != "/Post" {
		http.Error(w, "Error 404: file not found.", http.StatusNotFound)
		return
	}
	_, errt := template.ParseFiles("./html/posts.html")
	if errt != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error 500: Internal Server Error"))
		log.Println(http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case "GET":
		//--------------lancement du site--------------//
		templates1.ExecuteTemplate(w, "./html/posts.html", display{displayed: false})
	case "POST":
		templates1.ExecuteTemplate(w, "posts.html", display{displayed: false})
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var savePath string
	// debut
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	//--------------Traitement erreur 404 et 500--------------//
	if r.URL.Path != "/" && r.URL.Path != "/HomePage" && r.URL.Path != "/CreatePost" {
		http.Error(w, "Error 404: file not found.", http.StatusNotFound)
		return
	}

	_, errt := template.ParseFiles("./html/CreatePost.html")
	if errt != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error 500: Internal Server Error"))
		log.Println(http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "GET":
		//--------------lancement du site--------------//
		templates3.ExecuteTemplate(w, "CreatePost.html", display{displayed: false})
	case "POST":
		str, _ := ReadCookie(r, "logged")
		var users string
		users = str
		db, err1 := sql.Open("sqlite3", "database/users.db")
		if err1 != nil {
			return
		}
		defer db.Close()
		err := db.QueryRow("SELECT name FROM users WHERE status = ?", str).Scan(&users)
		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/HomePage", http.StatusSeeOther)
			return
			// Gérer l'erreur
		}
		templates3.ExecuteTemplate(w, "CreatePost.html", display{displayed: false})

		PostTitle := r.FormValue("NewPostTitle")
		PostText := r.FormValue("NewPostText")
		Categorie := r.FormValue("Categorie")

		fmt.Println(PostTitle, PostText, Categorie)

		uploadChecker := r.FormValue("checker")
		fmt.Println(uploadChecker)
		if uploadChecker != "off" {
			// Handle image upload
			file, fileHeader, err := r.FormFile("NewPostImage")
			if err != nil {
				fmt.Println(err)
				// Handle error when no image is uploaded
				// Display an error message to the user
				fmt.Fprintf(w, "Error uploading image:no image or image is too small < 250 Ko like you")
				return
			}
			defer file.Close()

			// Check image size
			maxFileSize := int64(20 * 1024 * 1024) // 20 MB
			if fileHeader.Size > maxFileSize {
				// Image size exceeds the limit, inform the user
				fmt.Fprintf(w, "Error: Image is too large. Max size is 20 MB.")
				return
			}

			// Check image file type
			fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
			if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" && fileExt != ".gif" {
				// Unsupported file type, inform the user
				fmt.Fprintf(w, "Error: Unsupported image file type.")
				return
			}

			// Create a unique filename for the image
			// You can use a library like uuid or a timestamp to generate unique filenames
			uniqueFilename := generateUniqueFilename(fileExt)

			// Create a file on the server to save the image
			savePath = "./uploads/" + uniqueFilename
			savedFile, err := os.Create(savePath)
			if err != nil {
				// Handle error when unable to create the file
				fmt.Fprintf(w, "Error saving image: %v", err)
				return
			}
			defer savedFile.Close()
			fmt.Println(savePath)
			// Copy the uploaded image data to the server file
			_, err = io.Copy(savedFile, file)
			if err != nil {
				// Handle error when unable to copy image data
				fmt.Fprintf(w, "Error saving image: %v", err)
				return
			}
			PostTitle = r.FormValue("NewPostTitle")
			PostText = r.FormValue("NewPostText")
			Categorie = r.FormValue("Categorie")
			fmt.Println(PostTitle, PostText, Categorie)

			if PostText != "" && PostTitle != "" {
				insertPost(PostTitle, PostText, Categorie, users, savePath)
			}
		} else {
			if PostText != "" && PostTitle != "" {
				insertPost(PostTitle, PostText, Categorie, users, savePath)
			}
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

// Function to generate a unique filename for the uploaded image
func generateUniqueFilename(extension string) string {
	// You can implement your logic here to generate a unique filename
	// For example, you can use a timestamp or a random string generator
	// to create a unique filename and ensure it doesn't clash with existing files.
	// Here's a simple example using a timestamp:
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d%s", timestamp, extension)
}

func Home(w http.ResponseWriter, r *http.Request) {
	//----------Token/Cookie--------------//
	cookie(w, r, mailCookie)
	//	Exportcomment()
	datas := exportPost()
	// debut
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	//--------------Traitement erreur 404 et 500--------------//
	// verif(w, r)
	if r.URL.Path != "/" && r.URL.Path != "/HomePage" {
		http.Error(w, "Error 404: file not found.", http.StatusNotFound)
		return
	}
	_, errt := template.ParseFiles("./html/HomePage.html")
	if errt != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error 500: Internal Server Error"))
		log.Println(http.StatusInternalServerError)
		return
	}
	// datas := exportPost()
	switch r.Method {
	case "GET":
		//--------------lancement du site--------------//
		err2 := templates1.Execute(w, &datas)
		if err2 != nil {
			log.Fatal(err2)
		}
	case "POST":
		// fmt.Println(datas)
		var email, password, name string
		err2 := templates1.Execute(w, &datas)
		if err2 != nil {
			log.Fatal(err2)
		}
		if r.FormValue("email") != "" && r.FormValue("password") != "" && r.FormValue("name") != "" {
			// Récupérer les valeurs de la première formulaire
			email = r.FormValue("email")
			password = r.FormValue("password")
			name = r.FormValue("name")
			// fmt.Printf("Email: %s, Password: %s\n", email, password)
			Signup(email, password, name)
		} else if r.FormValue("Sign-in-mail") != "" && r.FormValue("Sign-in-password") != "" {
			// Récupérer les valeurs de la deuxième formulaire
			email = r.FormValue("Sign-in-mail")
			password = r.FormValue("Sign-in-password")
			mailCookie = email
			// fmt.Printf("Email: %s, Password: %s\n", email, password)
			login(email, password)
			cookie(w, r, mailCookie)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func Server(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/main" {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	//--------------Traitement erreur 404 et 500--------------//
	//verif(w, r)
	_, err := template.ParseFiles("main.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error 500: Internal Server Error"))
		log.Println(http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case "GET":
		//--------------lancement du site--------------//
		templates.ExecuteTemplate(w, "main.html", display{displayed: false})
	case "POST":
		var email, password, name string
		if r.FormValue("email") != "" && r.FormValue("password") != "" && r.FormValue("name") != "" {
			// Récupérer les valeurs de la première formulaire
			email = r.FormValue("email")
			password = r.FormValue("password")
			name = r.FormValue("name")
			// fmt.Printf("Email: %s, Password: %s\n", email, password)
			verif := Signup(email, password, name)
			if !verif {
				dataSignup := "This user already exist"
				// fmt.Fprint(w, dataSignup)
				templates.ExecuteTemplate(w, "main.html", dataSignup)
				return
			}
		} else if r.FormValue("Sign-in-mail") != "" && r.FormValue("Sign-in-password") != "" {
			// Récupérer les valeurs de la deuxième formulaire
			email = r.FormValue("Sign-in-mail")
			mailCookie = email
			password = r.FormValue("Sign-in-password")
			// fmt.Printf("Email: %s, Password: %s\n", email, password)
			verif := login(email, password)
			if verif {
				utilisateur = name
				http.Redirect(w, r, "/HomePage", http.StatusSeeOther)
				db, err := sql.Open("sqlite3", "./database/users.db")
				if err != nil {
					fmt.Println("errorZ")
				}
				defer db.Close()
				searchEmail := "SELECT email, name FROM users"
				check, err := db.Query(searchEmail)
				if err != nil {
					log.Fatal(err)
				}
				for check.Next() {
					var ameil string
					var name string
					err := check.Scan(&ameil, &name)
					if err != nil {
						log.Fatal(err)
					}
					if email == ameil {
						fmt.Println(name)
						utilisateur = name
					}
				}
				return
			} else {
				datasignin := "Error during login"
				templates.ExecuteTemplate(w, "main.html", datasignin)
				return
			}
		}
		templates.ExecuteTemplate(w, "main.html", display{displayed: false})
		fmt.Fprintf(w, "")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func LikePage(w http.ResponseWriter, r *http.Request) {
	str, _ := ReadCookie(r, "logged")
	var userid string
	var postlike []string
	// userid = str

	if str == "" {
		http.Redirect(w, r, "/HomePage", http.StatusSeeOther)
	}

	db, err1 := sql.Open("sqlite3", "database/users.db")
	if err1 != nil {
		panic(err1)
	}

	err := db.QueryRow("SELECT name FROM users WHERE status = ?", str).Scan(&userid)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/HomePage", http.StatusSeeOther)
		return
		// Gérer l'erreur
	}

	querylike := "SELECT * FROM userlike"

	like, errlike := db.Query(querylike)
	if errlike != nil {
		panic(errlike)
	}

	for like.Next() {
		var likeid string
		var postid string
		var userpostid string

		errlike = like.Scan(&likeid, &userpostid, &postid)
		if errlike != nil {
			panic(errlike)
		}

		if userpostid == userid {
			postlike = append(postlike, postid)
		}
	}

	fmt.Println(postlike)

	datas := exportPostlikedis(postlike)
	templatelike.ExecuteTemplate(w, "like.html", &datas)
}

func DislikePage(w http.ResponseWriter, r *http.Request) {
	str, _ := ReadCookie(r, "logged")
	if str == "" {
		http.Redirect(w, r, "/HomePage", http.StatusSeeOther)
	}
	var userid string
	var postdislike []string
	// userid = str
	db, err1 := sql.Open("sqlite3", "database/users.db")
	if err1 != nil {
		panic(err1)
	}

	err := db.QueryRow("SELECT name FROM users WHERE status = ?", str).Scan(&userid)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/HomePage", http.StatusSeeOther)
		return
		// Gérer l'erreur
	}

	querydislike := "SELECT * FROM userdislike"

	dislike, errdislike := db.Query(querydislike)
	if errdislike != nil {
		panic(errdislike)
	}

	for dislike.Next() {
		var postiddislike string
		var userpostiddislike string
		var dislikeid string
		errdislike = dislike.Scan(&dislikeid, &postiddislike, &userpostiddislike)
		if errdislike != nil {
			panic(errdislike)
		}

		if userpostiddislike == userid {
			postdislike = append(postdislike, postiddislike)
		}
	}

	fmt.Println(postdislike)

	datas := exportPostlikedis(postdislike)

	templatedislke.ExecuteTemplate(w, "dislike.html", &datas)
}
