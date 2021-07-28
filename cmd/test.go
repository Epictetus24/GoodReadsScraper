package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type BookDeets struct {
	Name     string
	CoverURL string
	BookID   string
	Author   string
	Bookurl  string
	Status   string
	Blurb    string
	Rating   string
}

type UserLibrary struct {
	userid   string
	BookList []BookDeets
}

func main() {

	url := "https://www.goodreads.com/user/show/104959477"

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)

	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		os.Exit(1)
	}

	var UserLib UserLibrary

	// Save each .post-title as a list
	doc.Find(".secondcol-top").Each(func(i int, s *goquery.Selection) {
		var book BookDeets
		title := s.Find(".bookTitle")
		book.Name = title.Text()
		Author := s.Find(".authorName")
		book.Author = Author.Text()

		url, _ := title.Attr("href")
		book.Bookurl = "https://goodreads.com" + url

		//visit the book page, get the blurb, rating and cover.

		bkresp, _ := http.Get(book.Bookurl)
		defer bkresp.Body.Close()

		bkdoc, err := goquery.NewDocumentFromReader(bkresp.Body)
		if err != nil {
			os.Exit(1)
		}

		cover := bkdoc.Find(".editionCover>img")
		book.CoverURL, _ = cover.Attr("src")

		UserLib.BookList = append(UserLib.BookList, book)

	})

	fmt.Println(UserLib)

}
