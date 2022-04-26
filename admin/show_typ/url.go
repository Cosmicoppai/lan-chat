package show_typ

import (
	"lan-chat/middleware"
	"lan-chat/utils"
	"net/http"
)

var Routes = []utils.Route{
	utils.NewRoute(http.MethodGet, "/types", http.HandlerFunc(listTypes)),
	utils.NewRoute(http.MethodPost, "/type/", middleware.AuthMiddleware(addType)),
	utils.NewRoute(http.MethodPut, "/type/([0-9]+)", middleware.AuthMiddleware(updateTypName)),
	utils.NewRoute(http.MethodDelete, "type/([0-9]+)", middleware.AuthMiddleware(deleteType)),
}
