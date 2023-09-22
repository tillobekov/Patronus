package routes

import (
	"Patronus/controller"
	"Patronus/middleware"
	"Patronus/service"
	"github.com/gin-gonic/gin"
)

type OrderRouteController struct {
	orderController controller.OrderController
}

func NewRouteOrderController(orderController controller.OrderController) OrderRouteController {
	return OrderRouteController{orderController}
}

func (uc *OrderRouteController) OrderRoute(rg *gin.RouterGroup, userService service.UserService) {

	router := rg.Group("orders")
	router.Use(middleware.DeserializeUser(userService))
	router.POST("/", uc.orderController.PostOrder)

	router.GET("/all", uc.orderController.GetAllOrders)
	router.PUT("/:id", uc.orderController.CancelOrder)

	router.GET("/orderbook", uc.orderController.GetOrderBook)
}

//func CreatePost(c *gin.Context) {
//	var DB = database.ConnectDB()
//	var postCollection = database.GetUsersCollection(DB, "Posts")
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	user := new(model.User)
//	defer cancel()
//
//	if err := c.BindJSON(&user); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"message": err})
//		log.Fatal(err)
//		return
//	}
//
//	userPayload := model.User{
//		Id:    primitive.NewObjectID(),
//		Name:  user.Name,
//		Email: user.Email,
//	}
//
//	result, err := postCollection.InsertOne(ctx, userPayload)
//
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
//		return
//	}
//
//	c.JSON(http.StatusCreated, gin.H{"message": "Posted successfully", "Data": map[string]interface{}{"data": result}})
//}
