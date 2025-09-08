package user

import "demo/order/6-order-api-cart/pkg/db"

type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(db *db.Db) *UserRepository {
	return &UserRepository{
		Database: db,
	}
}

func (repo *UserRepository) AddUser(user *db.User) (*db.User, error) {
	err := repo.Database.DB.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) FindByPhone(phone string) (*db.User, error) {
	var user db.User
	err := repo.Database.DB.Where(&db.User{Phone: phone}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetSessionId(phone string) (string, error) {
	user, err := repo.FindByPhone(phone)
	if err != nil {
		user := db.NewUser(phone)
		user, err = repo.AddUser(user)
		if err != nil {
			return "", err
		}
		return user.SessionId, nil
	}
	user.GenerateSessionId()
	err = repo.Database.DB.Save(user).Error
	if err != nil {
		return "", err
	}
	return user.SessionId, nil
}

func (repo *UserRepository) VerifyCode(sessionId string, code int) (*db.User, error) {
	var user db.User
	err := repo.Database.DB.Where(&db.User{SessionId: sessionId, Code: code}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
