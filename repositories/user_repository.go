package repositories

import (
	"errors"
	"mygoshop/datamodels"
	"mygoshop/db"

	"gorm.io/gorm"
)

type IUser interface {
	Conn() error
	Insert(user *datamodels.User) (int64, error)
	Delete(int64) bool
	Update(user *datamodels.User) error
	SelectById(id int64) (user *datamodels.User, err error)
	Select(userName string) (user *datamodels.User, err error)
}

type UserManager struct {
	sqlConn *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUser {
	return &UserManager{
		sqlConn: db,
	}
}

func (u *UserManager) Conn() error {
	if u.sqlConn == nil {
		db, err := db.NewDbConn()
		if err != nil {
			return err
		}
		u.sqlConn = db
	}
	return nil
}

func (u *UserManager) Insert(user *datamodels.User) (userID int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}

	err = u.sqlConn.Create(user).Error
	if err != nil {
		return -1, err
	}
	return user.ID, nil
}

func (u *UserManager) Delete(int64) bool {
	// TODO:
	panic("implement me")
}

func (u *UserManager) Update(user *datamodels.User) error {
	//TODO:
	panic("implement me")
}

func (u *UserManager) Select(userName string) (*datamodels.User, error) {
	if userName == "" {
		return nil, errors.New("条件不能为空！")
	}

	if err := u.Conn(); err != nil {
		return nil, err
	}

	user := &datamodels.User{}
	err := u.sqlConn.Where("userName=?", userName).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserManager) SelectById(userID int64) (user *datamodels.User, err error) {
	if err = u.Conn(); err != nil {
		return nil, err
	}

	err = u.sqlConn.Where("ID=?", userID).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
