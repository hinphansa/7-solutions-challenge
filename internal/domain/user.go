package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// User represents the user entity in our domain
type User struct {
	ID        bson.ObjectID `json:"id" bson:"_id,omitempty" jsonschema:"title=ID,description=User ID"`
	Name      string        `json:"name" bson:"name" jsonschema:"title=Name,description=User Name,minLength=1"`
	Email     string        `json:"email" bson:"email" jsonschema:"title=Email,description=User Email,format=email"`
	Password  string        `json:"-" bson:"password" jsonschema:"title=Password,description=Bcrypt hash,minLength=8"` // "-" means this field won't be included in JSON responses
	CreatedAt time.Time     `json:"created_at" bson:"created_at" jsonschema:"title=CreatedAt,description=User Created At"`
}
