package ports

import (
	"context"

	"github.com/hinphansa/7-solutions-challenge/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*bson.ObjectID, error)
	GetByID(ctx context.Context, id bson.ObjectID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	List(ctx context.Context, pagination *Pagination) ([]domain.User, error)
	Update(ctx context.Context, id bson.ObjectID, user *domain.User) error
	Delete(ctx context.Context, id bson.ObjectID) error
	Count(ctx context.Context) (int64, error)
}

type UserService interface {
	Register(ctx context.Context, user *domain.User) (*bson.ObjectID, error)
	Login(ctx context.Context, email string, password string) (string, error) // Authentication
	GetByID(ctx context.Context, id bson.ObjectID) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	List(ctx context.Context, pagination *Pagination) ([]domain.User, error)
	Update(ctx context.Context, id bson.ObjectID, user *domain.User) error
	Delete(ctx context.Context, id bson.ObjectID) error
}
