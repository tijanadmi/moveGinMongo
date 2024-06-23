package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	db "github.com/tijanadmi/moveginmongo/repository"
	"github.com/tijanadmi/moveginmongo/token"
	"github.com/tijanadmi/moveginmongo/util"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      *db.MongoClient
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer( config util.Config, store *db.MongoClient) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,

	}
	/*if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}*/

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	/*if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}*/

	
	//router.POST("/users/login", server.loginUser)
	router.POST("/users/login", server.getUserByUsername)
	router.POST("/tokens/renew_access", server.renewAccessToken)
	router.GET("/halls", server.listHalls)
	router.PUT("/halls/:id",server.UpdateHall)
	router.POST("/halls", server.InsertHall)
	router.DELETE("/halls/:id", server.DeleteHall)
	
	// router.GET("/mrc/:id", server.getMrcById)
	// router.GET("/mrc", server.listMrcs)
	// router.GET("/tipprek/:id", server.getSTipPrekById)
	// router.GET("/tipprek", server.listTipPrek)
	// router.GET("/vrprek/:id", server.getSVrPrekById)
	// router.GET("/vrprek", server.listVrPrek)
	// router.GET("/uzrokprek/:id", server.getSUzrokPrekById)
	// router.GET("/uzrokprek", server.listUzrokPrek)
	// router.GET("/poduzrokprek/:id", server.getSPoduzrokPrekById)
	// router.GET("/poduzrokprek", server.listPoduzrokPrek)
	// router.GET("/mernamesta/:id", server.getSMernaMestaById)
	// router.GET("/mernamesta", server.listMernaMesta)

	// router.GET("/interruptionofdelivery/:id", server.getDDNInterruptionOfDeliveryById)
	// router.GET("/interruptionofproduction", server.listDDNInterruptionOfDeliveryP)
	// router.GET("/interruptionofusers", server.listDDNInterruptionOfDeliveryK)
	

	//authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	
	//authRoutes.GET("/accounts/:id", server.getAccount)
	

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}