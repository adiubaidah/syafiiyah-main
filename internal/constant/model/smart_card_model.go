package model

type SmartCardRequest struct {
	Uid string `json:"uid" validate:"required"`
}

type ListSmartCardRequest struct {
	OwnerRole string `form:"owner-role" binding:"omitempty,oneof=employee santri"`
	Q         string `form:"q"`
	Page      int32  `form:"page" binding:"omitempty,gte=1"`
	Limit     int32  `form:"limit" binding:"omitempty,gte=1"`
}

type UpdateSmartCardRequest struct {
	IsActive  bool   `json:"is_active"`
	OwnerRole string `json:"owner_role" binding:"omitempty,oneof=employee santri"`
	OwnerID   int32  `json:"owner_id"`
}

type SmartCard struct {
	ID        int32  `json:"id"`
	Uid       string `json:"uid"`
	CreatedAt string `json:"create_at"`
	IsActive  bool   `json:"is_active"`
}

type ListSmartCardResponse struct {
	Items      []SmartCardComplete `json:"items"`
	Pagination Pagination          `json:"pagination"`
}

type SmartCardComplete struct {
	SmartCard
	OwnerRole string           `json:"owned_by"`
	Details   SmartCardDetails `json:"details"`
}

type SmartCardDetails struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}
