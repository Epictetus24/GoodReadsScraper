package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func init() {
	os.Mkdir("Bookcovers", 444)
}

func GetBookCover(bookid string) {
	coverurl := ""
	http.Get(coverurl)
}

type BookDeets struct {
	Name     string
	CoverURL string
	BookID   string
	Author   string
	Bookurl  string
	Status   string
	Blurb    string
}

type UserLibrary struct {
	Userid   string
	BookList []BookDeets
}

func ParseBookDeets(userid string, resp http.Response) UserLibrary {

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		os.Exit(1)
	}

	var UserLib UserLibrary
	var book BookDeets
	UserLib.Userid = userid

	//find the book attributes
	doc.Find(".secondcol-top").Each(func(i int, s *goquery.Selection) {

		title := s.Find(".bookTitle")
		book.Name = title.Text()
		Author := s.Find(".authorName")
		book.Author = Author.Text()

		url, _ := title.Attr("href")
		book.Bookurl = "https://goodreads.com" + url

		//visit the book page and get the blurb and cover.
		bkresp, _ := http.Get(book.Bookurl)
		defer bkresp.Body.Close()

		bkdoc, err := goquery.NewDocumentFromReader(bkresp.Body)
		if err != nil {
			os.Exit(1)
		}

		cover := bkdoc.Find(".editionCover>img")
		book.CoverURL, _ = cover.Attr("src")

		book.Blurb = "blah blah"

		UserLib.BookList = append(UserLib.BookList, book)

	})
	fmt.Println(UserLib)

	return UserLib

}

func GetProfile(userid string) http.Response {

	url := "https://www.goodreads.com/user/show/" + userid

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)

	}

	return *resp
}

const mainpage = `<!doctype html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, user-scalabe=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">


<title>GoodReads Shelf</title>
</head>
	<h1>Viewing GoodReads Profile</h1>

	<h2>Currently Reading:</h2>
	<table><tr>
	<th>Cover</th>
	<th>Details</th>
	<th>Blurb</th>
	</tr>
	{{range .BookList}}
	<tr>
	<td><img src="{{.CoverURL}} "  width="153" /></td>
	<td>
	<b>Name {{.Name}}</b><br>
	Author {{.Author}}<br>
	<a href="{{.Bookurl}}">"Click for more details</a>
	</td>
	<td><span style="white-space: pre-wrap">{{.Blurb}}<span></td>
	</tr>
	{{end}}</table>



	<h3>Change User:</h3>

	<form action="/" method="POST">
	<label for="userid">set userid</label>
	<input type="text" name="userid" id="userid" placeholder="UserID from goodreads profile URL e.g. https://www.goodreads.com/user/show/234523434-username/">

	 <input type="submit">
	 </form>

	</html>`

func grhandler(w http.ResponseWriter, r *http.Request) {

	/*
		tpl, err := template.ParseFiles("templates/index.gohtml")
		if err != nil {
			log.Fatal(err.Error())
		}
	*/

	tpl := template.Must(template.New("index").Parse(mainpage))

	var Userid string

	if r.FormValue("userid") == "" {
		Userid = "104959477"

	} else {

		Userid = r.FormValue("userid")

	}

	resp := GetProfile(Userid)

	UserLib := ParseBookDeets(Userid, resp)

	err2 := tpl.Execute(w, UserLib)
	if err2 != nil {
		fmt.Println("Issue with template")
		log.Fatal(err2.Error())
	}

}

func main() {

	http.HandleFunc("/", grhandler)
	http.ListenAndServe(":8080", nil)

}
