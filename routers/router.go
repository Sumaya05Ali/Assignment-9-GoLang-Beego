package routers

import (
	"beeproject/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/cat", &controllers.CatController{})
	beego.Router("/breeds", &controllers.BreedsController{})
	beego.Router("/breeds/get", &controllers.BreedsController{}, "get:GetBreeds")
	beego.Router("/images/:breed_id", &controllers.BreedsController{}, "get:GetBreedImages")
	beego.Router("/save-favorite", &controllers.CatController{}, "post:SaveFavorite")
	beego.Router("/favorites", &controllers.CatController{}, "get:ShowFavorites")
}
