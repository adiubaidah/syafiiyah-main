package usecase

import (
	"context"
	"errors"

	"github.com/adiubaidah/syafiiyah-main/internal/constant/exception"
	"github.com/adiubaidah/syafiiyah-main/internal/constant/model"
	repo "github.com/adiubaidah/syafiiyah-main/internal/repository"
	"github.com/adiubaidah/syafiiyah-main/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserUseCase struct {
	store repo.Store
}

func NewUserUseCase(store repo.Store) *UserUseCase {
	return &UserUseCase{store: store}
}

func (c *UserUseCase) Create(ctx context.Context, request *model.CreateUserRequest) (*model.User, error) {
	hashedPassword, err := util.HashPassword(request.Password)

	if err != nil {
		return nil, err
	}

	createdUser, err := c.store.CreateUser(ctx, repo.CreateUserParams{
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

func (c *UserUseCase) List(ctx context.Context, request *model.ListUserRequest) (*[]model.UserComplete, error) {

	arg := repo.ListUserParams{
		Q:            pgtype.Text{String: request.Q, Valid: request.Q != ""},
		Role:         repo.NullRoleType{RoleType: request.Role, Valid: request.Role != ""},
		HasOwner:     pgtype.Bool{Bool: request.HasOwner == 1, Valid: request.HasOwner != 0},
		LimitNumber:  request.Limit,
		OffsetNumber: (request.Page - 1) * request.Limit,
		OrderBy:      repo.NullUserOrderBy{UserOrderBy: repo.UserOrderBy(request.Order), Valid: request.Order != ""},
	}
	users, err := c.store.ListUsers(ctx, arg)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return nil, exception.NewNotFoundError("User not found")
		}
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

func (c *UserUseCase) GetByID(ctx context.Context, userId int32) (*model.UserWithPassword, error) {
	user, err := c.store.GetUserById(ctx, pgtype.Int4{Int32: userId, Valid: userId != 0})
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
func (c *UserUseCase) GetByUsername(ctx context.Context, username string) (*model.UserWithPassword, error) {
	user, err := c.store.GetUserByUsername(ctx, pgtype.Text{String: username, Valid: username != ""})
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
func (c *UserUseCase) GetByEmail(ctx context.Context, email string) (*model.UserWithPassword, error) {
	user, err := c.store.GetUserByEmail(ctx, pgtype.Text{String: email, Valid: email != ""})
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

func (c *UserUseCase) Count(ctx context.Context, request *model.ListUserRequest) (int64, error) {
	count, err := c.store.CountUsers(ctx, repo.CountUsersParams{
		Q:        pgtype.Text{String: request.Q, Valid: request.Q != ""},
		HasOwner: pgtype.Bool{Bool: request.HasOwner == 1, Valid: request.HasOwner != -1},
		Role:     repo.NullRoleType{RoleType: request.Role, Valid: request.Role != ""},
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest, userId int32) (*model.User, error) {

	var newPassword string
	if request.Password != "" {
		hashedPassword, err := util.HashPassword(request.Password)
		if err != nil {
			return nil, err
		}
		newPassword = hashedPassword
	}

	updatedUser, err := c.store.UpdateUser(ctx, repo.UpdateUserParams{
		ID:       userId,
		Username: pgtype.Text{String: request.Username, Valid: request.Username != ""},
		Role:     repo.NullRoleType{RoleType: request.Role, Valid: true},
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

func (c *UserUseCase) Delete(ctx context.Context, userId int32) (*model.User, error) {
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
