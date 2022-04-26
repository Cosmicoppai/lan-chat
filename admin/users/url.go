package users

import (
	"lan-chat/middleware"
	"lan-chat/utils"
	"net/http"
)

var Routes = []utils.Route{
	utils.NewRoute(http.MethodPost, "/login", http.HandlerFunc(login)),
	utils.NewRoute(http.MethodGet, "/users", middleware.AdminMiddleware(listUsers)),
	utils.NewRoute(http.MethodPost, "/user/", middleware.AdminMiddleware(registerUser)),
	utils.NewRoute(http.MethodGet, "/user/([^/]+)", middleware.AdminMiddleware(listUser)),
	utils.NewRoute(http.MethodPut, "/user/", middleware.AdminMiddleware(updateUsername)),
	utils.NewRoute(http.MethodDelete, "/user/([^/]+)", middleware.AdminMiddleware(deleteUser)),
}
