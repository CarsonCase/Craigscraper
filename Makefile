makefile:

dev: 
	go build && ./Craigscraper

db:
	rm -f db/listings.db && touch db/listings.db