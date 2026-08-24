package repository

import (
	"context"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"myutilityx.com/db"
	"myutilityx.com/models"
)

type FilesRepository interface {
	AddFile(ctx context.Context, f models.File) error
	DeleteFile(ctx context.Context, fileId string, userId primitive.ObjectID) error
	GetFile(ctx context.Context, f models.File) (*models.File, error)
	GetAll(ctx context.Context, f models.File, userId primitive.ObjectID) ([]*models.File, error)
}

type MongoFilesRepo struct {
	db *mongo.Collection
}

func NewFileRepository() *MongoFilesRepo {
	db, _, _ := db.InitNew()
	filesCollection := db.Collection("files")
	return &MongoFilesRepo{
		db: filesCollection,
	}
}

func (fr *MongoFilesRepo) AddFile(ctx context.Context, f models.File) error {

	f.Id = primitive.NewObjectID()
	f.CreatedAt = time.Now()
	f.ContentType = filepath.Ext(f.FileName)
	if f.ContentType != "" {
		f.ContentType = f.ContentType[1:] // Remove the leading dot
	}
	_, err := fr.db.InsertOne(ctx, f)

	if err != nil {
		return err
	}
	return nil
}

func (fr *MongoFilesRepo) GetAll(ctx context.Context, userId primitive.ObjectID) ([]*models.File, error) {
	cursor, err := fr.db.Find(ctx, bson.M{"userId": userId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []*models.File

	for cursor.Next(ctx) {
		var file models.File

		err = cursor.Decode(&file)
		if err != nil {
			return nil, err
		}
		result = append(result, &file)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (fr *MongoFilesRepo) DeleteFile(ctx context.Context, fileId string, userId primitive.ObjectID) error {
	_, err := fr.db.DeleteOne(ctx, bson.M{"fileId": fileId, "userId": userId})
	return err
}
