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
	/*
		ADMIN/MANAGER ROUTERS
		router.POST("/admin/users", adminCreateUserHandler)
		router.PUT("/admin/user/update/:id", adminUpdateUserHandler)
		router.DELETE("/admin/user/delete/:id", adminDeleteUserHandler)
		router.POST("/admin/deals", adminCreateDealHandler)
		router.PUT("admin/deals/update/:id", adminUpdateDealsHandler)
		router.DELETE("/admin/deals/delete/:id", adminDeleteDealsHandler)
		router.GET("users", listUsersHandler)
	*/

	/*
		router.POST("/users/login", adminLogin)
		router.PUT("pitchrequest/update/:id", updatePitchReqHandler)
		router.GET("deals", getDealsHandler)
	*/

	/*
		SALES-REP
		router.POST("/sales/pitchReq", "salesCreatePitchReqHandler)
		router.PUT("sales/update/:username", salesUpdateuserHandler)
		router.DELETE("/sales/pitchReq/delete", salesDeletePitchReqHandler)

	*/
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