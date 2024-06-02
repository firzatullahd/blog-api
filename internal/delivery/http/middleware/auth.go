package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/firzatullahd/blog-api/internal/model"
	"github.com/firzatullahd/blog-api/internal/model/response"
	"github.com/golang-jwt/jwt/v5"
)

const UserDataKey = "user-data"

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		if strings.Contains(accessToken, "Bearer") {
			accessToken = strings.Replace(accessToken, "Bearer ", "", -1)
		}

		claims := model.MyClaim{}
		token, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.JWTSecretKey), nil
		})

		if err != nil {
			response.SetHTTPResponse(w, http.StatusUnauthorized, "token invalid", nil)
			return
		}

		if !token.Valid {
			response.SetHTTPResponse(w, http.StatusUnauthorized, "token invalid", nil)
			return
		}

		timeExp, err := token.Claims.GetExpirationTime()
		if err != nil {
			response.SetHTTPResponse(w, http.StatusUnauthorized, "token invalid", nil)
			return
		}

		res := timeExp.Compare(time.Now())
		if res == -1 {
			response.SetHTTPResponse(w, http.StatusUnauthorized, "token expired", nil)
			return
		}

		ctx := context.WithValue(r.Context(), UserDataKey, claims.UserData)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
