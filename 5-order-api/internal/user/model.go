package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Phone     string `json:"phone" gorm:"uniqueIndex"`
	Code      int
	SessionId string `json:"sessionId"`
}

func NewUser(phone string) *User {
	user := &User{
		Phone: phone,
		Code:  3245,
	}
	user.GenerateSessionId()
	return user
}

func (user *User) GenerateSessionId() {
	user.SessionId = uuid.NewString()
}
