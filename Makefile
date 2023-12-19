makefile:

dev: 
	go build && ./Craigscraper

db:
	rm db/listings.db && touch db/listings.db