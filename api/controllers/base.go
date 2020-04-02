package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/HackathonCovid/helpCovidBack/api/middlewares"
	"github.com/HackathonCovid/helpCovidBack/api/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
)

// Server : Struct
type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

var errList = make(map[string]string)

// Initialize : connection to our database
func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		log.Fatal("This is the error connecting to postgres:", err)
	} else {
		fmt.Printf("We are connected to the %s database", Dbdriver)
	}

	//database migration
	server.DB.Debug().AutoMigrate(
		&models.User{},
		&models.Mission{},
		&models.ResetPassword{},
		&models.Comment{},
		&models.Apply{},
	)

	server.Router = gin.Default()
	server.Router.Use(middlewares.CORSMiddleware())

	server.initializeRoutes()
}

// Run : Status
func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
