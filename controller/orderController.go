package controller

import (
	"Patronus/model"
	"Patronus/service"
	"Patronus/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderController struct {
	orderService  service.OrderService
	walletService service.WalletService
	exchange      model.Exchange
}

func NewOrderController(orderService service.OrderService, walletService service.WalletService, exchange model.Exchange) OrderController {
	return OrderController{orderService, walletService, exchange}
}

func (oc *OrderController) updateFilledOrders(filledOrders []*model.Order, orderFromOB *model.Order, coin model.CoinSymbol) (*model.OrderDBResponseModel, error) {
	var order model.OrderRequestModel
	for _, o := range filledOrders {
		tempOrder := model.OrderRequestModel{
			ID:     o.ID,
			Tofill: o.Size}
		_, _ = oc.orderService.Update(&tempOrder)
	}

	order.Tofill = orderFromOB.Size
	order.ID = orderFromOB.ID

	updatedOrder, err := oc.orderService.Update(&order)
	oc.exchange.OrderBooks[coin].ClearFilled()
	return updatedOrder, err
}

func (oc *OrderController) PostOrder(ctx *gin.Context) {
	var order *model.OrderRequestModel

	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	order.User = *ctx.MustGet("currentUser").(*model.UserDBResponseModel)
	savedOrder, _ := oc.orderService.Save(order)

	var wallet *model.Wallet
	wallet, _ = oc.walletService.FindUserWalletForNetwork(order.User.ID.Hex(), string(order.Coin))
	if wallet == nil {

		wallet, _ = oc.walletService.Save(&model.Wallet{
			User:    order.User.ID.Hex(),
			Address: string(order.Coin),
		})
	}

	ob := oc.exchange.OrderBooks[order.Coin]
	orderFromOB := model.NewOrder(savedOrder)

	if order.Type == util.OrderTypeLIMIT {

		limit, filledOrders := ob.PlaceLimitOrder(order.Price, orderFromOB)
		updatedOrder, err := oc.updateFilledOrders(filledOrders, orderFromOB, order.Coin)

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"order": model.FilteredOrderResponse(updatedOrder), "limitOrder": limit, "filledOrders": filledOrders, "err": err}})
		return
	} else if order.Type == util.OrderTypeMARKET {

		filledOrders, filledSize := ob.PlaceMarketOrder(orderFromOB)
		updatedOrder, err := oc.updateFilledOrders(filledOrders, orderFromOB, order.Coin)

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"order": model.FilteredOrderResponse(updatedOrder), "filledOrders": filledOrders, "filledSize": filledSize, "err": err}})
		return
	}
}

func (oc *OrderController) GetOrderBook(ctx *gin.Context) {

	ob := oc.exchange.OrderBooks[model.CoinETH]
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"asks": ob.Asks(), "bids": ob.Bids(), "marketAsks": ob.MarketAsks, "marketBids": ob.MarketBids}})
}

func (oc *OrderController) GetAllOrders(ctx *gin.Context) {
	user := *ctx.MustGet("currentUser").(*model.UserDBResponseModel)
	orders, err := oc.orderService.FindAll(user)

	var response []model.OrderResponseModel
	for _, order := range orders {
		response = append(response, model.FilteredOrderResponse(&order))
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"orders": response, "error": err}})

}

func (oc *OrderController) CancelOrder(ctx *gin.Context) {
	id := ctx.Param("id")

	order, _ := oc.orderService.FindOneById(id)
	ob := oc.exchange.OrderBooks[order.Coin]
	ob.CancelOrder(id, order.Bid)

	updatedOrder, error := oc.orderService.Cancel(id)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"order": model.FilteredOrderResponse(updatedOrder), "error": error}})

}
