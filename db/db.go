package db

import (
	"database/sql"
	"fmt"
	"mygoshop/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDbConn() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.SQLSet.User,
		config.SQLSet.Password,
		config.SQLSet.Host,
		config.SQLSet.Port,
		config.SQLSet.Dbname,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}
	return db, nil
}

// 获取一条返回值
func GetResultRow(rows *sql.Rows) map[string]string {
	colums, _ := rows.Columns()
	scanArgs := make([]interface{}, len(colums))
	values := make([][]byte, len(colums))
	for index := range values {
		scanArgs[index] = &values[index]
	}
	result := make(map[string]string)
	for rows.Next() {
		rows.Scan(scanArgs...)
		for index, val := range values {
			result[colums[index]] = string(val)
		}
	}
	return result
}

// 获取所有返回值
func GetAllResult(rows *sql.Rows) map[int]map[string]string {
	colums, _ := rows.Columns()
	scanArgs := make([]interface{}, len(colums))
	values := make([][]byte, len(colums))
	for index := range values {
		scanArgs[index] = &values[index]
	}
	result := make(map[int]map[string]string)
	result_index := 0
	for rows.Next() {
		rows.Scan(scanArgs...)
		for index, val := range values {
			result[result_index][colums[index]] = string(val)
		}
		result_index++
	}
	return result
}
