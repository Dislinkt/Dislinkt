package api

import (
	"fmt"

	events "github.com/dislinkt/common/saga/patch_user"
	"github.com/dislinkt/user_service/domain"
	uuid "github.com/satori/go.uuid"
)

func mapPatchUser(user events.User) *domain.User {
	id, _ := uuid.FromString(user.Id)
	userD := &domain.User{
		Id:       id,
		Username: &user.Username,
		Private:  user.Private,
	}
	fmt.Println("DOMAIN: ")
	fmt.Println(userD.Private)
	return userD
}
