// Author : rexdu
// Time : 2020-03-25 23:19
package datamodels

type User struct {
	ID           int64  `json:"id" form:"ID" sql:"ID"`
	NickName     string `json:"nickName" form:"nickName" sql:"nickName"`
	UserName     string `json:"userName" form:"userName" sql:"userName"`
	HashPassword string `json:"-" form:"password" sql:"password"`
}
