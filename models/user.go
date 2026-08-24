package models

import (
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"myutilityx.com/db"
	"myutilityx.com/utils"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Username   string             `bson:"username"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	İsVerified bool               `bson:"isVerified"`
	Role       utils.Role         `bson:"role"`
}

func (u *User) Save() error {
	database, ctx, err := db.Init()

	if err != nil {
		return err
	}
	userCollection := database.Database("myutilityx").Collection("users")
	u.Password, err = utils.HashPassword(u.Password)
	u.ID = primitive.NewObjectID()
	if err != nil {
		return err
	}
	result, err := userCollection.InsertOne(ctx, u)
	if err != nil {
		return err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		u.ID = oid
	} else {
		return errors.New("type assertion of InsertedID to primitive.ObjectID failed")
	}
	return err
}

func (u *User) ValidateCredintials(enteredPass string) error {

	err := u.FindByEmail()

	if err != nil {
		return err
	}

	if !u.İsVerified {
		return errors.New("please verify your email first")
	}

	isValid := utils.CheckPasswordHash(enteredPass, u.Password)

	if !isValid {
		return errors.New("invalid creditionals")
	}
	return nil
}

func (u User) Update(field bson.M) error {
	database, ctx, err := db.Init()

	if err != nil {
		return err
	}

	filter := bson.M{"_id": u.ID}

	_, err = database.Database(os.Getenv("MONGO_DB_NAME")).Collection("users").UpdateOne(ctx, filter, bson.M{"$set": field})

	return err
}

func (u *User) FindByEmail() error {
	database, ctx, err := db.Init()
	if err != nil {
		return err
	}

	filter := bson.M{"email": u.Email}

	err = database.Database(os.Getenv("MONGO_DB_NAME")).Collection("users").FindOne(ctx, filter).Decode(u)
	
	if err != nil {
		return err
	}
	return nil
}

func (u *User) FindById() error {
	database, ctx, err := db.Init()
	if err != nil {
		return err
	}
	filter := bson.M{"_id": u.ID}
	err = database.Database(os.Getenv("MONGO_DB_NAME")).Collection("users").FindOne(ctx, filter).Decode(u)
	return err
}

func (u User) VerifyAndUpdatePassword(oldPass string) error {
	database, ctx, err := db.Init()

	if err != nil {
		return err
	}

	filter := bson.M{"userid": u.ID}

	err = database.Database(os.Getenv("MONGO_DB_NAME")).Collection("users").FindOne(ctx, filter).Decode(&u)

	ok := utils.CheckPasswordHash(oldPass, u.Password)

	if !ok {
		return err
	}
	return err
}
