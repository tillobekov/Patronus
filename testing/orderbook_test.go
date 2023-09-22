package testing

import (
	"Patronus/model"
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

func TestLimit(t *testing.T) {
	l := model.NewLimit(10_000)
	buyOrderA := model.NewOrder(true, 5)
	buyOrderB := model.NewOrder(true, 5)
	buyOrderC := model.NewOrder(true, 5)

	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	l.DeleteOrder(buyOrderB)
	fmt.Println(l)
}

func TestPlaceLimitOrder(t *testing.T) {
	ob := model.NewOrderBook()

	sellOrderA := model.NewOrder(false, 10)
	sellOrderB := model.NewOrder(false, 5)
	ob.PlaceLimitOrder(10_000, sellOrderA)
	ob.PlaceLimitOrder(9_000, sellOrderB)

	assert(t, len(ob.Asks()), 2)
}

func TestPlaceMarketOrder(t *testing.T) {
	ob := model.NewOrderBook()

	sellOrder := model.NewOrder(false, 20)
	ob.PlaceLimitOrder(10_000, sellOrder)

	buyOrder := model.NewOrder(true, 10)
	matches := ob.PlaceMarketOrder(buyOrder)

	assert(t, len(matches), 1)
	assert(t, len(ob.Asks()), 1)
	assert(t, ob.AskTotalVolume(), 10.0)
	assert(t, matches[0].Ask, sellOrder)
	assert(t, matches[0].Bid, buyOrder)
	assert(t, matches[0].SizeFilled, 10.0)
	assert(t, matches[0].Price, 10_000.0)
	assert(t, buyOrder.IsFilled(), true)

	fmt.Printf("%+v", matches)
}

func TestPlaceMarketOrderMultiFill(t *testing.T) {
	ob := model.NewOrderBook()

	buyOrderA := model.NewOrder(true, 5)
	buyOrderB := model.NewOrder(true, 8)
	buyOrderC := model.NewOrder(true, 10)
	buyOrderD := model.NewOrder(true, 1)

	ob.PlaceLimitOrder(5_000, buyOrderC)
	ob.PlaceLimitOrder(5_000, buyOrderD)
	ob.PlaceLimitOrder(9_000, buyOrderB)
	ob.PlaceLimitOrder(10_000, buyOrderA)

	assert(t, ob.BidTotalVolume(), 24.0)

	sellOrder := model.NewOrder(false, 20)
	matches := ob.PlaceMarketOrder(sellOrder)

	assert(t, ob.BidTotalVolume(), 4.0)
	assert(t, len(matches), 3)
	assert(t, len(ob.Bids()), 1)

	fmt.Printf("%+v", matches)

}

func TestCancelOrder(t *testing.T) {
	ob := model.NewOrderBook()
	buyOrder := model.NewOrder(true, 4)
	ob.PlaceLimitOrder(10_000, buyOrder)
	assert(t, ob.BidTotalVolume(), 4.0)

	ob.CancelOrder(buyOrder)
	assert(t, ob.BidTotalVolume(), 0.0)
}
