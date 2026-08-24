package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"myutilityx.com/db"
	"myutilityx.com/graph/model"
)

type LinkRepository interface {
	FindAll() ([]*model.Link, error)
	FindAllByUserID(userId primitive.ObjectID) ([]*model.Link, error)
}

type database struct {
	client *mongo.Client
}

func New() LinkRepository {
	client, _, _ := db.Init()

	return &database{
		client: client,
	}
}

func (d *database) FindAll() ([]*model.Link, error) {
	client, _, _ := db.Init()

	linksCollection := client.Database("myutilityx").Collection("links")

	cursor, err := linksCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var result []*model.Link

	for cursor.Next(context.TODO()) {
		var l *model.Link
		err = cursor.Decode(&l)
		if err != nil {
			return nil, err
		}
		result = append(result, l)
	}

	return result, err
}

func (d *database) FindAllByUserID(userId primitive.ObjectID) ([]*model.Link, error) {
	client, _, _ := db.Init()
	linksCollection := client.Database("myutilityx").Collection("links")

	cursor, err := linksCollection.Find(context.TODO(), bson.M{"userid": userId})
	if err != nil {
		return nil, err
	}

	var links []*model.Link
	if err = cursor.All(context.TODO(), &links); err != nil {
		return nil, err
	}
	return links, nil
}
