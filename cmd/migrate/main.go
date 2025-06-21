package main

import (
	"context"
	"errors"
	"os"

	"github.com/hinphansa/7-solutions-challenge/pkg/logger"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	log := logger.New(logrus.DebugLevel).WithFields(logrus.Fields{
		"service": "migrate",
	})

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}
	mongoDBName := os.Getenv("MONGO_DB_NAME")
	if mongoDBName == "" {
		log.Fatal("MONGO_DB_NAME is not set")
	}
	ctx := context.Background()

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Error("Failed to connect to MongoDB")
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(mongoDBName)
	if err := ensureUserCollection(ctx, log, db); err != nil {
		log.Error("Failed to ensure user collection")
		log.Fatal(err)
	}

	log.Info("migration completed")
}

func ensureUserCollection(ctx context.Context, log logger.Logger, db *mongo.Database) error {
	// Prepare schema

	const (
		collectionName   = "users"
		validationLevel  = "strict"
		validationAction = "error"
	)
	schema := userSchema()
	validator := bson.M{"$jsonSchema": schema}
	collOpts := options.CreateCollection().
		SetValidator(validator).
		SetValidationLevel(validationLevel).
		SetValidationAction(validationAction)

	// Create collection (if errors is not duplicate key error or namespace exists error, return error)
	err := db.CreateCollection(ctx, collectionName, collOpts)
	if err != nil && !mongo.IsDuplicateKeyError(err) && !isNamespaceExists(err) {
		log.Errorf("Failed to create collection: %v", err)
		return err
	}

	// Update collection schema
	cmd := bson.D{
		{Key: "collMod", Value: collectionName},
		{Key: "validator", Value: validator},
		{Key: "validationLevel", Value: "strict"},
		{Key: "validationAction", Value: "error"},
	}
	if res := db.RunCommand(ctx, cmd); res.Err() != nil {
		log.Error("Failed to update collection schema")
		return res.Err()
	}

	// Create unique index on email
	idx := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("uniq_email"),
	}
	_, err = db.Collection(collectionName).Indexes().CreateOne(ctx, idx)
	if err != nil {
		log.Error("Failed to create unique index on email")
	}
	return err
}

func isNamespaceExists(err error) bool {
	var cmdErr mongo.CommandError
	return mongo.IsDuplicateKeyError(err) ||
		(errors.As(err, &cmdErr) && cmdErr.Code == 48) // NamespaceExists
}
