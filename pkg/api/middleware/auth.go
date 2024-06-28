package middleware

import (
	"context"
	"net/http"

	auth_service "github.com/Manjit2003/samespace/pkg/service/auth"
	"github.com/Manjit2003/samespace/pkg/utils"
)

type userContextKey string

const (
	UserKey = userContextKey("userId")
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		if token == "" {
			utils.SendResponse(w, http.StatusForbidden, utils.HTTPReponse{
				Error:   true,
				Message: "forbidden",
			})
		} else {

			userId, err := auth_service.GetUserFromJWT(token)

			if err != nil {
				utils.SendResponse(w, http.StatusForbidden, utils.HTTPReponse{
					Error:   true,
					Message: "forbidden (bad_token)",
				})
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, userId)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
	})
}
