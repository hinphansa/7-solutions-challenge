package main

import (
	"regexp"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// userSchema converts the Go-struct schema into Mongo-flavoured JSON-Schema.
func userSchema() bson.M {
	emailRe := regexp.MustCompile(`^[\w\-.]+@([\w\-]+\.)+[\w\-]{2,4}$`)
	schema := bson.M{
		"bsonType": "object",
		"required": []string{"name", "email", "password", "created_at"},
		"properties": bson.M{
			"name": bson.M{
				"bsonType":    "string",
				"minLength":   1,
				"description": "non-empty name",
			},
			"email": bson.M{
				"bsonType":    "string",
				"pattern":     emailRe.String(),
				"description": "valid e-mail",
			},
			"password": bson.M{
				"bsonType":  "string",
				"minLength": 60,
				"maxLength": 60, // bcrypt hash length
			},
			"created_at": bson.M{
				"bsonType": "date",
			},
		},
	}
	return schema
}
