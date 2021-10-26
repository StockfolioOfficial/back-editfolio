package main

// @securityDefinitions.apikey Auth-Jwt-Bearer
// @in header
// @name Authorization

// @title 에딧폴리오 public API
// @version 1.0
// @description 설명문 작성해야함
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://we.stockfolio.ai/
// @contact.email jlee@stockfolio.ai

// @license.name MIT License
// @license.url https://en.wikipedia.org/wiki/MIT_License

// @host localhost:8000
// @BasePath /
func main() {
	getApp().Start()
}