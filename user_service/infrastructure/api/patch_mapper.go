package api

import (
	events "github.com/dislinkt/common/saga/patch_user"
	"github.com/dislinkt/user_service/domain"
	"github.com/gofrs/uuid"
)

func mapPatchUser(user events.User) *domain.User {
	id, _ := uuid.FromString(user.Id)
	userD := &domain.User{
		Id:       id,
		Username: &user.Username,
		Private:  user.Private,
	}
	return userD
}
