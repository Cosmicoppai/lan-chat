package shows

import (
	"lan-chat/middleware"
	"lan-chat/utils"
	"net/http"
)

var Routes = []utils.Route{
	utils.NewRoute(http.MethodGet, "/shows", http.HandlerFunc(listShows)),
	utils.NewRoute(http.MethodGet, "/show/([0-9]+)", http.HandlerFunc(listShow)),
	utils.NewRoute(http.MethodPost, "/show/", middleware.AuthMiddleware(createShow)),
	utils.NewRoute(http.MethodPut, "/show/([0-9]+)", middleware.AuthMiddleware(updateShow)),
	utils.NewRoute(http.MethodDelete, "/show/([0-9]+)", middleware.AuthMiddleware(deleteShow)),
}
