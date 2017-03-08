package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type TestMysql struct {
	db *sql.DB
}

/* 初始化数据库引擎 */
func Init() (*TestMysql, error) {
	test := new(TestMysql)
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/soft_logs?charset=utf8")
	//第一个参数 ： 数据库引擎
	//第二个参数 : 数据库DSN配置。Go中没有统一DSN,都是数据库引擎自己定义的，因此不同引擎可能配置不同
	//本次演示采用http://code.google.com/p/go-mysql-driver
	if err != nil {
		fmt.Println("database initialize error : ", err.Error())
		return nil, err
	}
	test.db = db
	return test, nil
}

/* 测试数据库数据添加 */
func (test *TestMysql) Create() {
	if test.db == nil {
		return
	}
	stmt, err := test.db.Prepare("insert into test(name,age)values(?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmt.Close()
	if result, err := stmt.Exec("张三", 20); err == nil {
		if id, err := result.LastInsertId(); err == nil {
			fmt.Println("insert id : ", id)
		}
	}
	if result, err := stmt.Exec("李四", 30); err == nil {
		if id, err := result.LastInsertId(); err == nil {
			fmt.Println("insert id : ", id)
		}
	}
	if result, err := stmt.Exec("王五", 25); err == nil {
		if id, err := result.LastInsertId(); err == nil {
			fmt.Println("insert id : ", id)
		}
	}
}

/* 测试数据库数据更新 */
func (test *TestMysql) Update() {
	if test.db == nil {
		return
	}
	stmt, err := test.db.Prepare("update test set name=?,age=? where age=?")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmt.Close()
	if result, err := stmt.Exec("周七", 40, 25); err == nil {
		if c, err := result.RowsAffected(); err == nil {
			fmt.Println("update count : ", c)
		}
	}
}

/* 测试数据库数据读取 */
func (test *TestMysql) Read() {
	if test.db == nil {
		return
	}
	rows, err := test.db.Query("select id,name,age from test limit 0,5")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()
	fmt.Println("")
	cols, _ := rows.Columns()
	for i := range cols {
		fmt.Print(cols[i])
		fmt.Print("\t")
	}
	fmt.Println("")
	var id int
	var name string
	var age int
	for rows.Next() {
		if err := rows.Scan(&id, &name, &age); err == nil {
			fmt.Print(id)
			fmt.Print("\t")
			fmt.Print(name)
			fmt.Print("\t")
			fmt.Print(age)
			fmt.Print("\t\r\n")
		}
	}
}

/* 测试数据库删除 */
func (test *TestMysql) Delete() {
	if test.db == nil {
		return
	}
	stmt, err := test.db.Prepare("delete from test where age=?")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmt.Close()
	if result, err := stmt.Exec(20); err == nil {
		if c, err := result.RowsAffected(); err == nil {
			fmt.Println("remove count : ", c)
		}
	}
}

func (test *TestMysql) Close() {
	if test.db != nil {
		test.db.Close()
	}
}

func (test *TestMysql) aa() (gameIds []int) {
	rows, err := test.db.Query("select game_id FROM netbar_game_log_20151123 where netbar_id=579 limit 10")
	if err != nil {
		panic(err.Error())
		return
	}
	defer rows.Close()

	var gameId int

	for rows.Next() {
		if err := rows.Scan(&gameId); err == nil {
			fmt.Println("aa")
			gameIds = append(gameIds, gameId)
		} else {
			panic(err)
		}
	}
	fmt.Println("bbb")
	return
}

func main() {
	if test, err := Init(); err == nil {
		//		test.Create()
		//		test.Update()
		//		test.Read()
		//		test.Delete()
		//		test.Read()
		//		test.Close()
		test.aa()
	}
}
