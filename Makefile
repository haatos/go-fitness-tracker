build:
	go build -ldflags "-s -w" -o bin/fitnesstracker main.go

migrate-down:
	goose -dir migrations sqlite3 ./db.sqlite3 down

migrate-down-to:
	goose -dir migrations sqlite3 ./db.sqlite3 down-to $(target)

migrate-up:
	goose -dir migrations sqlite3 ./db.sqlite3 up

migrate-up-to:
	goose -dir migrations sqlite3 ./db.sqlite3 up-to $(target)
