makefile:

dev: 
	go build && ./Craigscraper

db:
	@mkdir -p db
	@rm -f db/listings.db 
	@touch db/listings.db
	@echo "Created New Database: db/listings.db"