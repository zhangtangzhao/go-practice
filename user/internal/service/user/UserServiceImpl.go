package user

import (
	"context"
	"github.com/go-practice/api/user"
	userDao "github.com/go-practice/user/internal/data/user"
)

type UserServiceImpl struct {

}

func (userService *UserServiceImpl) GetUser(ctx context.Context, r *user.Request) (*user.Response, error) {

	entity := userDao.Get(r.Id)

	resp := user.Response{
		Code:0,
		Message: "success",
		UserEntry: &entity,
	}
	return &resp,nil
}