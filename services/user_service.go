// Author : rexdu
// Time : 2020-03-25 23:46
package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"seckill/datamodels"
	"seckill/repositories"
)

type IUserService interface {
	IsPwdSuccess(userName, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userID int64, err error)
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func NewUserService(repo repositories.IUserRepository) IUserService {
	return &UserService{UserRepository: repo}
}

func (u *UserService) IsPwdSuccess(userName, pwd string) (user *datamodels.User, isOK bool) {
	var err error
	if userName == "" || pwd == "" {
		return &datamodels.User{}, false
	}
	user, err = u.UserRepository.Select(userName)
	if err != nil {
		return
	}
	isOK, _ = ValidatePassword(pwd, user.HashPassword)
	if !isOK {
		return &datamodels.User{}, false
	}
	return
}

func (u *UserService) AddUser(user *datamodels.User) (userID int64, err error) {
	hashed, err := GeneratePassword(user.HashPassword)
	if err != nil {
		return userID, err
	}
	user.HashPassword = string(hashed)
	return u.UserRepository.Insert(user)
}

func ValidatePassword(userPwd, hashed string) (isOK bool, err error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPwd)); err != nil {
		return false, errors.New("密码比对错误")
	}
	return true, nil
}

func GeneratePassword(pwd string) (hashed []byte, err error) {
	return bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
}
