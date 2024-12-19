package routers

import (
	"cat-voting-app/controllers"

	"github.com/beego/beego/v2/server/web"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	web.Router("/api/cats", &controllers.CatController{}, "get:FetchCatImages")
	web.Router("/api/breeds", &controllers.CatController{}, "get:FetchBreeds")
}
