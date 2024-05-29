package repositories

import (
	"bn-survey-point/domain/entities"
	. "bn-survey-point/domain/datasources"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AlertMessageRepositories struct {
	Context    context.Context
	Collection *mongo.Collection
}

type IAlertMessageRepositories interface {
	GetSurveyCredits() (int32, error)
	GetAll() (*entities.AlertMessage, error)
}

func NewAlertMessageRepositories(db *MongoDB) IAlertMessageRepositories {
	return &AlertMessageRepositories{
		Context:    db.Context,
		Collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("alert_message"),
	}
}

func (repo AlertMessageRepositories) GetSurveyCredits() (int32, error) {
	var result entities.AlertMessage
	filter := bson.M{}

	err := repo.Collection.FindOne(repo.Context, filter).Decode(&result)

	if err != nil {
		return 0, err
	}

	credit := result.Teacher.Point

	return credit, nil
}

func (repo AlertMessageRepositories) GetAll() (*entities.AlertMessage, error) {
	var result entities.AlertMessage
	filter := bson.M{}

	err := repo.Collection.FindOne(repo.Context, filter).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
