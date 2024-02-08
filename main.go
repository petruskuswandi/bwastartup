package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	webHandler "github.com/petruskuswandi/bwastartup.git/web/handler"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/petruskuswandi/bwastartup.git/handler"
	"github.com/petruskuswandi/bwastartup.git/helper"
	"github.com/petruskuswandi/bwastartup.git/payment"
	"github.com/petruskuswandi/bwastartup.git/repository"
	"github.com/petruskuswandi/bwastartup.git/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := repository.NewRepositoryUser(db)
	campaignRepository := repository.NewRepositoryCampaign(db)
	transactionRepository := repository.NewRepositoryTransaction(db)

	userService := service.NewServiceUser(userRepository)
	campaignService := service.NewServiceCampaign(campaignRepository)
	authService := service.NewServiceAuth()
	paymentService := payment.NewService()
	transactionService := service.NewServiceTransaction(transactionRepository, campaignRepository, paymentService, sync.Mutex{})

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	userWebHandler := webHandler.NewUserHandler(userService)
	campaginWebHandler := webHandler.NewCampaignHandler(campaignService, userService)

	router := gin.Default()
	router.Use(cors.Default())

	router.HTMLRender = loadTemplates("./web/templates")

	router.Static("/images", "./images")
	router.Static("/css", "./web/assets/css")
	router.Static("/js", "./web/assets/js")
	router.Static("/webfonts", "./web/assets/webfonts")

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	router.GET("/users", userWebHandler.Index)
	router.GET("/users/new", userWebHandler.New)
	router.POST("/users/", userWebHandler.Create)
	router.GET("/users/edit/:id", userWebHandler.Edit)
	router.POST("/users/update/:id", userWebHandler.Update)
	router.GET("/users/avatar/:id", userWebHandler.NewAvatar)
	router.POST("/users/avatar/:id", userWebHandler.CreateAvatar)

	router.GET("/campaigns", campaginWebHandler.Index)
	router.GET("/campaigns/new", campaginWebHandler.New)
	router.POST("/campaigns", campaginWebHandler.Create)
	router.GET("/campaigns/image/:id", campaginWebHandler.NewImage)
	router.POST("/campaigns/image/:id", campaginWebHandler.CreateImage)
	router.GET("/campaigns/edit/:id", campaginWebHandler.Edit)
	router.POST("/campaigns/update/:id", campaginWebHandler.Update)
	router.GET("/campaigns/show/:id", campaginWebHandler.Show)

	router.Run()
}

func authMiddleware(authService service.ServiceAuth, userService service.ServiceUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := claim["user_id"].(string)
		if err != nil {
			fmt.Println("Error converting userID to string:", err)
			return
		}
		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}

}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		panic(err.Error())
	}

	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}

// jalankan migrate
// migrate -database "mysql://root@tcp(localhost:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local" -path db/migrations up

// reset table/hapus table
// migrate -database "mysql://root@tcp(localhost:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local" -path db/migrations down
