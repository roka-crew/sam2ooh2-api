package payload

import "github.com/roka-crew/sam2ooh2-api/pkg/query"

// store layer
type CreateGroupParam struct {
	Title       string
	Author      string
	PageCount   int
	Publisher   string
	Description string
}

type ListGroupsParam struct {
	// conditions
	IDs          []uint
	Titles       []string
	Authors      []string
	PageCounts   []int
	Publishers   []string
	Descriptions []string

	// sort
	Sorts query.Sorts

	// options
	Limit  int
	Offset int
}

type PatchGroupParam struct {
	// conditions
	ID uint

	// fields
	Title       string
	Author      string
	PageCount   int
	Publisher   string
	Description string
}

type DeleteGroupParam struct {
	// conditions
	ID uint

	// options
	IsHardDelete bool
}

type AddUserToGroupParam struct {
	UserID  uint
	GroupID uint
}

type RemoveUserFromGroup struct {
	UserID  uint
	GroupID uint
}

type CountGroupUsers struct {
	GroupID uint
}

type IsUserInGroup struct {
	UserID  uint
	GroupID uint
}
