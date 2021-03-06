package repository

import (
	"clean_code_demo/domain"
	"clean_code_demo/repository/mongodb"
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repo struct {
	DB    *mongo.Database
	asset domain.Repository
}

var onceAsset sync.Once

func (r *Repo) Asset() domain.Repository {
	onceAsset.Do(func() {
		r.asset = mongodb.NewAssetRepository(context.Background(), r.DB)
	})
	return r.asset
}
