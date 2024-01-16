package main

import (
	"Patronus/blockchain"
	"Patronus/config"
	"Patronus/controller"
	"Patronus/model"
	"Patronus/routes"
	"Patronus/service"
	impl3 "Patronus/service/impl"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
)

//func main() {
//	router := gin.Default()
//	router.POST("/", routes.CreatePost)
//	router.Run("localhost:3000")
//}

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	redisclient *redis.Client

	userService  service.UserService
	authService  service.AuthService
	orderService service.OrderService
	coinService  service.CoinService
	//limitOrderService service.LimitOrderService
	walletService      service.WalletService
	transactionService service.TransactionService

	UserController  controller.UserController
	AuthController  controller.AuthController
	OrderController controller.OrderController
	CoinController  controller.CoinController

	UserRouteController  routes.UserRouteController
	AuthRouteController  routes.AuthRouteController
	OrderRouteController routes.OrderRouteController
	CoinRouteController  routes.CoinRouteController

	authCollection  *mongo.Collection
	orderCollection *mongo.Collection
	coinCollection  *mongo.Collection
	//limitOrderCollection *mongo.Collection

	transactionCollection *mongo.Collection
	walletCollection      *mongo.Collection
	managerSet            blockchain.ManagerSet

	Exchange model.Exchange
)

func init() {
	databaseName := "PatronusDB"
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()

	// Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.MongoDBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	// Connect to Redis
	redisclient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := redisclient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisclient.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client connected successfully...")

	// Collections
	authCollection = mongoclient.Database(databaseName).Collection("users")
	userService = impl3.NewUserServiceImpl(authCollection, ctx)
	authService = impl3.NewAuthService(authCollection, ctx)
	AuthController = controller.NewAuthController(authService, userService)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controller.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)

	Exchange = *model.NewExchange()
	//limitOrderCollection = mongoclient.Database(databaseName).Collection("limitOrders")
	//limitOrderService = impl3.NewLimitOrderServiceImpl(limitOrderCollection, ctx)

	walletCollection = mongoclient.Database(databaseName).Collection("wallets")
	walletService = impl3.NewWalletServiceImpl(walletCollection, ctx)
	transactionCollection = mongoclient.Database(databaseName).Collection("transactions")
	transactionService = impl3.NewTransactionServiceImpl(transactionCollection, ctx)
	orderCollection = mongoclient.Database(databaseName).Collection("orders")
	orderService = impl3.NewOrderServiceImpl(orderCollection, ctx)
	OrderController = controller.NewOrderController(orderService, walletService, transactionService, Exchange)
	OrderRouteController = routes.NewRouteOrderController(OrderController)

	coinCollection = mongoclient.Database(databaseName).Collection("coins")
	coinService = impl3.NewCoinServiceImpl(coinCollection, ctx)
	CoinController = controller.NewCoinController(coinService, Exchange)
	CoinRouteController = routes.NewCoinRouteController(CoinController)

	managerSet = blockchain.NewManagerSet()

	ethManager := blockchain.NewEthereumManager("HTTP://127.0.0.1:8545")
	managerSet.AddManager(&ethManager, "ETH")

	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoclient.Disconnect(ctx)

	value, err := redisclient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", "http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	router.GET("/order", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)
	OrderRouteController.OrderRoute(router, userService)
	CoinRouteController.MarketRoute(router, userService)

	log.Fatal(server.Run(":" + config.Port))
}
