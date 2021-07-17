package main

import (
	"clean_code_demo/controller/asset"
	"clean_code_demo/core"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	sys := core.New()
	sys.Initialize()

	controller := asset.NewAsset(*sys)

	r := gin.Default()
	r.GET("/ping", controller.Ping)
	r.POST("/upload", controller.Upload)

	err := r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
