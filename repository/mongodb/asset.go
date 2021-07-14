package mongodb

import (
	"clean_code_demo/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Asset struct {
	FileSize        int64     `bson:"file_size,omitempty"`
	Type            int64     `bson:"type,omitempty"`
	UploadedTime    time.Time `bson:"uploaded_time,omitempty"`
	ModifyTime      time.Time `bson:"modify_time,omitempty"`
	TransferredTime time.Time `bson:"transferred_time,omitempty"`
	ID              string    `bson:"id,omitempty"`
	Status          string    `bson:"status,omitempty"`
	FileName        string    `bson:"file_name,omitempty"`
	Label           string    `bson:"label,omitempty"`
	Description     string    `bson:"description,omitempty"`
	HashedFileName  string    `bson:"hashed_file_name,omitempty"`
}

type assetRepo struct {
	collection *mongo.Collection
	ctx        context.Context
}

func (repo *assetRepo) Save(asset domain.Asset) error {
	result := Asset{}
	err := copier.Copy(&result, &asset)
	if err != nil {
		fmt.Println(err)
		return err
	}
	result.HashedFileName = asset.HashedFileName()
	filter := bson.M{
		"id": asset.ID,
	}

	update := bson.M{"$set": &result}
	opts := options.Update().SetUpsert(true)
	_, err = repo.collection.UpdateOne(repo.ctx, filter, update, opts)
	return err
}

func (repo *assetRepo) NewID() string {
	return uuid.NewString()
}

func NewAssetRepository(ctx context.Context, db *mongo.Database) domain.Repository {
	repo := &assetRepo{
		collection: db.Collection("assets"),
		ctx:        ctx,
	}

	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "label", Value: "text"}, {Key: "description", Value: "text"}},
		},
	}

	_, _ = repo.collection.Indexes().CreateMany(
		ctx,
		models,
	)

	return repo
}
