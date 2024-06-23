package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tijanadmi/moveginmongo/models"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
}

func newUserResponse(user *models.User) userResponse {
	return userResponse{
		Username:          user.Username,
	}
}



type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

// func (server *Server) loginUser(ctx *gin.Context) {
// 	var req loginUserRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	user, err := server.store.Users.GetUserByUsername(ctx,req.Username)
// 	if err != nil {
		
// 		ctx.JSON(http.StatusNotFound, errorResponse(err))
// 		return
// 	}

// 	err = util.CheckPassword(req.Password, user.Password)
// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
// 		return
// 	}

// 	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
// 		user.Username,
// 		/*user.Role,*/
// 		"depositor",
// 		server.config.AccessTokenDuration,
// 	)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
// 		user.Username,
// 		//user.Role,
// 		"depositor",
// 		server.config.RefreshTokenDuration,
// 	)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	/*session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
// 		ID:           refreshPayload.ID,
// 		Username:     user.Username,
// 		RefreshToken: refreshToken,
// 		UserAgent:    ctx.Request.UserAgent(),
// 		ClientIp:     ctx.ClientIP(),
// 		IsBlocked:    false,
// 		ExpiresAt:    refreshPayload.ExpiredAt,
// 	})
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}*/

// 	rsp := loginUserResponse{
// 		//SessionID:             session.ID,
// 		AccessToken:           accessToken,
// 		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
// 		RefreshToken:          refreshToken,
// 		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
// 		User:                  newUserResponse(user),
// 	}
// 	ctx.JSON(http.StatusOK, rsp)
// }

func (server *Server) getUserByUsername(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.Users.GetUserByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}