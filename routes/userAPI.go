package routes

import (
	"Patronus/controller"
	"Patronus/middleware"
	"Patronus/service"
	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controller.UserController
}

func NewRouteUserController(userController controller.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup, userService service.UserService) {

	router := rg.Group("users")
	router.Use(middleware.DeserializeUser(userService))
	router.GET("/me", uc.userController.GetMe)
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
