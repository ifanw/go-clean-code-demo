package main

import (
	"clean_code_demo/application"
	"clean_code_demo/core"
	"clean_code_demo/domain"
	"clean_code_demo/repository/awss3"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"time"
)

func parseUploadFiles(ctx *gin.Context, key string) (multipart.File, *multipart.FileHeader, error) {
	file, header, err := ctx.Request.FormFile(key)
	if err != nil {
		return nil, nil, err
	}

	if header == nil {
		specErr := domain.NewSpecBuilder()
		specErr.Append(domain.NewBadRequest("header is missing"))
		return nil, nil, specErr.GetIfAny()
	}

	fmt.Printf("filename: %s, size: %d, mimeType: %s \n", header.Filename, header.Size, header.Header)

	return file, header, err
}

func main() {
	sys := core.New()
	sys.Initialize()

	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/upload", func(ctx *gin.Context) {

		label := ctx.PostForm("label")
		description := ctx.PostForm("description")

		fileAsset, headerAsset, err := parseUploadFiles(ctx, "asset")
		if err != nil {
			ctx.JSON(400, gin.H{
				"err": err,
			})
		}

		// hardcode for demo
		config := awss3.AWSConfig{
			Address: "http://localhost:4566",
			Bucket:  "bucket.name", // you need to create
			Region:  "us-west-2",
			Profile: "localstack",
			ID:      "id",
			Secret:  "secret",
		}

		// make channel for files
		ch := make(chan io.Reader, 1)
		defer close(ch)

		ch <- fileAsset

		client, err := awss3.NewStorageClient(config, ch, time.Second*5)
		if err != nil {
			ctx.JSON(400, gin.H{
				"err": err,
			})
		}

		repoAsset := sys.GetRepository().Asset()
		appSrv := application.NewAssetAppService(repoAsset, client)

		cmd := application.UploadFileCommand{
			FileName:    headerAsset.Filename,
			Label:       label,
			FileSize:    headerAsset.Size,
			Description: description,
		}

		result, err := appSrv.UploadFile(cmd)
		if err != nil {
			ctx.JSON(400, gin.H{
				"err": err,
			})
		}

		fmt.Printf("result: %+v\n", result)
		ctx.JSON(200, gin.H{
			"success": result,
		})

	})

	err := r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		fmt.Println(err)
	}
}
