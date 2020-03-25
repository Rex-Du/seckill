// Author : rexdu
// Time : 2020-03-25 23:22
package repositories

import (
	"database/sql"
	"errors"
	"seckill/common"
	"seckill/datamodels"
)

type IUserRepository interface {
	Conn() (err error)
	Select(userName string) (user *datamodels.User, err error)
	Insert(user *datamodels.User) (userID int64, err error)
}

type UserManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func NewUserRepository(table string, db *sql.DB) IUserRepository {
	return &UserManagerRepository{table: table, mysqlConn: db}
}

func (u *UserManagerRepository) Conn() (err error) {
	if u.mysqlConn == nil {
		mysql, errMysql := common.NewMysqlConn()
		if errMysql != nil {
			return errMysql
		}
		u.mysqlConn = mysql
	}
	if u.table == "" {
		u.table = "user"
	}
	return
}

func (u *UserManagerRepository) Select(userName string) (user *datamodels.User, err error) {
	if userName == "" {
		return &datamodels.User{}, errors.New("条件不能为空！")
	}
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "select * from " + u.table + " where userName=?"

	rows, err := u.mysqlConn.Query(sql, userName)
	defer rows.Close()
	if err != nil {
		return &datamodels.User{}, err
	}
	result := common.GetResultRow(rows)

	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}

func (u *UserManagerRepository) Insert(user *datamodels.User) (userID int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	sql := "insert into " + u.table + " set nickName=?,userName=?,passWord=?"
	stmt, err := u.mysqlConn.Prepare(sql)
	if err != nil {
		return userID, err
	}

	result, err := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if err != nil {
		return userID, err
	}
	return result.LastInsertId()
}

func (u *UserManagerRepository) SelectByID(userID int64) (user *datamodels.User, err error) {
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "select * from " + u.table + " where userID=?"

	rows, err := u.mysqlConn.Query(sql, userID)
	defer rows.Close()
	if err != nil {
		return &datamodels.User{}, err
	}
	result := common.GetResultRow(rows)

	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}
