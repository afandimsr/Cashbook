package main

import "github.com/afandimsr/cashbook-backend/internal/bootstrap"

// @title           CashBook - Personal Finance API
// @version         1.0
// @description     API Server for CashBook Personal Finance Manager. This API handles transactions, budgeting, category management, and financial reporting.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Afandi
// @contact.email  mohamadafandi71@gmail.com

// @host      localhost:8181
// @BasePath  /api/v1
func main() {
	bootstrap.Run()
}
