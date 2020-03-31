// Author : rexdu
// Time : 2020-03-25 00:00
package repositories

import (
	"database/sql"
	"seckill/common"
	"seckill/datamodels"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(order *datamodels.Order) (orderID int64, err error)
	Delete(orderID int64) bool
	Update(order *datamodels.Order) error
	SelectByKey(orderID int64) (order *datamodels.Order, err error)
	SelectAll() (orders []*datamodels.Order, err error)
	SelectAllWithInfo() (OrderMap map[int]map[string]string, err error)
}

type OrderManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewOrderManager(table string, db *sql.DB) IOrderRepository {
	return &OrderManager{table: table, mysqlConn: db}
}

func (o *OrderManager) Conn() (err error) {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "order"
	}
	return
}

func (o *OrderManager) Insert(order *datamodels.Order) (orderID int64, err error) {
	if err = o.Conn(); err != nil {
		return
	}
	// 表名order是mysql的关键字，所以在sql语句中要写成`order`才不会报错
	sql := "INSERT `" + o.table + "` set userID=?,productID=?,orderStatus=?"
	stmt, errSql := o.mysqlConn.Prepare(sql)
	if errSql != nil {
		return 0, errSql
	}
	result, errStmt := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if errStmt != nil {
		return 0, errStmt
	}
	return result.LastInsertId()
}

func (o *OrderManager) Delete(orderID int64) bool {
	if err := o.Conn(); err != nil {
		return false
	}
	sql := "delete from " + o.table + " where ID=?"
	stmt, errSql := o.mysqlConn.Prepare(sql)
	if errSql != nil {
		return false
	}
	_, errStmt := stmt.Exec(orderID)
	if errStmt != nil {
		return false
	}
	return true
}

func (o *OrderManager) Update(order *datamodels.Order) (err error) {
	if err := o.Conn(); err != nil {
		return err
	}
	sql := "update " + o.table + " set userID=?,productID=?,orderStatus=? where ID=?"
	stmt, errSql := o.mysqlConn.Prepare(sql)
	if errSql != nil {
		return errSql
	}
	_, errStmt := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus, order.ID)
	if errStmt != nil {
		return errStmt
	}
	return
}

func (o *OrderManager) SelectByKey(orderID int64) (order *datamodels.Order, err error) {
	var rows *sql.Rows
	if err = o.Conn(); err != nil {
		return &datamodels.Order{}, err
	}
	sql := "select * from " + o.table + " where ID=" + strconv.FormatInt(orderID, 10)
	if rows, err = o.mysqlConn.Query(sql); err != nil {
		return &datamodels.Order{}, err
	}
	defer rows.Close()
	result := common.GetResultRow(rows)
	if len(result) == 0 {
		return &datamodels.Order{}, nil
	}
	order = &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return
}

func (o *OrderManager) SelectAll() (orders []*datamodels.Order, err error) {
	var rows *sql.Rows
	if err = o.Conn(); err != nil {
		return nil, err
	}
	sql := "select * from " + o.table
	if rows, err = o.mysqlConn.Query(sql); err != nil {
		return nil, err
	}
	defer rows.Close()
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}
	for _, v := range result {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orders = append(orders, order)
	}
	return
}

func (o *OrderManager) SelectAllWithInfo() (OrderMap map[int]map[string]string, err error) {
	var rows *sql.Rows
	if err = o.Conn(); err != nil {
		return nil, err
	}
	sql := "select o.ID, p.productName,o.orderStatus from seckill.order as o left join product as p on o.productID=p.ID "
	if rows, err = o.mysqlConn.Query(sql); err != nil {
		return nil, err
	}
	defer rows.Close()
	return common.GetResultRows(rows), err
}
