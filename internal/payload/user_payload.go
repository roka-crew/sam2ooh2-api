package payload

import "github.com/roka-crew/sam2ooh2-api/pkg/query"

// to service layer
type CreateUserRequest struct {
	Nickname  string  `json:"nickname"  validate:"required,min=2,max=12"`
	Biography *string `json:"biography" validate:"max=14"`
}

type CreateUserResponse struct {
	UserID    uint   `json:"userID"`
	Nickname  string `json:"nickname"`
	Biography string `json:"biography"`
}

// to store layer
type CreateUserParam struct {
	Nickname  string
	Biography *string
}

type ListUsersParam struct {
	// conditions
	IDs         []uint
	Nicknames   []string
	Biographies []string

	// sort
	Sorts query.Sorts

	// optinos
	Limit  int
	Offset int
}

type PatchUsersParam []PatchUserParam

type PatchUserParam struct {
	// conditions
	ID uint

	// fields
	Nickname  *string
	Biography *string
}

type DeleteUserParam struct {
	// conditions
	ID       uint
	Nickname string

	// option
	IsHardDelete bool
}
