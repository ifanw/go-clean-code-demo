package main

import (
	"clean_code_demo/controller/asset"
	"clean_code_demo/core"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	sys := core.New()
	sys.Initialize()

	asset := asset.NewAsset(*sys)

	r := gin.Default()
	r.GET("/ping", asset.Ping)
	r.POST("/upload", asset.Upload)

	err := r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		fmt.Println(err)
	}
}
