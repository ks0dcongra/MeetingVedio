package repository

// import (
// 	"time"

// 	"gorm.io/gorm"
// )

type UserRepositoryInterface interface {
}
type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create User
// func (h *UserRepository) Create(data *model.Student) (id int, result *gorm.DB) {
// 	student := model.Student{
// 		Name:           data.Name,
// 		Password:       data.Password,
// 		Student_number: data.Student_number,
// 		CreatedTime:    time.Now(),
// 		UpdatedTime:    time.Now()}
// 	result = database.DB.Create(&student)
// 	return student.Id, result
// }
