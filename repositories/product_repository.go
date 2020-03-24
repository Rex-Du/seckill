// Author : rexdu
// Time : 2020-03-23 00:02
package repositories

import (
	"database/sql"
	"seckill/common"
	"seckill/datamodels"
	"strconv"
)

// 1.开发对应的接口
type IProduct interface {
	Conn() error
	Insert(product *datamodels.Product) (int64, error)
	Delete(id int64) bool
	Update(product *datamodels.Product) error
	SelectByKey(id int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}

// 2.实现定义的接口
type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

// 实现ProductManager的构造方法，也可以用于自检ProductManager是否实现了IProduct接口
func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: db}
}

func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}
	return
}

func (p *ProductManager) Insert(product *datamodels.Product) (productID int64, err error) {
	if err = p.Conn(); err != nil {
		return
	}
	sql := "insert into product set productName=?,productNum=?,productImage=?,productUrl=?"
	stmt, errSql := p.mysqlConn.Prepare(sql)
	if errSql != nil {
		return 0, errSql
	}
	result, errStmt := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if errStmt != nil {
		return 0, errStmt
	}
	return result.LastInsertId()

}

func (p *ProductManager) Delete(id int64) bool {
	if err := p.Conn(); err != nil {
		return false
	}
	sql := "delete from product where ID=?"
	stmt, errSql := p.mysqlConn.Prepare(sql)
	if errSql != nil {
		return false
	}
	_, errStmt := stmt.Exec(id)
	if errStmt != nil {
		return false
	}
	return true
}

func (p *ProductManager) Update(product *datamodels.Product) (err error) {
	if err := p.Conn(); err != nil {
		return err
	}
	sql := "update product set productName=?,productNum=?,productImage=?,productUrl=? where ID=?"
	stmt, errSql := p.mysqlConn.Prepare(sql)
	if errSql != nil {
		return errSql
	}
	_, errStmt := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl, product.ID)
	if errStmt != nil {
		return errStmt
	}
	return
}

func (p *ProductManager) SelectByKey(id int64) (product *datamodels.Product, err error) {
	var rows *sql.Rows
	if err = p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}
	sql := "select * from " + p.table + " where ID=" + strconv.FormatInt(id, 10)
	if rows, err = p.mysqlConn.Query(sql); err != nil {
		return &datamodels.Product{}, err
	}
	defer rows.Close()
	result := common.GetResultRow(rows)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}
	product = &datamodels.Product{}
	common.DataToStructByTagSql(result, product)
	return
}

func (p *ProductManager) SelectAll() (products []*datamodels.Product, err error) {
	var rows *sql.Rows
	if err = p.Conn(); err != nil {
		return nil, err
	}
	sql := "select * from " + p.table
	if rows, err = p.mysqlConn.Query(sql); err != nil {
		return nil, err
	}
	defer rows.Close()
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}
	for _, v := range result {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		products = append(products, product)
	}
	return
}
