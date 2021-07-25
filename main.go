package main

import (
	"database/sql"
	"fmt"
	"runtime"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// 用户模型, 表结构, 需要一个结构来接收查询结果集
	type User struct {
		Id   int32
		Name string
		Age  int8
	}

	// 保存用户信息列表
	var user User

	// 查询一条记录时, 不能使用类似if err := db.QueryRow().Scan(&...); err != nil {}的处理方式
	// 因为查询单条数据时, 可能返回var ErrNoRows = errors.New("sql: no rows in result set")该种错误信息
	// 而这属于正常错误, 不使用Wrap, 也不抛给上层
	err = db.QueryRow(`
        SELECT id,name,age WHERE id = ?
    `, 2).Scan(
		&user.Id, &user.Name, &user.Age,
	)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		// 使用该方式可以打印出运行时的错误信息, 该种错误是编译时无法确定的
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	fmt.Println(user.Id, user.Name, user.Age)
}
