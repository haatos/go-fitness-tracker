build:
	go build -ldflags "-s -w" -o bin/fitnesstracker main.go

build-tw:
	npx tailwindcss -o ./static/css/tw.css --minify

deploy: build-tw build
	scp -i $(id) -r 'migrations' ubuntu@$(ip):/opt/fitness-tracker/
	scp -i $(id) -r 'static' ubuntu@$(ip):/opt/fitness-tracker/
	scp -i $(id) -r 'templates' ubuntu@$(ip):/opt/fitness-tracker/
	scp -i $(id) './bin/fitnesstracker' ubuntu@$(ip):/opt/fitness-tracker/
	scp -i $(id) '.env' ubuntu@$(ip):/opt/fitness-tracker/

migrate-down:
	goose -dir migrations sqlite3 ./db.sqlite3 down

migrate-down-to:
	goose -dir migrations sqlite3 ./db.sqlite3 down-to $(target)

migrate-up:
	goose -dir migrations sqlite3 ./db.sqlite3 up

migrate-up-to:
	goose -dir migrations sqlite3 ./db.sqlite3 up-to $(target)
