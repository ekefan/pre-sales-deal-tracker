package server

import (
	"fmt"
	"net/http"
	"time"

	db "github.com/ekefan/deal-tracker/internal/db/sqlc"
	"github.com/ekefan/deal-tracker/internal/utils"
	"github.com/gin-gonic/gin"
)

// CreateUsrReq holds fields that must be provided by client to create user
type CreateUsrReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Role     string `json:"role" binding:"required,valid-role"`
	FullName string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// CreateUsrRep holds fields that must be provided to client after creating a user
type CreateUsrResp struct {
	Role      string `json:"Role"`
	Username  string `json:"username"`
	CreatedAt int64  `json:"created_at"`
}

// adminCreateUserHandler http handler for the api end point for creating a new user
// must receive a CreateUserReq with username, role fullname, email and password
// on update to the handler, default password would be used so there wouldn't be a need
// to provide password in request
// FIXME: in the swagger.yml this is placed in the wrong section.
// DONE: You said I can ignore this comment
// FIXME: I can see from the endpoint address that you're not adhering to the REST-API standard.
// DONE: I have tried to adhere to the REST-API standard, this is really new to me
// If I miss anyone during this review can we talk about it, the more I write API's and read those guideline the better I will get at it
// QUESTION: How do you get information you need from long texts? Does it take long, what method do you use?
// You may want to try to follow this resource:
// - https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/
func (s *Server) adminCreateUserHandler(ctx *gin.Context) {
	var req CreateUsrReq
	// validate and bind request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Hash Password to prevent saving user password in database
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	args := db.CreateNewUserParams{
		Username: req.Username,
		Role:     req.Role,
		FullName: req.FullName,
		Email:    req.Email,
		Password: passwordHash,
	}

	user, err := s.Store.CreateNewUser(ctx, args)
	if err != nil {
		if pqErrHandler(ctx, "user", err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := CreateUsrResp{
		Role:      user.Role,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Unix(),
	}
	// FIXME: this is not compliant to REST standard. For example, return http.StatusCreated with the id of the resource created. Some framework (like Gin maybe) should also set the Location Response Header with the URI of the newly added resource.
	// DONE: I am now returning a StatusCreated code

	// We can talk about this, and I will adhere to a better design next time...
	// resourceLocation := fmt.Sprintf("/users/%s", user.Username)
	// ctx.Header("Location", resourceLocation)
	ctx.JSON(http.StatusCreated, resp)
}

// AdminUpdateUsrReq holds the field - ID to unmarshall json requests
// ID here is the id of the user to be updated
// They are all required, however if no new values are passed... the current
// the current user fields will be used
type AdminUpdateUsrReq struct {
	ID       int64  `json:"user_id" binding:"numeric,gt=0"`
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,alphanum"`
}

// FIXME: do not leave commented out code around.
// DONE: removed commented code

// adminUpdateUserHandler http handler for the api end point for updating a user
// FIXME: the PUT should entirely replace the resource. => PUT /users/1 + the request payload.
// DONE: Using PATCH as the endpoint handler doesn't update the entire user resource
func (s *Server) adminUpdateUserHandler(ctx *gin.Context) {
	var (
		req    AdminUpdateUsrReq
		newUsr db.User
	)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// get access token
	if !authAccess(ctx, utils.AdminRole) {
		return
	}
	usr, err := s.Store.GetUserForUpdate(ctx, req.ID)
	if err != nil {
		if sqlNoRowsHandler(ctx, err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// Set update time to time now....
	args := db.AdminUpdateUserParams{
		ID:        usr.ID,
		FullName:  req.Fullname,
		Email:     req.Email,
		Username:  req.Username,
		UpdatedAt: time.Now(),
	}

	if usr.Role != utils.SalesRole {
		newUsr, err = s.Store.AdminUpdateUser(ctx, args)
		if err != nil {
			if pqErrHandler(ctx, "user", err) {
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else {
		err = s.Store.UpdateUserTxn(ctx, db.UpdateUsrTxnArgs{
			ID:       usr.ID,
			Fullname: req.Fullname,
			Email:    req.Email,
			Username: req.Username,
		})
		if err != nil {
			if pqErrHandler(ctx, "user", err) {
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		newUsr, err = s.Store.GetUser(ctx, args.Username)
		if err != nil {
			if pqErrHandler(ctx, "user", err) {
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	resp := db.User{
		ID:              newUsr.ID,
		Username:        newUsr.Username,
		Role:            newUsr.Role,
		FullName:        newUsr.FullName,
		Email:           newUsr.Email,
		UpdatedAt:       newUsr.UpdatedAt,
		PasswordChanged: newUsr.PasswordChanged,
		CreatedAt:       newUsr.CreatedAt,
	}
	ctx.JSON(http.StatusOK, resp)
}

// AdminDeleteUserReq holds field user id that is to be deleted
type AdminDeleteUserReq struct {
	ID int64 `uri:"id" binding:"required"`
	// AdminRole string `uri:"admin_role" binding:"required"`
}

// adminDeleteUserhandler http handler for the api end point for Deleting a user
// FIXME: DELETE /users/1.
// DONE: using DELETE admin/users/1
// FIXME: return 204 NoContent.
// DONE: returning the correct response for a deleted resource 204
func (s *Server) adminDeleteUserHandler(ctx *gin.Context) {
	var req AdminDeleteUserReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// authenticated access
	if !authAccess(ctx, utils.AdminRole) {
		return
	}
	exists, err := s.Store.AdminUserExists(ctx, req.ID)
	if err != nil || !exists {
		ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user doesn't exist")))
		return
	}

	err = s.Store.AdminDeleteUser(ctx, req.ID)
	if err != nil {
		if sqlNoRowsHandler(ctx, err) {
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{
		"message": "successful",
	})
}


type ListUsersReq struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// FIXME: this has nothing to do with the Deals resource. Should be put in the adminUser.go
// Also you should clarify the "admin" prefix meaning. Are those endpoints restricted? Or what's else?
// listUsershandler http handler for the api end point for getting list of users currently
// DONE: moved the listUsersHandler from adminUser.go file
func (s *Server) listUsersHandler(ctx *gin.Context) {
	var req ListUsersReq
	// FIXME: if I set page_id = 0 (or outside of the allowed boundaries), please adjust it to be a default value. page_id might also start from '0'. In case you go away from conventions, you need to be declarative and put it in the documentation.
	// "(Required) The page number or offset from which to start retrieving Users. Determines where the current page of results starts in the overall list." => doesn't state the range of allowed values
	// DONE: I have put that in the documentation now...
	// For the pagination, this is a huge design issue for me and I got it wrong, lets talk about it during a session

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !multipleAuthAccess(ctx, []string{utils.AdminRole, utils.ManagerRole}) {
		return
	}
	// FIXME: if you're accepting pagination info, you should return a paginated result, not only the collection of resources.
	// DONE: still related to my poor design...
	// I wanted to use pagination, but the organisation is small, currently there are only 10 full time employees management inclusive
	args := db.AdminViewUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	users, err := s.Store.AdminViewUsers(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("i don't know")))
		return
	}
	ctx.JSON(http.StatusOK, users)
}