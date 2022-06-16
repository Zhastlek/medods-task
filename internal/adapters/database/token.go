package database

import (
	"context"
	"fmt"
	"log"
	"medods/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (d *db) GetOne(userGUID string, bindTokens string) (*model.UserToken, error) {
	log.Println("GET ONE MONGO DB", userGUID, bindTokens)
	filter := bson.M{
		"user_guid":   userGUID,
		"bind_tokens": bindTokens,
	}

	mongoResult := d.collection.FindOne(context.Background(), filter)
	if mongoResult.Err() != nil {
		log.Printf("Failed in find one method: %v\n", mongoResult.Err())
		return nil, mongoResult.Err()
	}
	userTokens := &model.UserToken{}
	if err := mongoResult.Decode(&userTokens); err != nil {
		log.Printf("Failed in decode method: %v\n", err)
		return nil, err
	}
	return userTokens, nil
}

func (d *db) UpdateOne(old *model.UserToken, newToken *model.UserToken) error {
	filter := bson.M{
		"refresh_token": bson.M{"$eq": old.RefreshToken},
		"bind_tokens":   bson.M{"$eq": old.BindTokens},
		"user_guid":     bson.M{"$eq": old.UserGUID},
	}
	update := bson.M{
		"$set": bson.M{
			"user_guid":     newToken.UserGUID,
			"bind_tokens":   newToken.BindTokens,
			"refresh_token": newToken.RefreshToken},
	}
	mongoResult, err := d.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Failed in update method: %v\n", err)
		return err
	}
	fmt.Println("UpdateOne() result:", mongoResult)
	return nil
}

func (d *db) CreateOne(ut *model.UserToken) error {
	mongoResult, err := d.collection.InsertOne(context.Background(), ut)
	if err != nil {
		log.Printf("Error in create one method: %v\n", err)
		return err
	}
	oid, ok := mongoResult.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println(oid.Hex())
		return err
	}
	log.Println(mongoResult.InsertedID)
	return nil
}
