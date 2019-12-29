package docparser

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"telegram/internal/datastructures"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var wg sync.WaitGroup

func filePathWalkDir(root string, c chan string) {
	var counter int = 0
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		counter++
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			c <- path
			//TODO(savsunenko): DONT FORGET TO FIX ME

		}
		return nil
	})
	close(c)
	wg.Done()

}

func processDoc(docStream chan datastructures.Document, c chan string) {

	for i := range c {

		file, err := os.Open(i)
		if err != nil {
			log.Fatal(err)
		}

		doc, err := html.Parse(file)
		if err != nil {
			fmt.Println(err)
		}

		var f func(*html.Node)
		var document datastructures.Document
		document = datastructures.Document{FullText: ""}

		//TODO Make external function from this
		f = func(n *html.Node) {
			if n.DataAtom == atom.Meta && n.Attr[0].Val == "og:title" {
				document.Header = n.Attr[1].Val
			}
			if n.DataAtom == atom.Meta && n.Attr[0].Val == "og:site_name" {
				document.WebsiteName = n.Attr[1].Val
			}
			if n.DataAtom == atom.Meta && n.Attr[0].Val == "og:url" {
				u, err := url.Parse(n.Attr[1].Val)
				if err != nil {
					//TODO: REMOVE ME
					log.Fatal(err)
				}
				document.WebsiteAddress = u.Hostname()
			}
			if n.DataAtom == atom.Meta && n.Attr[0].Val == "article:published_time" {
				document.Time = n.Attr[1].Val
			}

			if n.Parent != nil && n.Parent.DataAtom == atom.H1 {
				document.H1 = n.Data
			}
			if n.Parent != nil && n.Parent.DataAtom == atom.H2 {
				document.H2 = n.Data
			}

			if n.Type == html.TextNode && n.Parent.Data == "p" {
				document.FullText += n.Data + " "
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}

		f(doc)
		document.Path = i
		document.FileName = filepath.Base(i)
		docStream <- document
		file.Close()
	}
	close(docStream)
	wg.Done()
}

// ParseDoc processes files in path
func ParseDoc(path string, result []datastructures.Document) []datastructures.Document {
	docStream := make(chan datastructures.Document, 100)
	individualFiles := make(chan string, 100)

	wg.Add(1)
	go filePathWalkDir(path, individualFiles)
	wg.Add(1)
	go processDoc(docStream, individualFiles)

	for i := range docStream {
		result = append(result, i)
	}
	wg.Wait()
	return result
}
