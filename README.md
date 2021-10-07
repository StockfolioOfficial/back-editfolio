# back-editfolio
나중에 설명문 누가 작성좀 해줘

# How to Start

## requires
- [make](https://www.gnu.org/software/make/)
- [golang 1.17 🔺](https://golang.org/)
- [mysql 5 🔺](https://www.mysql.com/)

## Commands
```bash
# pwd
.../back-editfolio
# make init
# make generate
... process ...
# go run .
```

# Used

### HTTP Router
[Echo Framework](https://echo.labstack.com/)

### ORM
[GORM](https://gorm.io/)

### Etc
- [google/wire](https://github.com/google/wire) - Compile-time Dependency Injection for Go
- [go-playground/validator](https://github.com/go-playground/validator) - About 💯Go Struct and Field validation, including Cross Field, Cross Struct, Map, Slice and Array diving
  - [document](https://pkg.go.dev/github.com/go-playground/validator/v10)
- [swagger](https://swagger.io/) - Automatically generate RESTful API documentation with Swagger 2.0
  - [swaggo/swag](https://github.com/swaggo/swag#declarative-comments-format)
  - [swaggo/echo-swagger](https://github.com/swaggo/echo-swagger)
- [sirupsen/logrus](https://github.com/sirupsen/logrus) - Structured, pluggable logging for Go.

# License
[`MIT License`](./LICENSE)