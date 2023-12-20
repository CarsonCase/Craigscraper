makefile:

dev: 
	go build && ./Craigscraper

db:
	rm -r db && touch db/listings.db