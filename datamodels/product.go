// Author : rexdu
// Time : 2020-03-22 23:55
package datamodels

type Product struct {
	ID           int64  `json:"id" sql:"ID" imooc:"id"`
	ProductName  string `json:"ProductName" sql:"productName" imooc:"ProductName"`
	ProductNum   int64  `json:"ProductNum" sql:"productNum" imooc:"ProductNum"`
	ProductImage string `json:"ProductImage" sql:"productImage" imooc:"ProductImage"`
	ProductUrl   string `json:"ProductUrl" sql:"productUrl" imooc:"ProductUrl"`
	Price        string `json:"Price" sql:"price" imooc:"Price"`
}
