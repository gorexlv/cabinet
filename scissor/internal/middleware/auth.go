package middleware

import (
	"net/http"
	"strings"

	"github.com/emicklei/go-restful/v3"
	"github.com/gorexlv/cabinet/scissor/pkg/jwt"
)

func AuthMiddleware(jwtSecret string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		authHeader := req.HeaderParameter("Authorization")
		if authHeader == "" {
			resp.WriteError(http.StatusUnauthorized, restful.NewError(http.StatusUnauthorized, "missing authorization header"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			resp.WriteError(http.StatusUnauthorized, restful.NewError(http.StatusUnauthorized, "invalid authorization header"))
			return
		}

		claims, err := jwt.ValidateToken(parts[1], jwtSecret)
		if err != nil {
			resp.WriteError(http.StatusUnauthorized, restful.NewError(http.StatusUnauthorized, "invalid token"))
			return
		}

		// 将用户ID添加到请求上下文中
		req.SetAttribute("user_id", claims.UserID)
		chain.ProcessFilter(req, resp)
	}
}
