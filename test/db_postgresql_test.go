package test

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/sql_factory"
	_ "goskeleton/bootstrap"
	"testing"
)

//	测试 postgre 之前，首先请去  app/utils/sql_factory/client.go 第 8 行， 打开 被注释的驱动，否则 postgre 无法操作
//  database/db_demo_postgre.sql 有最简洁的创建表命令,您可以快速初始化一个 db_goskeleton 数据库以及 scheme，例如： web , 然后快速使用demo文件创建相关表.
// 本次测试使用最快捷的方式，只要保证 postgre 驱动初始化 ok 以及连接有效即可
// 实际应用请在 app/model 里面建表，整个操作与 mysql 类似

// 查询类
func TestSelectPostgre(t *testing.T) {

	postgreConn := sql_factory.GetOneSqlClient("postgre", "Read")
	if postgreConn == nil {
		return
	}
	sql := "SELECT name, sex, age, addr, remark, created_at, updated_at FROM web.tb_users "
	rows, err := postgreConn.Query(sql)
	if err == nil {
		var userName, addr, sex, age, remark, createdAt, updatedAt string
		for rows.Next() {
			_ = rows.Scan(&userName, &sex, &age, &addr, &remark, &createdAt, &updatedAt)
			fmt.Println(userName, sex, age, addr, remark, createdAt, updatedAt)
		}
		_ = rows.Close()
	} else {
		fmt.Println(err.Error())
	}
}

//执行类： 以修改数据为例，其他类似
func TestUpdatePostgre(t *testing.T) {

	postgreConn := sql_factory.GetOneSqlClient("postgre", "Write")
	if postgreConn == nil {
		return
	}
	sql := "update    web.tb_users   set  created_at=current_date ,updated_at=current_date ,remark='数据修改测试,postgre'  where   id=3  "
	result, err := postgreConn.Exec(sql)
	if err == nil {
		effectiveNum, err2 := result.RowsAffected()
		if err2 == nil {
			fmt.Println("修改数据音响行数：", effectiveNum)
		} else {
			t.Errorf("修改数据发生错误：%s", err2.Error())
		}

	} else {
		t.Errorf("执行sql发生错误：%s", err.Error())
	}
}

// 测试读写分离
// 您可以在config/config.yml> PostgreSql 配置不同的数据库ip等信息，测试  读 和 写 所操作的数据库

func TestPostgreReadWrite(t *testing.T) {
	sqlservConn := model.CreateBaseSqlFactory("postgresql")
	fmt.Printf("获取sql数据库的指针:%#+v\n", sqlservConn)
	if sqlservConn == nil {
		t.Error("单元测试失败")
	}
	sql := "update   web.tb_users   set  created_at=current_date ,updated_at=current_date ,remark='数据修改测试_postgresql'  where   id=3  "
	// 这里的操作会在 Write 对应的数据库进行
	effectiveRowNums := sqlservConn.ExecuteSql(sql)
	fmt.Println("影响的行数：", effectiveRowNums)

	sql = "select    name, sex, age, addr, remark, created_at, updated_at from  web.tb_users "
	// 这里的操作会在 Read 对应的数据库进行
	rows := sqlservConn.QuerySql(sql)
	if rows != nil {
		var userName, sex, age, addr, remark, createdAt, updatedAt string
		for rows.Next() {
			_ = rows.Scan(&userName, &sex, &age, &addr, &remark, &createdAt, &updatedAt)
			fmt.Println(userName, sex, age, addr, remark, createdAt, updatedAt)
			variable.ZapLog.Sugar().Info(userName, sex, age, addr, remark, createdAt, updatedAt)
		}
		_ = rows.Close()
	} else {
		fmt.Println("获取postgresql连接失败")
	}

}
