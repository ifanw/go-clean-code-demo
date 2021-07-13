package application

import (
	"clean_code_demo/domain"
	"fmt"
	"github.com/google/uuid"
	assert2 "github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestAssetAppService_UploadFile(t *testing.T) {
	assert := assert2.New(t)
	type args struct {
		file  string
		label string
	}
	tests := []struct {
		name string
		args args
	}{
		/**
		Given a file named "go-go-gopher.jpg"
		When I upload the file
		Then the file name should be hashed
		*/
		{
			name: "hash real file name",
			args: args{
				file:  "go-go-gopher.jpg",
				label: "hooray",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appSrv := NewAssetAppService(newInmemoryRepo(), newInMemStorageClient())
			r, _ := appSrv.UploadFile(UploadFileCommand{
				FileName:     tt.args.file,
				Label:        tt.args.label,
				FileSize:     0,
				UploadedTime: time.Now(),
			})
			//goland:noinspection SpellCheckingInspection
			assert.Equal(r.HashedFileName, "4TB8htgio73XpfqVI5h--ts6zzMLe_zqMzHRvSDxXKE.jpg")
		})
	}
}

func TestAssetAppService_UploadFile_only_png_jpg_jpeg_mp4_mp3_pdf_Accepted(t *testing.T) {
	var memSClient *inMemStorageClient
	var memRepo *inMemoryRepo
	var assert *assert2.Assertions
	type args struct {
		file     string
		label    string
		filesize int64
		userid   string
	}
	tests := []struct {
		name     string
		args     args
		expected func(asset *UploadFileResult, err error)
	}{
		/**
		```
		Given a file named "go-go-gopher.jpg"
		When I upload the file
		The upload should be success
		```
		*/
		{
			name: "only jpg/jpeg/png/mp4/mp3/pdf accepted",
			args: args{
				file:   "go-go-gopher.jpg",
				label:  "hooray",
				userid: "XX12345",
			},
			expected: func(asset *UploadFileResult, err error) {
				assert.Nil(err)
			},
		},
		/**
		```
		Given a file named "go-go-gopher.JPG"
		When user upload the file
		The upload should be success
		```
		*/
		{
			name: "only jpg/jpeg/png/mp4/mp3/pdf accepted",
			args: args{
				file:   "go-go-gopher.JPG",
				label:  "hooray",
				userid: "XX12345",
			},
			expected: func(asset *UploadFileResult, err error) {
				assert.Nil(err)
			},
		},

		/**
		```
		Given a file named "go-go-gopher.jpeg"
		When user upload the file
		The upload should be success
		```
		*/
		{
			name: "only jpg/jpeg/png/mp4/mp3/pdf accepted",
			args: args{
				file:   "go-go-gopher.jpeg",
				label:  "hooray",
				userid: "XX12345",
			},
			expected: func(asset *UploadFileResult, err error) {
				assert.Nil(err)
			},
		},

		/**
		```
		Given a file named "go-go-gopher.mp3"
		When user upload the file
		The upload should be success
		```
		*/
		{
			name: "should mp3 format support",
			args: args{
				file:   "go-go-gopher.mp3",
				label:  "hooray",
				userid: "XX12345",
			},
			expected: func(asset *UploadFileResult, err error) {
				assert.Nil(err)
			},
		},

		/**
		```
		Given a file named "go-go-gopher.pdf"
		When user upload the file
		The upload should be success
		```
		*/
		{
			name: "should pdf format support",
			args: args{
				file:   "go-go-gopher.pdf",
				label:  "hooray",
				userid: "XX12345",
			},
			expected: func(asset *UploadFileResult, err error) {
				assert.Nil(err)
			},
		},

		//```
		//Given a file named "go-go-gopher.gif"
		//When user upload the file
		//Then upload should be fail
		// And user should get the error Message "only accept jpg/png/mp4"
		//```
		{
			name: "go-go-gopher.gif not support",
			args: args{
				file:   "go-go-gopher.gif",
				label:  "hooray",
				userid: "XX12345",
			},
			expected: func(_ *UploadFileResult, err error) {
				message := err.Error()
				fmt.Println(message)
				assert.Contains(message, "not supported")
			},
		},
		// Given a file named "go-go-gopher.jpg"
		//And the file size is 11 mb
		//When user upload the file
		//Then the upload should be failure
		{
			name: "should fail if image file size is 11mb",
			args: args{
				file:     "go-go-gopher.jpg",
				label:    "hooray",
				userid:   "XX12345",
				filesize: 10*1024*1024 + 1,
			},
			expected: func(_ *UploadFileResult, err error) {
				message := err.Error()
				assert.Contains(message, ".jpg file size limitation is 10mb, uploaded file is too large")
			},
		},
		// Given a file named "go-go-gopher.jpg"
		// And the file size is 9 mb
		// When user upload the file
		// Then the upload should be success
		{
			name: "should success if image file size smaller 10mb",
			args: args{
				file:     "go-go-gopher.jpg",
				label:    "hooray",
				userid:   "XX12345",
				filesize: 10 * 1024 * 1024,
			},
			expected: func(asset *UploadFileResult, err error) {
				assert.Nil(err)
				assert.NotNil(asset)
			},
		},
		{
			name: "should success if video file size smaller 100mb",
			args: args{
				file:     "go-go-gopher.mp4",
				label:    "hooray",
				userid:   "XX12345",
				filesize: 100 * 1024 * 1024,
			},
			expected: func(asset *UploadFileResult, err error) {
				assert.Nil(err)
				assert.NotNil(asset)
			},
		},
		{
			name: "should fail if video file size is 100mb + 1b",
			args: args{
				file:     "go-go-gopher.mp4",
				label:    "hooray",
				userid:   "XX12345",
				filesize: 100*1024*1024 + 1,
			},
			expected: func(asset *UploadFileResult, err error) {
				message := err.Error()
				assert.Contains(message, ".mp4 file size limitation is 100mb, uploaded file is too large")
			},
		},
		{
			// Given a file named "go-go-gopher.jpg"
			// and a Asset file "go-go-gopher.mp4"
			// then this file should be save in S3

			name: "should save to s3 if file is valid",
			args: args{
				file:     "go-go-gopher.mp4",
				label:    "hooray",
				userid:   "XX12345",
				filesize: 100 * 1024 * 1024,
			},
			expected: func(_ *UploadFileResult, err error) {
				assert.EqualValues(1, len(memSClient.asset))
			},
		},
		{
			name: "the status of asset should be uploaded when asset uploaded",
			args: args{
				file:     "go-go-gopher.mp4",
				label:    "hooray",
				userid:   "XX12345",
				filesize: 100 * 1024 * 1024,
			},
			expected: func(asset *UploadFileResult, err error) {
				assert.NotEmpty(asset.ID)
				assert.Equal("TRANSFERRED", asset.AssetStatus)
			},
		},
		{
			name: "should lower case extension if file extension is capital",
			args: args{
				file:     "go-go-gopher.MP4",
				label:    "hooray",
				userid:   "XX12345",
				filesize: 100 * 1024 * 1024,
			},
			expected: func(asset *UploadFileResult, err error) {
				assert.Equal("go-go-gopher.mp4", memRepo.assets["go-go-gopher.mp4"].FileName)
			},
		},
		{
			name: "should lower case extension if file extension is capital",
			args: args{
				file:     "go-go-gopher.MP3",
				label:    "hooray",
				userid:   "XX12345",
				filesize: 100 * 1024 * 1024,
			},
			expected: func(asset *UploadFileResult, err error) {
				assert.Equal("go-go-gopher.mp3", memRepo.assets["go-go-gopher.mp3"].FileName)
			},
		},
		{
			name: "should lower case extension if file extension is capital",
			args: args{
				file:     "go-go-gopher.PDF",
				label:    "hooray",
				userid:   "XX12345",
				filesize: 100 * 1024 * 1024,
			},
			expected: func(asset *UploadFileResult, err error) {
				assert.Equal("go-go-gopher.pdf", memRepo.assets["go-go-gopher.pdf"].FileName)
			},
		},
	}
	for _, tt := range tests {
		memRepo = newInmemoryRepo()
		memSClient = newInMemStorageClient()
		t.Run(tt.name, func(t *testing.T) {
			assert = assert2.New(t)
			appSrv := NewAssetAppService(memRepo, memSClient)
			command := UploadFileCommand{
				FileName: tt.args.file,
				Label:    tt.args.label,
				FileSize: tt.args.filesize,
			}
			result, err := appSrv.UploadFile(command)
			tt.expected(result, err)
		})
	}
}

