package routes

import (
	"Patronus/controller"
	"Patronus/service"
	"github.com/gin-gonic/gin"
)

type CoinRouteController struct {
	coinController controller.CoinController
}

func NewCoinRouteController(coinController controller.CoinController) CoinRouteController {
	return CoinRouteController{coinController: coinController}
}

func (mc CoinRouteController) MarketRoute(rg *gin.RouterGroup, userService service.UserService) {

	router := rg.Group("coins")
	//router.Use(middleware.DeserializeUser(userService))
	router.POST("/add", mc.coinController.PostCoin)
}
