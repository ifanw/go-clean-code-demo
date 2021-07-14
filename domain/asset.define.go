package domain

import (
	"context"
	"time"
)

type AssetID string
type AssetType int
type AssetStatus string
type AssetShareType string
type ShareTo string

const (
	AssetTypeImage AssetType = 1
	AssetTypeVideo AssetType = 2
)

const (
	AssetStatusUploaded    AssetStatus = "UPLOADED"
	AssetStatusTransferred AssetStatus = "TRANSFERRED"
)

var availableExt = map[string]AssetType{
	".png":  AssetTypeImage,
	".jpg":  AssetTypeImage,
	".jpeg": AssetTypeImage,
	".mp4":  AssetTypeVideo,
	".mp3":  AssetTypeVideo,
	".pdf":  AssetTypeVideo,
}

var fileSizeMap = map[AssetType]int64{
	AssetTypeImage: 10 * 1024 * 1024,  // 10 mb
	AssetTypeVideo: 100 * 1024 * 1024, // 100 mb
}

type Asset struct {
	FileSize        int64
	Type            AssetType
	UploadedTime    time.Time
	ModifyTime      time.Time
	TransferredTime time.Time
	ID              AssetID
	Status          AssetStatus
	FileName        string
	Label           string
	Description     string
	hashedFileName  string
}

type NewAssetParam struct {
	FileSize        int64
	UploadedTime    time.Time
	TransferredTime time.Time
	Status          AssetStatus
	FileName        string
	ID              string
	Label           string
	Description     string
}

type Repository interface {
	Save(asset Asset) error
	NewID() string
}

type StorageClient interface {
	// Save a new object to a bucket and returns its URL to view/download.
	Save(ctx context.Context, fileName string) (string, error)
	// Delete an existing object from a bucket.
	Delete(ctx context.Context, fileName string) error
	// List all objects in a bucket.
	List(ctx context.Context) ([]string, error)
}
