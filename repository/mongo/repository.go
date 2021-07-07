package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	Db *mongo.Database
}

func useUsersCollection(r *Repository) *mongo.Collection {
	return r.Db.Collection("users")
}

func usePostsCollection(r *Repository) *mongo.Collection {
	return r.Db.Collection("posts")
}

func useCommentsCollection(r *Repository) *mongo.Collection {
	return r.Db.Collection("comments")
}