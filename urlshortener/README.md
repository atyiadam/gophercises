URL Shortener

source: https://github.com/gophercises/urlshort

Solution, including all the bonus tasks as well.

```
$ docker-compose up -d
$ goose -dir migrations/ postgres "user=user password=password dbname=urlshortener_db host=localhost port=5432 sslmode=disable" up
$ go run main/main.go
```
