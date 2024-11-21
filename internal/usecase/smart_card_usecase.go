package usecase

import (
	"context"

	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/exception"
	"github.com/adiubaidah/rfid-syafiiyah/internal/constant/model"
	db "github.com/adiubaidah/rfid-syafiiyah/internal/storage/persistence"
	"github.com/jackc/pgx/v5/pgtype"
)

type SmartCardUseCase interface {
	CreateSmartCard(ctx context.Context, request *model.SmartCardRequest) (model.SmartCard, error)
	ListSmartCards(ctx context.Context, request *model.ListSmartCardRequest) ([]model.SmartCardComplete, error)
	CountSmartCards(ctx context.Context, request *model.ListSmartCardRequest) (int64, error)
	UpdateSmartCard(ctx context.Context, request *model.UpdateSmartCardRequest, id int32) (model.SmartCardComplete, error)
	DeleteSmartCard(ctx context.Context, id int32) (model.SmartCard, error)
}

type service struct {
	store db.Store
}

func NewSmartCardUseCase(store db.Store) SmartCardUseCase {
	return &service{store: store}
}

func (c *service) CreateSmartCard(ctx context.Context, request *model.SmartCardRequest) (model.SmartCard, error) {
	createdSmartCard, err := c.store.CreateSmartCard(ctx, db.CreateSmartCardParams{
		Uid:        request.Uid,
		IsActive:   true,
		SantriID:   pgtype.Int4{Valid: false},
		EmployeeID: pgtype.Int4{Valid: false},
	})

	if err != nil {
		if exception.DatabaseErrorCode(err) == exception.ErrCodeUniqueViolation {
			return model.SmartCard{}, exception.NewUniqueViolationError("", err)
		}
		return model.SmartCard{}, err
	}

	return model.SmartCard{
		ID:        createdSmartCard.ID,
		Uid:       createdSmartCard.Uid,
		CreatedAt: createdSmartCard.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		IsActive:  createdSmartCard.IsActive,
	}, nil
}

func (c *service) ListSmartCards(ctx context.Context, request *model.ListSmartCardRequest) ([]model.SmartCardComplete, error) {
	listSmartCard, err := c.store.ListSmartCards(ctx, db.ListSmartCardsParams{
		Q:            pgtype.Text{String: request.Q, Valid: request.Q != ""},
		IsActive:     pgtype.Bool{Bool: true, Valid: true},
		IsSantri:     pgtype.Bool{Bool: request.OwnerRole == "santri", Valid: request.OwnerRole != ""},
		IsEmployee:   pgtype.Bool{Bool: request.OwnerRole == "employee", Valid: request.OwnerRole != ""},
		OffsetNumber: request.Limit * (request.Page - 1),
		LimitNumber:  request.Limit,
	})

	if err != nil {
		return nil, err
	}

	var result []model.SmartCardComplete

	for _, smartCard := range listSmartCard {

		ownerRole := ""
		detailsId := int32(0)
		detailsName := ""
		if smartCard.SantriID.Valid {
			ownerRole = "santri"
			detailsId = smartCard.SantriID.Int32
			detailsName = smartCard.SantriName.String
		} else if smartCard.EmployeeID.Valid {
			ownerRole = "employee"
			detailsId = smartCard.EmployeeID.Int32
			detailsName = smartCard.EmployeeName.String
		}

		result = append(result, model.SmartCardComplete{
			SmartCard: model.SmartCard{
				ID:        smartCard.ID,
				Uid:       smartCard.Uid,
				CreatedAt: smartCard.CreatedAt.Time.Format("2006-01-02 15:04:05"),
				IsActive:  smartCard.IsActive,
			},
			OwnerRole: ownerRole,
			Details: model.SmartCardDetails{
				ID:   detailsId,
				Name: detailsName,
			},
		})
	}

	return result, nil
}

func (c *service) CountSmartCards(ctx context.Context, request *model.ListSmartCardRequest) (int64, error) {
	count, err := c.store.CountSmartCards(ctx, db.CountSmartCardsParams{
		Q:          pgtype.Text{String: request.Q, Valid: request.Q != ""},
		IsActive:   pgtype.Bool{Bool: true, Valid: true},
		IsSantri:   pgtype.Bool{Bool: request.OwnerRole == "santri", Valid: request.OwnerRole != ""},
		IsEmployee: pgtype.Bool{Bool: request.OwnerRole == "employee", Valid: request.OwnerRole != ""},
	})

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *service) UpdateSmartCard(ctx context.Context, request *model.UpdateSmartCardRequest, id int32) (model.SmartCardComplete, error) {
	var ownerRole string
	var detailsId int32
	var detailsName string

	if request.OwnerRole == "santri" {
		santri, err := c.store.GetSantri(ctx, request.OwnerID)
		if err != nil {
			return model.SmartCardComplete{}, err
		}
		detailsId = santri.ID
		detailsName = santri.Name
		ownerRole = "santri"
	} else if request.OwnerRole == "employee" {
		employee, err := c.store.GetEmployee(ctx, request.OwnerID)
		if err != nil {
			return model.SmartCardComplete{}, err
		}
		detailsId = employee.ID
		detailsName = employee.Name
		ownerRole = "employee"
	} else {
		return model.SmartCardComplete{}, exception.NewValidationError("invalid owner role")
	}

	updatedSmartCard, err := c.store.UpdateSmartCard(ctx, db.UpdateSmartCardParams{
		ID:         id,
		IsActive:   pgtype.Bool{Bool: request.IsActive, Valid: true},
		SantriID:   pgtype.Int4{Int32: request.OwnerID, Valid: ownerRole == "santri"},
		EmployeeID: pgtype.Int4{Int32: request.OwnerID, Valid: ownerRole == "employee"},
	})

	if err != nil {
		return model.SmartCardComplete{}, err
	}

	return model.SmartCardComplete{
		SmartCard: model.SmartCard{
			ID:        updatedSmartCard.ID,
			Uid:       updatedSmartCard.Uid,
			CreatedAt: updatedSmartCard.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			IsActive:  updatedSmartCard.IsActive,
		},
		OwnerRole: string(ownerRole),
		Details: model.SmartCardDetails{
			ID:   detailsId,
			Name: detailsName,
		},
	}, nil
}

func (c *service) DeleteSmartCard(ctx context.Context, id int32) (model.SmartCard, error) {
	deletedSmartCard, err := c.store.DeleteSmartCard(ctx, id)
	if err != nil {
		return model.SmartCard{}, err
	}

	return model.SmartCard{
		ID:        deletedSmartCard.ID,
		Uid:       deletedSmartCard.Uid,
		CreatedAt: deletedSmartCard.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		IsActive:  deletedSmartCard.IsActive,
	}, nil
}