func TestAssetAppService_Upload_Error_handling(t *testing.T) {
	var assert *assert2.Assertions
	//type args struct {
	//	file      string
	//	label     string
	//	filesize  int64
	//	userid    string
	//}
	tests := []struct {
		name    string
		execute func()
	}{
		{
			name: "should return error if s3 transfer fail",
			execute: func() {
				memSClient := newFailInMemStorageClient()
				memRepo := newInmemoryRepo()
				appSrv := NewAssetAppService(memRepo, memSClient)
				command := UploadFileCommand{
					FileName: "go-go-gopher.jpg", Label: "Label",
					FileSize: 1024 * 1024 * 1,
				}
				_, err := appSrv.UploadFile(command)
				assert.Error(err, "transfer to storage fail")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert = assert2.New(t)
			tt.execute()
		})
	}
}

func TestAssetAppService_Add_Fields(t *testing.T) {
	var assert *assert2.Assertions
	tests := []struct {
		name    string
		execute func()
	}{
		{
			name: "should add extra fields",
			execute: func() {
				memSClient := newInMemStorageClient()
				imRepo := newInmemoryRepo()
				appSrv := NewAssetAppService(imRepo, memSClient)
				now := time.Now()
				command := UploadFileCommand{
					FileName:     "go-go-gopher.jpg",
					Label:        "Label",
					FileSize:     1024 * 1024 * 1,
					Description:  "Description",
					UploadedTime: now,
				}
				_, _ = appSrv.UploadFile(command)
				result := imRepo.assets["go-go-gopher.jpg"] //cause the repos save twice
				assert.Equal(result.UploadedTime, now)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert = assert2.New(t)
			tt.execute()
		})
	}
}

