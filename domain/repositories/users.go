package repositories

import (
	. "bn-survey-point/domain/datasources"
	"bn-survey-point/domain/entities"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type usersRepository struct {
	Context    context.Context
	Collection *mongo.Collection
}

type IUsersRepository interface {
	UserExist(uid string) bool
	AddCredits(uid string, credit int32) error
}

func NewUsersRepository(db *MongoDB) IUsersRepository {
	return &usersRepository{
		Context:    db.Context,
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("users"),
	}
}

func (repo usersRepository) UserExist(uid string) bool {
	var result bson.M

	filter := bson.M{"uid": uid}

	err := repo.Collection.FindOne(repo.Context, filter).Decode(&result)

	if err != nil || result == nil {
		return false
	}

	return true
}

func (repo usersRepository) AddCredits(uid string, credit int32) error {
	filter := bson.M{"uid": uid}

	_, err := repo.FindOneById(uid)

	updated := bson.M{"$inc": bson.M{"credits": credit}}

	if err != nil {
		return err
	}

	_, err = repo.Collection.UpdateOne(repo.Context, filter, updated)

	if err != nil {
		return err
	}

	return nil
}

func (repo usersRepository) FindOneById(uid string) (*entities.UserProfile, error) {
	var bsonResult bson.M

	var result entities.UserProfile

	filter := bson.M{"uid": uid}

	if err := repo.Collection.FindOne(repo.Context, filter).Decode(bsonResult); err != nil {
		return nil, err
	}

	bsonByte, err := bson.Marshal(bsonResult)

	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(bsonByte, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil

}
