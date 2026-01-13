package main

import "github.com/afandimsr/cashbook-backend/internal/bootstrap"

// @title           Go Gin API
// @version         1.0
// @description     Clean Architecture Go Gin API
// @termsOfService  http://swagger.io/terms/

// @contact.name   Afandi
// @contact.email  mohamadafandi71@gmail.com

// @host      localhost:8181
// @BasePath  /api/v1
func main() {
	bootstrap.Run()
}
