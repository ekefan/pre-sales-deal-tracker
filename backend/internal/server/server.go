package server

import (
	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server contains fields required by the server instance
type Server struct {
	Router *gin.Engine
	Store  db.Store
}

// NewServer create a server instance, having a router that connect api endpoints
func NewServer(store db.Store) *Server {
	return &Server{
		Store: store,
	}
}

// SetupRouter ini
func (s *Server) SetupRouter() {
	router := gin.Default()

	// ADMIN
	router.POST("/user", s.adminCreateUserHandler)
	router.PUT("/user/update/", s.adminUpdateUserHandler)
	router.DELETE("/user/delete/:id", s.adminDeleteUserHandler)
	router.POST("/admin/deals", s.adminCreateDealHandler)
	router.PUT("admin/deals/update", s.adminUpdateDealHandler)
	router.DELETE("/admin/deals/delete/:deal_id", s.adminDeleteDealHandler)
	router.GET("/users", s.listUsersHandler)

	// ADMINSALES
	router.POST("/users/login", s.userLogin)
	router.PUT("pitchrequest/update", s.updatePitchReqHandler)
	router.GET("/deals", s.getDealsHandler)


	//SALES-REP
	router.POST("/sales/pitchReq", s.salesCreatePitchReqHandler)
	router.PUT("/sales/update/:username", s.salesUpdateuserHandler)
	router.GET("/pitchrequest/", s.salesViewPitchRequests)
	router.DELETE("/sales/pitchReq/delete/:pitch_id", s.salesDeletePitchReqHandler)

	s.Router = router
}

// hostAddress string
func (s *Server) StartServer(hostAddress string) error {
	err := s.Router.Run(hostAddress)
	if err != nil {
		return err
	}
	return nil
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
