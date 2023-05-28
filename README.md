# A backend service demonstrating REST APIs with SQL (sqlite) in Golang

```
$ cd $HOME/go/src/Golang-REST-SQL
$ go run main.go
```

go to Postman to test CURD REST APIs:


| Method   |      URL      |   Payload  |
|----------|:-------------:|:----------:|
| GET      | http://localhost:8000/books      |    |
| GET      | http://localhost:8000/books/:id  |    |
| POST     | http://localhost:8000/books      | {"title": "Wonderland", "author_id": 2, "published": 2023 }    |
| PUT      | http://localhost:8000/books/:id  |    |
| DELETE   | http://localhost:8000/books/:id  |    |

