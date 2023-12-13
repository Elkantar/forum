package forum

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var userLog map[string]string
var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)
var user []map[string]string

type Post struct {
	ID       int
	Username string
	Title    string
	Text     string
	Image    string
	Likes    int
	Dislikes int
	Date     int
	Comment  map[string][]string
}
type RepPost struct {
	IDrep   int
	PostID  int
	UserRep string
	RepText string
}
type RowPost struct {
	Row      []Post
	Response []RepPost
}

func parseBrowser(userAgent string) string {
	// Liste de certains navigateurs courants
	browsers := map[string]string{
		"Chrome":  "Google Chrome",
		"Firefox": "Mozilla Firefox",
		"Safari":  "Apple Safari",
		"Edge":    "Microsoft Edge",
	}
	// Recherche du navigateur dans l'en-tête "User-Agent"
	for key, value := range browsers {
		if strings.Contains(userAgent, key) {
			return value
		}
	}
	// Si le navigateur n'est pas trouvé, retourne une valeur par défaut
	return "Navigateur inconnu"
}

func GenerateToken() uuid.UUID {
	Token, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	return Token
}

func Exportcomment(IDValue int) map[string][]string {
	IDValueSTR := strconv.Itoa(IDValue)
	db, err := sql.Open("sqlite3", "database/users.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Récupérer toutes les réponses de la table "reponse"
	rows, err := db.Query("SELECT * FROM reponse")
	if err != nil {
		panic(err)
	}
	// Créer une variable de type "map[string][]string[]" pour stocker les données
	reponses := make(map[string][]string)
	// Parcourir toutes les réponses et les stocker dans la variable "reponses"
	for rows.Next() {
		var repID int
		var postID, repUser, repText string
		err = rows.Scan(&repID, &postID, &repUser, &repText)
		if err != nil {
			panic(err)
		}
		// Stocker les données dans la variable "reponses"*
		if postID == IDValueSTR {
			reponses[repUser] = append(reponses[repUser], repText)
		}
	}
	return reponses
}

func importRep(user, text string, idPost int) {
	db, err1 := sql.Open("sqlite3", "database/users.db")
	if err1 != nil {
		return
	}
	defer db.Close()
	query := "INSERT INTO reponse (postID, repText , repUser) VALUES (?, ?, ?)"
	_, err := db.Exec(query, idPost, text, user)
	if err != nil {
		return
	}
}

func exportRep(IDvalue int) RowPost {
	datas := RowPost{}
	// datas := Post{}
	db, err1 := sql.Open("sqlite3", "database/users.db")
	if err1 != nil {
		panic(err1)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM reponse")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// Parcourir les lignes de résultat
	for rows.Next() {
		var IDrep int
		var UserRep string
		var RepText string
		var PostID string
		IDvalueInt := strconv.Itoa(IDvalue)
		// var userCreate string
		// Extraire les valeurs des colonnes de chaque ligne
		err := rows.Scan(&IDrep, &PostID, &UserRep, &RepText)
		if err != nil {
			panic(err)
		}
		if PostID == IDvalueInt {
			poststruc := RepPost{
				IDrep:   IDrep,
				UserRep: UserRep,
				RepText: RepText,
			}
			datas.Response = append(datas.Response, poststruc)
		}
	}
	// datareturn.Response = append(datareturn.Response, datas.Comment...)
	// fmt.Println(datas)
	return datas
}

func insertPost(title, text, image, User, img string) {
	db, err1 := sql.Open("sqlite3", "database/users.db")
	// User = utilisateur
	if err1 != nil {
		return
	}
	defer db.Close()
	query := "INSERT INTO posts (post_title, post_text, post_categorie, Username,image) VALUES (?, ?, ?, ?,?)"
	_, err := db.Exec(query, title, text, image, User, img)
	if err != nil {
		return
	}
}

func ReadCookie(r *http.Request, name string) (string, error) {
	userAgent := r.UserAgent()
	browser := parseBrowser(userAgent)
	fmt.Println("Navigateur: " + browser)
	// Read the cookie as normal.
	cookie, err := r.Cookie("logged")
	if err != nil {
		return "errCookie", err
	}
	// Return the decoded cookie value.
	return string(cookie.Value), nil
}

/*
Permet d'exporter les post de la base de donné en fonction de leur id et ainsi leur attribué une url unique
*/
func exportPostDetails(IdValue int) RowPost {
	datas := RowPost{}
	db, err1 := sql.Open("sqlite3", "database/users.db")
	if err1 != nil {
		panic(err1)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// Parcourir les lignes de résultat
	for rows.Next() {
		var postID int
		var postTitle string
		var postText string
		var postCategorie string
		var likes int
		var dislikes int
		var user string
		var images string
		// Extraire les valeurs des colonnes de chaque ligne
		err := rows.Scan(&postID, &postTitle, &postText, &postCategorie, &likes, &dislikes, &user, &images)
		if err != nil {
			panic(err)
		}
		if postID == IdValue {
			// datas2 = exportRep(postID)
			poststruc := Post{
				ID: postID,
				// Username: userCreate,
				Title:    postTitle,
				Text:     postText,
				Likes:    likes,
				Dislikes: dislikes,
				Username: user,
				Image:    images,
				Comment:  Exportcomment(postID),
			}
			datas.Row = append(datas.Row, poststruc)
		}
	}
	// fmt.Println(datas)
	return datas
}

func exportPost() RowPost {
	datas := RowPost{}
	db, err1 := sql.Open("sqlite3", "database/users.db")
	if err1 != nil {
		panic(err1)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// Parcourir les lignes de résultat
	for rows.Next() {
		var postID int
		var postTitle string
		var postText string
		var postCategorie string
		var likes int
		var dislikes int
		var user string
		var images string
		// Extraire les valeurs des colonnes de chaque ligne
		err := rows.Scan(&postID, &postTitle, &postText, &postCategorie, &likes, &dislikes, &user, &images)
		if err != nil {
			panic(err)
		}
		poststruc := Post{
			ID: postID,
			// Username: userCreate,
			Title:    postTitle,
			Text:     postText,
			Likes:    likes,
			Dislikes: dislikes,
			Username: user,
			Image:    images,
		}
		datas.Row = append(datas.Row, poststruc)
	}
	return datas
}

func exportPostlikedis(likePost []string) RowPost {
	datas := RowPost{}
	db, err1 := sql.Open("sqlite3", "database/users.db")
	if err1 != nil {
		panic(err1)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// Parcourir les lignes de résultat
	for rows.Next() {
		var postID int
		var postTitle string
		var postText string
		var postCategorie string
		var likes int
		var dislikes int
		var user string
		var images string
		// Extraire les valeurs des colonnes de chaque ligne
		err := rows.Scan(&postID, &postTitle, &postText, &postCategorie, &likes, &dislikes, &user, &images)
		if err != nil {
			panic(err)
		}
		for _, ch := range likePost {
			postIdstr := strconv.Itoa(postID)
			if ch == postIdstr {
				poststruc := Post{
					ID: postID,
					// Username: userCreate,
					Title:    postTitle,
					Text:     postText,
					Likes:    likes,
					Dislikes: dislikes,
					Username: user,
					Image:    images,
				}
				datas.Row = append(datas.Row, poststruc)
			}
		}
	}
	return datas
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	postId := r.FormValue("likepostid")
	fmt.Println("postId :", postId)
	str, _ := ReadCookie(r, "logged")
	var userid string
	userid = str
	db, err1 := sql.Open("sqlite3", "database/users.db")
	if err1 != nil {
		panic(err1)
	}
	defer db.Close()
	err := db.QueryRow("SELECT name FROM users WHERE status = ?", str).Scan(&userid)
	if err != nil {
		fmt.Print("db ")
		fmt.Println(err)
		http.Redirect(w, r, "/HomePage", http.StatusSeeOther)
		return
		// Gérer l'erreur
	}
	count := 0
	countverif := 0
	userlikeverif, _ := db.Query("SELECT * FROM userdislike")
	userlike, _ := db.Query("SELECT * FROM userlike")
	var ver bool
	ver = true
	for userlike.Next() {
		var userlikeID int
		var iduser string
		var idpost string
		errscan := userlike.Scan(&userlikeID, &iduser, &idpost)
		if errscan != nil {
			panic(errscan)
		}
		fmt.Println("----------------------------")
		fmt.Println(postId)
		fmt.Println(idpost)
		fmt.Println("----------------------------")
		if postId == idpost && userid == iduser {
			count++
		} else {
			count++
			countverif++
		}
		fmt.Println(count)
		fmt.Println(countverif)
	}
	if count > countverif {
		querymoins := "UPDATE posts SET likes = likes - 1 WHERE post_id = ?"
		_, errmoins := db.Exec(querymoins, postId)
		if errmoins != nil {
			panic(errmoins)
		}
		querymoinsuser := "DELETE FROM userlike WHERE postid = ? AND userid = ?"
		_, err1 := db.Exec(querymoinsuser, postId, userid)
		if err1 != nil {
			panic(err1)
		}
	} else {
		querymoinsuser := "INSERT INTO userlike (postid, userid) VALUES (?,?)"
		_, err1 := db.Exec(querymoinsuser, postId, userid)
		if err1 != nil {
			// fmt.Println("ici")
			panic(err1)
		}
		query := "UPDATE posts SET likes = likes + 1 WHERE post_id = ?"
		_, err := db.Exec(query, postId, userid)
		if err != nil {
			panic(err)
		}
		for userlikeverif.Next() {
			var userdislikeID int
			var iduser string
			var idpost string
			errscan := userlikeverif.Scan(&userdislikeID, &idpost, &iduser)
			if errscan != nil {
				panic(errscan)
			}
			if userid == iduser {
				ver = false
			}
		}
		if !ver {
			querylike := "UPDATE posts SET dislikes = dislikes - 1 WHERE post_id = ?"
			_, errmoins := db.Exec(querylike, postId)
			if errmoins != nil {
				panic(errmoins)
			}
			querymoinsuserdislike := "DELETE FROM userdislike WHERE postid = ? AND userid = ?"
			_, err1 := db.Exec(querymoinsuserdislike, postId, userid)
			if err1 != nil {
				panic(err1)
			}
		}
	}
	http.Redirect(w, r, "/detail?"+postId, http.StatusFound)
}

func veriflike(user string, db2 *sql.Rows) bool {
	for db2.Next() {
		var userdislikeID int
		var iduser string
		var idpost string
		errscan := db2.Scan(&userdislikeID, &idpost, &iduser)
		if errscan != nil {
			panic(errscan)
		}
		if iduser == user {
			return false
		}
	}
	return true
}

func DislikePost(w http.ResponseWriter, r *http.Request) {
	postId := r.FormValue("postid")
	fmt.Println("postId :", postId)
	str, _ := ReadCookie(r, "logged")
	var userid string
	userid = str
	db, err1 := sql.Open("sqlite3", "database/users.db")
	if err1 != nil {
		panic(err1)
	}
	defer db.Close()
	err := db.QueryRow("SELECT name FROM users WHERE status = ?", str).Scan(&userid)
	if err != nil {
		fmt.Print("db ")
		fmt.Println(err)
		http.Redirect(w, r, "/HomePage", http.StatusSeeOther)
		return
		// Gérer l'erreur
	}
	count := 0
	countverif := 0
	userveriflike, _ := db.Query("SELECT * FROM userlike")
	userlike, _ := db.Query("SELECT * FROM userdislike")
	var ver bool
	ver = true
	for userlike.Next() {
		var userdislikeID int
		var iduser string
		var idpost string
		errscan := userlike.Scan(&userdislikeID, &idpost, &iduser)
		if errscan != nil {
			panic(errscan)
		}
		fmt.Println("----------------------------")
		fmt.Println(postId)
		fmt.Println(idpost)
		fmt.Println("----------------------------")
		if postId == idpost && userid == iduser {
			count++
		} else {
			count++
			countverif++
		}
		fmt.Println(count)
		fmt.Println(countverif)
	}
	if count > countverif {
		querymoins := "UPDATE posts SET dislikes = dislikes - 1 WHERE post_id = ?"
		_, errmoins := db.Exec(querymoins, postId)
		if errmoins != nil {
			panic(errmoins)
		}
		querymoinsuser := "DELETE FROM userdislike WHERE postid = ? AND userid = ?"
		_, err1 := db.Exec(querymoinsuser, postId, userid)
		if err1 != nil {
			panic(err1)
		}
	} else {
		querymoinsuser := "INSERT INTO userdislike (postid, userid) VALUES (?,?)"
		_, err1 := db.Exec(querymoinsuser, postId, userid)
		if err1 != nil {
			fmt.Println("ici")
			panic(err1)
		}
		query := "UPDATE posts SET dislikes = dislikes + 1 WHERE post_id = ?"
		_, err := db.Exec(query, postId, userid)
		if err != nil {
			panic(err)
		}
		for userveriflike.Next() {
			var userdislikeID int
			var iduserver string
			var idpostver string
			errscan := userveriflike.Scan(&userdislikeID, &iduserver, &idpostver)
			if errscan != nil {
				panic(errscan)
			}
			if userid == iduserver {
				ver = false
			}
		}
		if !ver {
			querylike := "UPDATE posts SET likes = likes - 1 WHERE post_id = ?"
			_, errmoins := db.Exec(querylike, postId)
			if errmoins != nil {
				panic(errmoins)
			}
			querymoinsuserlike := "DELETE FROM userlike WHERE postid = ? AND userid = ?"
			_, err1 := db.Exec(querymoinsuserlike, postId, userid)
			if err1 != nil {
				panic(err1)
			}
		}
	}
	http.Redirect(w, r, "/detail?"+postId, http.StatusFound)
}

func userExists(db *sql.DB, username string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func createUser(username, password string, name string) error {
	// Vérifier si l'utilisateur existe déjà
	db, err := sql.Open("sqlite3", "database/users.db")
	if err != nil {
		return err
	}
	defer db.Close()
	exists, err := userExists(db, username)
	if err != nil {
		// return err
	}
	if exists {
		return fmt.Errorf("user '%s' already exists", username)
	}
	// Hasher le mot de passe avec bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Insérer l'utilisateur dans la base de données
	_, err = db.Exec("INSERT INTO users (email, password, name) VALUES (?, ?, ?)", username, hashedPassword, name)
	if err != nil {
		return err
	}
	return nil
}

func authenticateUser(username, password string) (bool, error) {
	// Get the hashed password from the database
	db, err := sql.Open("sqlite3", "database/users.db")
	if err != nil {
		return false, err
	}
	defer db.Close()
	var hashedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE email = ?", username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		} else {
			return false, err
		}
	}
	// Compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}

func Signup(email, passwords, name string) bool {
	var username, password string
	username = email
	password = passwords
	err := createUser(username, password, name)
	if err != nil {
		// log.Fatal(err)
		return false
	}
	fmt.Println("Signup successful!")
	return true
}

func login(email, passwords string) bool {
	var username, password string
	username = email
	password = passwords
	authenticated, err := authenticateUser(username, password)
	if err != nil {
		log.Fatal(err)
	}
	if authenticated {
		fmt.Println("Login successful!")
		return true
	} else {
		fmt.Println("Login failed!")
		return false
	}
	return false
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	// Supprimer le cookie en définissant sa valeur sur une chaîne vide et une date d'expiration antérieure
	expiration := time.Now().Add(-time.Hour)
	cookie := http.Cookie{Name: "user", Value: "", Expires: expiration}
	http.SetCookie(w, &cookie)
	// Rediriger l'utilisateur vers la page de connexion
	http.Redirect(w, r, "/", http.StatusFound) // faire la redirection sur la page de log
}
