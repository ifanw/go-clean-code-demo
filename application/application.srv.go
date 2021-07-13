package application

import (
	"clean_code_demo/domain"
	"fmt"
	"golang.org/x/net/context"
	"time"
)

type AppService struct {
	repo          domain.AssetRepository
	storageClient domain.AssetStorageClient
}

func (assetApp *AppService) UploadFile(command UploadFileCommand) (*UploadFileResult, error) {
	ts := time.Now()
	asset, err := domain.NewAsset(domain.NewAssetParam{
		ID:           assetApp.repo.NewID(),
		FileName:     command.FileName,
		FileSize:     command.FileSize,
		Status:       domain.AssetStatusUploaded,
		UploadedTime: command.UploadedTime,
		Description:  command.Description,
		Label:        command.Label,
	})
	if err != nil {
		return nil, err
	}

	err = assetApp.repo.Save(*asset)
	if err != nil {
		return nil, err
	}

	// upload asset file
	urlAsset, err := assetApp.storageClient.Save(context.Background(), asset.HashedFileName())
	if err != nil {
		return nil, err
	}
	fmt.Printf("File uploaded success, url: %s\n", urlAsset)

	sp := time.Now().Sub(ts)
	asset.Transferred(asset.UploadedTime.Add(sp))
	err = assetApp.repo.Save(*asset)
	if err != nil {
		return nil, err
	}

	result := UploadFileResult{
		HashedFileName: asset.HashedFileName(),
		AssertURL:      urlAsset,
		ID:             string(asset.ID),
		AssetStatus:    string(asset.Status),
	}

	return &result, nil
}

func NewAssetAppService(repo domain.AssetRepository,
	storageClient domain.AssetStorageClient) *AppService {
	return &AppService{
		repo:          repo,
		storageClient: storageClient,
	}
}