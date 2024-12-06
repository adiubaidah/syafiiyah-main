package usecase

import (
	"context"
	"errors"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, request *model.CreateUserRequest) (*model.User, error)
	ListUsers(ctx context.Context, request *model.ListUserRequest) (*[]model.UserComplete, error)
	GetUser(ctx context.Context, userId int32, username string) (*model.UserWithPassword, error)
	CountUsers(ctx context.Context, request *model.ListUserRequest) (int64, error)
	UpdateUser(ctx context.Context, request *model.UpdateUserRequest, userId int32) (*model.User, error)
	DeleteUser(ctx context.Context, userId int32) (*model.User, error)
}

type userService struct {
	store db.Store
}

func NewUserUseCase(store db.Store) UserUseCase {
	return &userService{store: store}
}

func (c *userService) CreateUser(ctx context.Context, request *model.CreateUserRequest) (*model.User, error) {
	hashedPassword, err := util.HashPassword(request.Password)

	if err != nil {
		return nil, err
	}

	createdUser, err := c.store.CreateUser(ctx, db.CreateUserParams{
		Username: request.Username,
		Role:     request.Role,
		Password: hashedPassword,
	})
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:       createdUser.ID,
		Username: createdUser.Username.String,
		Role:     createdUser.Role.RoleType,
	}, nil

}

func (c *userService) ListUsers(ctx context.Context, request *model.ListUserRequest) (*[]model.UserComplete, error) {

	arg := db.ListUserParams{
		Q:            pgtype.Text{String: request.Q, Valid: request.Q != ""},
		Role:         db.NullRoleType{RoleType: request.Role, Valid: true},
		HasOwner:     pgtype.Bool{Bool: request.HasOwner == 1, Valid: request.HasOwner != 0},
		LimitNumber:  request.Limit,
		OffsetNumber: (request.Page - 1) * request.Limit,
		OrderBy:      db.NullUserOrderBy{UserOrderBy: db.UserOrderBy(request.Order), Valid: request.Order != ""},
	}
	users, err := c.store.ListUsers(ctx, arg)
	if err != nil {
		return nil, err
	}

	var userComplete []model.UserComplete
	for _, user := range users {
		userComplete = append(userComplete, model.UserComplete{
			ID:       user.ID,
			Username: user.Username,
			Role:     string(user.Role),
			UserDetails: model.UserDetails{
				ID:   user.IDOwner.Int32,
				Name: user.NameOwner.String,
			},
		})
	}

	return &userComplete, nil
}

func (c *userService) GetUser(ctx context.Context, userId int32, username string) (*model.UserWithPassword, error) {
	user, err := c.store.GetUser(ctx, db.GetUserParams{
		Username: pgtype.Text{String: username, Valid: username != ""},
		ID:       pgtype.Int4{Int32: userId, Valid: userId != 0},
	})
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("User not found")
		}

		return nil, err
	}

	return &model.UserWithPassword{
		ID:       user.ID,
		Username: user.Username.String,
		Role:     user.Role.RoleType,
		Password: user.Password.String,
	}, nil
}

func (c *userService) CountUsers(ctx context.Context, request *model.ListUserRequest) (int64, error) {
	count, err := c.store.CountUsers(ctx, db.CountUsersParams{
		Q:        pgtype.Text{String: request.Q, Valid: request.Q != ""},
		HasOwner: pgtype.Bool{Bool: request.HasOwner == 1, Valid: request.HasOwner != -1},
		Role:     db.NullRoleType{RoleType: request.Role, Valid: true},
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *userService) UpdateUser(ctx context.Context, request *model.UpdateUserRequest, userId int32) (*model.User, error) {

	var newPassword string
	if request.Password != "" {
		hashedPassword, err := util.HashPassword(request.Password)
		if err != nil {
			return nil, err
		}
		newPassword = hashedPassword
	}

	updatedUser, err := c.store.UpdateUser(ctx, db.UpdateUserParams{
		ID:       userId,
		Username: pgtype.Text{String: request.Username, Valid: request.Username != ""},
		Role:     db.NullRoleType{RoleType: request.Role, Valid: true},
		Password: pgtype.Text{String: newPassword, Valid: newPassword != ""},
	})
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("User not found")
		}
		return nil, err
	}

	return &model.User{
		ID:       updatedUser.ID,
		Username: updatedUser.Username.String,
		Role:     updatedUser.Role.RoleType,
	}, nil
}

func (c *userService) DeleteUser(ctx context.Context, userId int32) (*model.User, error) {
	userDeleted, err := c.store.DeleteUser(ctx, userId)
	if err != nil {

		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("User not found")
		}

		return nil, err
	}

	return &model.User{
		ID:       userDeleted.ID,
		Username: userDeleted.Username.String,
		Role:     userDeleted.Role.RoleType,
	}, nil
}
