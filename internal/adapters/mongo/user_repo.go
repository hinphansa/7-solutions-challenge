package mongo

import (
	"context"
	"errors"

	"github.com/hinphansa/7-solutions-challenge/internal/domain"
	"github.com/hinphansa/7-solutions-challenge/internal/ports"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// compile time check to ensure mongoRepository implements ports.UserRepository
var _ ports.UserRepository = (*userRepository)(nil)

const (
	collectionName = "users"
)

type userRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *userRepository {
	return &userRepository{coll: db.Collection(collectionName)}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*bson.ObjectID, error) {
	res, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID.(bson.ObjectID)
	return &id, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var result *domain.User
	if err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *userRepository) GetByID(ctx context.Context, id bson.ObjectID) (*domain.User, error) {
	var result *domain.User
	if err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	cursor, err := r.coll.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []domain.User
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *userRepository) List(ctx context.Context, pagination *ports.Pagination) ([]domain.User, error) {
	limit := pagination.Limit
	skip := pagination.Offset * pagination.Limit
	cursor, err := r.coll.Find(ctx, bson.D{{}}, options.Find().SetLimit(limit).SetSkip(skip))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []domain.User
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *userRepository) Update(ctx context.Context, id bson.ObjectID, user *domain.User) error {
	updateFields := bson.M{}
	if user.Email != "" {
		updateFields["email"] = user.Email
	}
	if user.Name != "" {
		updateFields["name"] = user.Name
	}

	if len(updateFields) == 0 {
		return errors.New("no fields to update")
	}

	_, err := r.coll.UpdateByID(ctx, id, bson.M{"$set": updateFields})
	return err
}

func (r *userRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	return r.coll.CountDocuments(ctx, nil)
}
