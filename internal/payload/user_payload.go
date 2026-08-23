package payload

import "github.com/roka-crew/sam2ooh2-api/pkg/query"

// to service layer

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
