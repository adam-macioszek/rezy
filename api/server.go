package api

import (
	"errors"
	"log"

	"github.com/adam-macioszek/rezy/config"
	db "github.com/adam-macioszek/rezy/db/sqlc"
	"github.com/adam-macioszek/rezy/token"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     config.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.JWTMaker
}

func NewServer(config config.Config, store db.Store) (*Server, error) {
	maker, err := token.NewJWTMaker(config.ApiTokenKey)
	if err != nil {
		return &Server{}, errors.New("cannot create token key")
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: maker,
	}
	router := gin.Default()

	router.GET("/reservation", server.listReservation)
	router.GET("/reservation/:id", server.getReservation)
	router.POST("/book", server.createBooking)
	authroute := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authroute.POST("/reservation", server.createReservation)
	server.router = router

	return server, nil
}

func (server *Server) Start(address string) error {
	apiToken, _ := server.tokenMaker.CreateToken("root", server.config.ApiTokenDuration)
	log.Println(apiToken)
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error: ": err.Error()}
}
