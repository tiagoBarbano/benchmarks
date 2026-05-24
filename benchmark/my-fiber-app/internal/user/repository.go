package user

import (
	"context"
	"errors"

	"my-fiber-app/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id string) (*DtoUserResponse, error)
	Update(ctx context.Context, id string, user *User) (*User, error)
	Delete(ctx context.Context, id string) error
}

type mongoRepository struct {
	collection *mongoDriver.Collection
}

func NewRepository() Repository {
	return &mongoRepository{
		collection: mongo.DB.Collection("users"),
	}
}

func (r *mongoRepository) Create(ctx context.Context, user *User) (*User, error) {
	res, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (r *mongoRepository) GetAll(ctx context.Context) ([]User, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *mongoRepository) GetByID(ctx context.Context, id string) (*DtoUserResponse, error) {
	// oid, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	return nil, errors.New("invalid id")
	// }

	var user DtoUserResponse
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongoDriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *mongoRepository) Update(ctx context.Context, id string, user *User) (*User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	update := bson.M{
		"$set": bson.M{
			"name":  user.Name,
			"email": user.Email,
		},
	}

	_, err = r.collection.UpdateByID(ctx, oid, update)
	if err != nil {
		return nil, err
	}
	user.ID = oid
	return user, nil
}

func (r *mongoRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
