package services

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/wanchanok6698/web-auth/api/v1/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (us *AuthService) isUserNameTaken(ctx context.Context, userName string) (bool, error) {
	var user models.User
	err := us.Collection.FindOne(ctx, bson.M{"username": userName}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}

		logrus.Error("Error checking username:", err)
		return false, err
	}

	return true, nil
}
