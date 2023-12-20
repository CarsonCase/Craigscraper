makefile:

dev: 
	go build && ./Craigscraper

db:
	rm -f db && touch db/listings.db