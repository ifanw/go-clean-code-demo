package mongodb

import (
	"clean_code_demo/domain"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test_asset_index(t *testing.T) {
	assertions := assert.New(t)
	asset := domain.Asset{
		ID:          domain.AssetID(uuid.NewString()),
		Description: "this a very long description",
		Label:       "there are so many labels here",
		FileName:    "demo.jpg",
	}

	// mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[defaultauthdb][?options]]
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://developer:developer@localhost:27017/demo"))
	if err != nil {
		assertions.Fail(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		assertions.Fail(err.Error())
	}
	defer client.Disconnect(ctx)

	DB := client.Database("demo")
	repo := NewAssetRepository(context.Background(), DB)
	repo.Save(asset)
}
