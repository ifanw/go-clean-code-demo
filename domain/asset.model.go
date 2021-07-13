package domain

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func (asset *Asset) HashedFileName() string {
	if asset.hashedFileName == "" {
		extFile := filepath.Ext(asset.FileName)

		f := sha256.Sum256([]byte(asset.FileName))
		hashed := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(f[:])
		asset.hashedFileName = hashed + extFile
	}
	return asset.hashedFileName
}

func (asset *Asset) Transferred(t time.Time) {
	asset.Status = AssetStatusTransferred
	asset.TransferredTime = t
}

func fileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}

func NewAsset(param NewAssetParam) (*Asset, error) {
	specError := NewSpecBuilder()

	getAssetTypeAndValidate := func(fileName string) (AssetType, string) {
		ext := strings.ToLower(filepath.Ext(fileName))
		assetType, exists := availableExt[ext]

		if !exists {
			specError.AppendBadRequest(fmt.Sprintf("%q not supported", ext))
		}

		if param.FileSize > fileSizeMap[assetType] {
			specError.AppendBadRequest(fmt.Sprintf("%s file size limitation is %dmb, uploaded file is too large", ext, fileSizeMap[assetType]/1024/1024))
		}

		lowerCaseExtFileName := fileNameWithoutExtension(param.FileName) + ext
		return assetType, lowerCaseExtFileName
	}

	assetType, fileName := getAssetTypeAndValidate(param.FileName)

	m := &Asset{
		ID:              AssetID(param.ID),
		FileName:        fileName,
		FileSize:        param.FileSize,
		Type:            assetType,
		Status:          param.Status,
		UploadedTime:    param.UploadedTime,
		TransferredTime: param.TransferredTime,
		Label:           param.Label,
		Description:     param.Description,
	}

	return m, specError.GetIfAny()
}
