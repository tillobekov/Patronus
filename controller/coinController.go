package controller

import (
	"Patronus/model"
	"Patronus/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CoinController struct {
	coinService service.CoinService
	exchange    model.Exchange
}

func NewCoinController(coinService service.CoinService, exchange model.Exchange) CoinController {
	return CoinController{coinService: coinService, exchange: exchange}
}

func (mc CoinController) PostCoin(ctx *gin.Context) {
	var coinModel *model.CoinRequestModel

	if err := ctx.ShouldBindJSON(&coinModel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	coinFromDB, err := mc.coinService.Save(coinModel)

	//fmt.Errorf("err %+v", err)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"coin": coinFromDB, "err": err}})
}
