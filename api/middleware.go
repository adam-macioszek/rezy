package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/adam-macioszek/rezy/token"
	"github.com/gin-gonic/gin"
)

const authorizationTypreBearer = "bearer"
const authorizationPayloadKey = "authorization_payload"

func authMiddleware(tokenMaker token.JWTMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authorizationHeader := ctx.GetHeader("authorization")
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is empty")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))

		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypreBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		}
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()

	}
}
