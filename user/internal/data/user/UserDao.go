package user

import (
	"database/sql"
	"fmt"
	"github.com/go-practice/api/user"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init(){
	db, _ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/my_order")
	//设置数据库最大连接数
	db.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(10)
	//验证连接
	if err := db.Ping(); err != nil {
		fmt.Println("open database fail")
		return
	}
	fmt.Println("connnect success")
}

func Get(uid int64) user.UserEntry{
	rows := db.QueryRow("select id,`name` from `user` where `id` = ? limit 1",uid)
	var id int64
	var name string
	rows.Scan(&id,&name)
	return user.UserEntry{
		Id: id,
		Name: name,
	}
}