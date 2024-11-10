package usecase

import (
	"context"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/adiubaidah/rfid-syafiiyah/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, request *model.CreateUserRequest) (model.UserResponse, error)
	ListUsers(ctx context.Context, request *model.ListUserRequest) ([]model.UserComplete, error)
	GetUser(ctx context.Context, userId int32) (model.UserResponse, error)
	CountUsers(ctx context.Context, request *model.ListUserRequest) (int32, error)
	UpdateUser(ctx context.Context, request *model.UpdateUserRequest, userId int32) (model.UserResponse, error)
	DeleteUser(ctx context.Context, userId int32) (model.UserResponse, error)
}

type userService struct {
	store db.Store
}

func NewUserUseCase(store db.Store) UserUseCase {
	return &userService{store: store}
}

func (c *userService) CreateUser(ctx context.Context, request *model.CreateUserRequest) (model.UserResponse, error) {
	hashedPassword, err := util.HashPassword(request.Password)

	if err != nil {
		return model.UserResponse{}, err
	}

	createdUser, err := c.store.CreateUser(ctx, db.CreateUserParams{
		Username: request.Username,
		Role:     db.UserRole(request.Role),
		Password: hashedPassword,
	})
	if err != nil {
		return model.UserResponse{}, err
	}

	return model.UserResponse{
		ID:       createdUser.ID,
		Username: createdUser.Username.String,
		Role:     string(createdUser.Role.UserRole),
	}, nil

}

func (c *userService) ListUsers(ctx context.Context, request *model.ListUserRequest) ([]model.UserComplete, error) {

	arg := db.ListUserParams{
		Q:            pgtype.Text{String: request.Q, Valid: request.Q != ""},
		Role:         db.NullUserRole{UserRole: db.UserRole(request.Role), Valid: request.Role != ""},
		HasOwner:     pgtype.Bool{Bool: request.HasOwner == 1, Valid: request.HasOwner != 0},
		LimitNumber:  request.Limit,
		OffsetNumber: (request.Page - 1) * request.Limit,
		OrderBy:      db.NullUserOrderBy{UserOrderBy: db.UserOrderBy(request.Order), Valid: request.Order != ""},
	}
	// fmt.Println("argument", arg)
	// fmt.Println("q", request.Q)
	// fmt.Println("hasOwner", request.HasOwner)
	// fmt.Println("role", request.Role)
	// fmt.Println("limit", request.Limit)
	// fmt.Println("page", request.Page)
	// fmt.Println("order", request.Order)
	// userRoleValidation := db.UserRole
	users, err := c.store.ListUsers(ctx, arg)
	if err != nil {
		return []model.UserComplete{}, err
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

	return userComplete, nil
}

func (c *userService) GetUser(ctx context.Context, userId int32) (model.UserResponse, error) {
	user, err := c.store.GetUserByID(ctx, userId)
	if err != nil {
		return model.UserResponse{}, err
	}

	return model.UserResponse{
		ID:       user.ID,
		Username: user.Username.String,
		Role:     string(user.Role.UserRole),
	}, nil
}

func (c *userService) CountUsers(ctx context.Context, request *model.ListUserRequest) (int32, error) {
	count, err := c.store.CountUsers(ctx, db.CountUsersParams{
		Q:        pgtype.Text{String: request.Q, Valid: request.Q != ""},
		HasOwner: pgtype.Bool{Bool: request.HasOwner == 1, Valid: request.HasOwner != -1},
		Role:     db.NullUserRole{UserRole: db.UserRole(request.Role), Valid: request.Role != ""},
	})
	if err != nil {
		return 0, err
	}

	return int32(count), nil
}

func (c *userService) UpdateUser(ctx context.Context, request *model.UpdateUserRequest, userId int32) (model.UserResponse, error) {

	var newPassword string
	if request.Password != "" {
		hashedPassword, err := util.HashPassword(request.Password)
		if err != nil {
			return model.UserResponse{}, err
		}
		newPassword = hashedPassword
	}

	updatedUser, err := c.store.UpdateUser(ctx, db.UpdateUserParams{
		ID:       userId,
		Username: pgtype.Text{String: request.Username, Valid: request.Username != ""},
		Role:     db.NullUserRole{UserRole: db.UserRole(request.Role), Valid: request.Role != ""},
		Password: pgtype.Text{String: newPassword, Valid: newPassword != ""},
	})
	if err != nil {
		return model.UserResponse{}, err
	}

	return model.UserResponse{
		ID:       updatedUser.ID,
		Username: updatedUser.Username.String,
		Role:     string(updatedUser.Role.UserRole),
	}, nil
}

func (c *userService) DeleteUser(ctx context.Context, userId int32) (model.UserResponse, error) {
	userDeleted, err := c.store.DeleteUser(ctx, userId)
	if err != nil {
		return model.UserResponse{}, err
	}

	return model.UserResponse{
		ID:       userDeleted.ID,
		Username: userDeleted.Username.String,
		Role:     string(userDeleted.Role.UserRole),
	}, nil
}
