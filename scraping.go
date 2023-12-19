package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"golang.org/x/net/html"
)

type Context struct {
	RequestCount int
	InProgress   int
	Complete     int
	Proxies      []string
}

func (c *Context) incrementInProgress() {
	c.InProgress++
}

func (c *Context) incrementComplete() {
	c.Complete++
}

// scrapePage scrapes the specified URL for listings.
//
// The function takes a URL, a counter, a listing channel, and a page channel
// as arguments. The counter is used to track the number of listings that are
// currently being scraped. The listing channel is used to send listings to
// the main function. The page channel is used to signal that a page has been
// scraped.
// There are too layers of callback functions, 1 to find li elements, and another to find the price data WITHIN that element
func scrapePage(url string, context *Context, listingChan chan Listing, pageChan chan bool) {
	c := colly.NewCollector()
	c.OnHTML("li.cl-static-search-result", func(e *colly.HTMLElement) {
		listing := Listing{}
		title := e.Attr("title")
		price := e.ChildText(".price")
		link := e.ChildAttr("a", "href")
		listing.Title = title
		listing.Price = price
		listing.Link = link

		listingChan <- listing
	})
	rp, err := proxy.RoundRobinProxySwitcher("http://p.webshare.io:9999/")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)
	c.Visit(url)

	pageChan <- true
}

func (c *Context) getRespWithProxy(getURL string) (resp *http.Response, err error) {
	proxyUrl, err := url.Parse("http://p.webshare.io:9999/")
	proxy := http.ProxyURL(proxyUrl)
	transport := &http.Transport{Proxy: proxy}
	client := &http.Client{Transport: transport}
	req, err := http.NewRequest("GET", getURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	c.RequestCount++
	return resp, nil
}

// getHTMLPage gets the HTML page for the specified URL.
//
// The function returns an HTML document object.
func (c *Context) getHTMLPage(url string) *html.Node {
	response, err := c.getRespWithProxy(url)
	if err != nil {
		log.Fatal("Get Error: ", err, "\nRequestCount: ", c.RequestCount)
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
