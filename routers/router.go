package routers

import (
	"cat-voting-app/controllers"

	"github.com/beego/beego/v2/server/web"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	// beego.Router("/", &controllers.MainController{}, "get:HomePage")
	web.Router("/api/cats", &controllers.CatController{}, "get:FetchCatImages")
	web.Router("/api/breeds", &controllers.CatController{}, "get:FetchBreeds")
	web.Router("/api/add-to-favourites", &controllers.CatController{}, "post:AddToFavourites")
	web.Router("/api/get-favourites", &controllers.CatController{}, "get:GetFavourites")
	web.Router("/api/remove-favourite", &controllers.CatController{}, "delete:RemoveFavourite")
	web.Router("/api/vote", &controllers.CatController{}, "post:Vote")
	web.Router("/api/votes", &controllers.CatController{}, "get:GetVotes")
}
