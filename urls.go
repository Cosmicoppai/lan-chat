package main

import (
	"lan-chat/admin/show_typ"
	"lan-chat/admin/shows"
	"lan-chat/admin/users"
	"lan-chat/admin/videos"
	"lan-chat/movieHandler"
	"lan-chat/suggestions"
	"lan-chat/utils"
	"net/http"
)

var Routes = []utils.Route{
	utils.NewRoute(http.MethodGet, "/static/(.*)", http.HandlerFunc(StaticPageHandler)),
	utils.NewRoute(http.MethodGet, "/send-suggestion", http.HandlerFunc(suggestions.FormHandler)),
	utils.NewRoute(http.MethodGet, "/list-movies", http.HandlerFunc(movieHandler.ListVideos)),
	utils.NewRoute(http.MethodGet, "/file/(.*)", http.HandlerFunc(movieHandler.GetFile)),
	utils.NewRoute(http.MethodGet, "/bwahahaha/(.*)", http.StripPrefix("/bwahahaha", TemplateHandler("./templates/admin"))),
}

var AppRoutes = [][]utils.Route{
	shows.Routes,
	users.Routes,
	show_typ.Routes,
	videos.Routes,

	// add TemplateHandler at last, otherwise all requests will be routed to serve html files
	{utils.NewRoute(http.MethodGet, "/(.*)", TemplateHandler("./templates"))},
}
