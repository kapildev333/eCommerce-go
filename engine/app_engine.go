package engine

import (
	"eCommerce-go/db"
	config "eCommerce-go/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func RunEngine() {
	_ = godotenv.Load()
	r := gin.Default()

	ConfigRoutes(r)
	config.InitLogger()
	config.LoadConfig()

	db.InitDB()
	err := r.Run(":8080")
	if err != nil {
		fmt.Println("Error while running the server")
		return
	}
}
