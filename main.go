package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// searchCity searches the specified city for listings.
//
// The function takes a city URL, a number of pages to search, a counter, a
// listing channel, and a done channel as arguments. The counter is used to
// track the number of listings that are currently being scraped. The listing
// channel is used to send listings to the main function. The done channel is
// used to signal that a city has been scraped.
func searchCity(cityUrl string, pagesToSearch int, context *Context, listings chan Listing, done chan bool) {

	doneScrapingPage := make(chan bool)
	for i := 0; i <= pagesToSearch; i++ {
		go scrapePage(cityUrl+"/search/cta#search=1~gallery~"+string(i)+"~0", context, listings, doneScrapingPage)
	}

	for i := 0; i <= pagesToSearch; i++ {
		<-doneScrapingPage
		fmt.Println("Finished page:\t", i)
	}
	done <- true

}

// getCities gets the cities from the specified URL.
//
// The function takes a URL as an argument. It returns a slice of strings,
// where each string is the URL of a city.
func getCities(url string, context *Context) (cities []string) {
	doc := context.getHTMLPage(url)

	findHTML(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {

			for _, a := range n.Attr {
				if a.Key == "href" {
					cities = append(cities, a.Val)
					break
				}
			}
		}
	})
	return cities
}

// storeValues stores the listings in the specified database.
//
// The function takes a listing channel, a done channel, and a database URL
// as arguments. The listing channel is used to receive listings from the
// scraping functions. The done channel is used to signal that all listings
// have been stored. The database URL is the URL of the database to store the
// listings in.
func storeValues(lc chan Listing, done chan bool, db *sql.DB) {
	var wg sync.WaitGroup
	for {
		val, ok := <-lc

		if !ok {
			break
		}
		wg.Add(1)
		err := insertListing(db, val)
		wg.Done()
		if err != nil {
			log.Fatal(err)
		}
	}
	wg.Wait()
	done <- true
}

// main is the entry point for the program.
//
// The function does the following:
//
//  1. Gets the cities from the specified URL.
//  2. Creates a counter to track the number of listings that are currently
//     being scraped.
//  3. Creates a listing channel to send listings to the main function.
//  4. Creates a done channel to signal that all listings have been scraped.
//  5. Creates a goroutine to store the listings in the database.
//  6. Creates a goroutine to search each city for listings.
//  7. Waits for all goroutines to finish.
//  8. Prints the total number of listings that were scraped.
func main() {
	fmt.Println("Setting up server")
	db := SetupDB()

	fmt.Println("Server setup complete")

	pagesToScan := 4
	batchSize := 5

	context := Context{}
	context.SetUpScraper("proxies_list.txt")
	cities := getCities("https://geo.craigslist.org/iso/us", &context)[1:414][0:125]
	startTime := time.Now()
	lc := make(chan Listing)
	done := make(chan bool)

	go storeValues(lc, done, db)

	cityStatus := make(chan bool)
	for i := 0; i < len(cities); i += batchSize {

		for j := i; j < i+batchSize; j++ {
			city := cities[j]
			fmt.Println("Searching:\t" + city)
			go searchCity(city, pagesToScan, &context, lc, cityStatus)
		}

		for j := i; j < i+batchSize; j++ {
			city := cities[j]
			<-cityStatus
			fmt.Println("Finished city:\t", city)
		}

	}

	// todo. If cities size isn't divisible by batch size, do the last cities

	close(lc)
	fmt.Println("Scraping Completed in: ", time.Since(startTime), " seconds!")
	<-done
	fmt.Println("Completed in: ", time.Since(startTime), " seconds!")
	fmt.Println("Print Complete\t#", context.Complete, "elements printed")
}
