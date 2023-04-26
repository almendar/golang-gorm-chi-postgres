package main

import (
	"fmt"
	"net/http"

	"github.com/almendar/golang-gorm-chi-postgres/dogowners"
	"github.com/almendar/golang-gorm-chi-postgres/shared"

	_ "github.com/almendar/golang-gorm-chi-postgres/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// swagger embed files
//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:3000
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	db, err := shared.DefaultGormHandle()
	if err != nil {
		fmt.Printf("failed to connect to database: %s\n", err)
		return
	}
	dbStorage := dogowners.NewDatabase(db)
	svc := dogowners.NewService(dbStorage)
	handler := dogowners.NewHttpHandlers(svc)

	handler.R.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/swagger/doc.json"), //The url pointing to API definition
	))
	http.ListenAndServe("localhost:3000", handler)
}
