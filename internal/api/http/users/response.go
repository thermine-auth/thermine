package users

import (
	"xermess/internal/model"
	"xermess/internal/store"
)

// pageResponse is a page of users, with enough about the page for the panel
// to say "1 of 4" without asking again.
type pageResponse struct {
	Users  []model.User `json:"users"`
	Total  int64        `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

func newPageResponse(users []model.User, total int64, query store.UserQuery) pageResponse {
	return pageResponse{
		Users:  users,
		Total:  total,
		Limit:  query.Limit,
		Offset: query.Offset,
	}
}
