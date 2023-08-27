package service

import ()

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

// CreateUser
// func (h *UserService) CreateUser(data *model.Student) (student_id int, status string) {
// 	pwd := []byte(data.Password)
// 	hash := NewUserService().HashAndSalt(pwd)
// 	data.Password = hash
// 	student_id, db := repository.NewUserRepository().Create(data)
// 	if db.Error != nil {
// 		return -1, responses.Error
// 	}
// 	return student_id, responses.Success
// }
