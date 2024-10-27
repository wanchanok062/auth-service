package services

import (
	"context"
	"errors"
	"time"

	"github.com/wanchanok6698/web-auth/api/v1/models"
	"github.com/wanchanok6698/web-auth/config"
	"github.com/wanchanok6698/web-auth/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Collection *mongo.Collection
}

func NewAuthService() (*AuthService, error) {
	collection, err := config.UserCollection()
	if err != nil {
		return nil, errors.New("failed to get auth collection: " + err.Error())
	}
	return &AuthService{Collection: collection}, nil
}

func (as *AuthService) GetUserByID(ctx context.Context, id string) (*models.GetUser, error) {
	var user models.GetUser
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	err = as.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *AuthService) RegisterUser(ctx context.Context, user models.User) (string, string, error) {
	isTaken, err := us.isUserNameTaken(ctx, user.UserName)
	if err != nil {
		return "", "", err
	}
	if isTaken {
		return "", "", errors.New("username already taken")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	user.Password = string(hashPassword)

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Id = primitive.NewObjectID()

	_, err = us.Collection.InsertOne(ctx, user)
	if err != nil {
		return "", "", err
	}

	token, err := util.GenerateJWTToken(user.Id.Hex())
	if err != nil {
		return "", "", err
	}

	return user.Id.Hex(), token, nil
}
