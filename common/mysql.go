// Author : rexdu
// Time : 2020-03-23 00:16
package common

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//创建mysql 连接
func NewMysqlConn() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "root:rootroot@tcp(192.168.124.129:3306)/seckill?charset=utf8&loc=Local")
	return
}

//获取返回值，获取一条
func GetResultRow(rows *sql.Rows) map[string]string {
	if rows == nil {
		fmt.Println("rows实例不能为nil")
		return nil
	}
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	record := make(map[string]string)
	for rows.Next() {
		//将行数据保存到record字典
		rows.Scan(scanArgs...)
		for i, v := range values {
			if v != nil {
				// fmt.Println(reflect.TypeOf(v))
				record[columns[i]] = string(v.([]byte))
			}
		}
	}
	return record
}

//获取所有
func GetResultRows(rows *sql.Rows) map[int]map[string]string {
	//返回所有列
	columns, _ := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(columns))
	//这里表示一行填充数据
	scans := make([]interface{}, len(columns))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for rows.Next() {
		//填充数据
		rows.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := columns[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result[i] = row
		i++
	}
	return result
}