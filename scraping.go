package main

import (
	"log"
	"net/http"

	"golang.org/x/net/html"
)

// scrapePage scrapes the specified URL for listings.
//
// The function takes a URL, a counter, a listing channel, and a page channel
// as arguments. The counter is used to track the number of listings that are
// currently being scraped. The listing channel is used to send listings to
// the main function. The page channel is used to signal that a page has been
// scraped.
// There are too layers of callback functions, 1 to find li elements, and another to find the price data WITHIN that element
func scrapePage(url string, counter *Counter, listingChan chan Listing, pageChan chan bool) {
	doc := getHTMLPage(url)
	listing := Listing{}
	findHTML(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			for _, a := range n.Attr {
				if a.Key == "title" {
					counter.incrementInProgress()
					listing = Listing{Title: a.Val}
					break
				}
			}

			findHTML(n, func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "div" {
					for _, a := range n.Attr {
						if a.Key == "class" && a.Val == "price" {
							counter.incrementComplete()
							listing.Price = n.FirstChild.Data
							listingChan <- listing
							break
						}
					}
				}
			})
		}
	})
	pageChan <- true
}

// getHTMLPage gets the HTML page for the specified URL.
//
// The function returns an HTML document object.
func getHTMLPage(url string) *html.Node {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Get Error")
	}

	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		log.Fatal("Read Response Error")
	}
	return doc
}

// findHTML recursively traverses the HTML tree, calling the `foo` function
// on each node.
//
// The `foo` function takes a single argument, the HTML node. It can be used
// to perform any desired operations on the node, such as extracting text.
func findHTML(n *html.Node, foo func(n *html.Node)) {
	foo(n)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findHTML(c, foo)
	}
}