func newInMemStorageClient() *inMemStorageClient {
	return &inMemStorageClient{
		mustFail: false,
		asset:    make(map[string]inMemClientStruct, 0),
	}
}

func newFailInMemStorageClient() *inMemStorageClient {
	return &inMemStorageClient{
		mustFail: true,
		asset:    make(map[string]inMemClientStruct, 0),
	}
}

type inMemClientStruct struct {
	filename string
	folder   string
}

type inMemStorageClient struct {
	mustFail bool
	asset    map[string]inMemClientStruct
}

func (client *inMemStorageClient) Save(_ context.Context, fileName string) (string, error) {
	if client.mustFail {
		return "", fmt.Errorf("upload to storage fail")
	}

	client.asset[fileName] = inMemClientStruct{filename: fileName}

	return "path/" + fileName, nil
}

func (client *inMemStorageClient) Delete(_ context.Context, fileName string) error {
	if client.mustFail {
		return fmt.Errorf("delete fail")
	}

	delete(client.asset, fileName)

	return fmt.Errorf("delete fail")
}

func (client *inMemStorageClient) List(_ context.Context) ([]string, error) {
	if client.mustFail {
		return nil, fmt.Errorf("delete fail")
	}
	var tmpObjects []string
	return tmpObjects, nil
}

type inMemoryRepo struct {
	assets map[string]domain.Asset
}

func (repo *inMemoryRepo) Save(asset domain.Asset) error {
	repo.assets[asset.FileName] = asset
	return nil
}

func (repo *inMemoryRepo) NewID() string {
	return uuid.NewString()
}

func (repo *inMemoryRepo) Count() int {
	return len(repo.assets)
}

func newInmemoryRepo() *inMemoryRepo {
	return &inMemoryRepo{
		make(map[string]domain.Asset, 0),
	}
}
