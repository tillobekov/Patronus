package controller

import (
	"Patronus/blockchain"
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

func (oc *OrderController) updateFilledOrders(filledOrders []*model.Order, order *model.Order) (*model.Order, error) {
	for _, o := range filledOrders {
		_, _ = oc.orderService.Update(o)
	}

	updatedOrder, err := oc.orderService.Update(order)
	oc.exchange.OrderBooks[order.Coin].ClearFilled()
	return updatedOrder, err
}

func (oc *OrderController) PostOrder(ctx *gin.Context) {
	var order *model.Order

	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	user := *ctx.MustGet("currentUser").(*model.UserDBResponseModel)
	order.UserID = user.ID.Hex()

	// Create a wallet
	manager := *blockchain.Managers.ByKey[order.Coin]
	wallet, _ := oc.walletService.FindUserWalletForNetwork(order.UserID, order.Coin)
	if wallet == nil {
		wallet = manager.CreateNewWallet()
		wallet.User = order.UserID
		wallet.Network = order.Coin
		wallet.Balance = "0"
		wallet, _ = oc.walletService.Save(wallet)
	}

	// Order processing
	order, _ = oc.orderService.Save(order)
	ob := oc.exchange.OrderBooks[order.Coin]

	if order.Type == util.OrderTypeLIMIT {
		limit, filledOrders := ob.PlaceLimitOrder(order.Price, order)
		updatedOrder, err := oc.updateFilledOrders(filledOrders, order)

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"order": model.FilteredOrderResponse(updatedOrder), "limitOrder": limit, "filledOrders": filledOrders, "err": err}})
		return
	} else if order.Type == util.OrderTypeMARKET {
		filledOrders, filledSize := ob.PlaceMarketOrder(order)
		updatedOrder, err := oc.updateFilledOrders(filledOrders, order)

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
	orders, err := oc.orderService.FindAll(user.ID.Hex())

	var response []model.Order
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
