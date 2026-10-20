// 命令 thsdk-demo 演示 thsdk-go 的主要公开 API：
// 认证（账号密码 / 二维码）、K 线、分时、盘口、快讯、问财选股。
//
// 运行（在 thsdk-go 工程根目录）：
//
//	go run ./examples/demo -username 你的账号 -password 你的密码
//
// 不提供账号密码时尝试环境变量 THS_USERNAME / THS_PASSWORD；
// 都没有时回退到临时账户（需要运行组件支持）。
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	thsdk "thsdk"
)

func main() {
	username := flag.String("username", "", "同花顺账号（可选）")
	password := flag.String("password", "", "同花顺密码（可选）")
	code := flag.String("code", "USHA600519", "完整证券代码（四位市场前缀 + 证券代码）")
	flag.Parse()

	user := *username
	if user == "" {
		user = os.Getenv("THS_USERNAME")
	}
	pass := *password
	if pass == "" {
		pass = os.Getenv("THS_PASSWORD")
	}

	if user != "" {
		ok, err := thsdk.Auth(user, pass, "")
		if err != nil {
			log.Fatalf("账号密码认证失败: %v", err)
		}
		fmt.Printf("账号密码认证完成: %v\n", ok)
	} else {
		ok, err := thsdk.AuthQRCode()
		if err != nil {
			log.Fatalf("二维码登录失败: %v", err)
		}
		fmt.Printf("二维码登录完成: %v\n", ok)
	}
	defer func() {
		_ = thsdk.Logout()
	}()

	price, err := thsdk.GetPrice(*code, &thsdk.GetPriceOptions{Count: 5})
	if err != nil {
		log.Fatalf("GetPrice: %v", err)
	}
	fmt.Printf("最近 5 根日 K（索引 %s）：\n%s\n", price.IndexName, price)

	intraday, err := thsdk.IntradayData(*code, nil)
	if err != nil {
		log.Fatalf("IntradayData: %v", err)
	}
	fmt.Printf("当日分时，共 %d 行\n", intraday.RowCount())

	book, err := thsdk.Depth(*code)
	if err != nil {
		log.Fatalf("Depth: %v", err)
	}
	fmt.Printf("五档盘口（前 5 行）：\n%s\n", firstRows(book, 5))

	flash, err := thsdk.News(&thsdk.NewsOptions{Page: 1, PageSize: 3})
	if err != nil {
		log.Fatalf("News: %v", err)
	}
	fmt.Printf("7x24 快讯（前 3 条）：\n%s\n", firstRows(flash, 3))

	wencai, err := thsdk.WencaiNLP("市盈率小于 10 的贵州股票", &thsdk.WencaiNLPOptions{Limit: 10})
	if err != nil {
		log.Fatalf("WencaiNLP: %v", err)
	}
	fmt.Printf("问财选股结果，共 %d 行：\n%s\n", wencai.RowCount(), firstRows(wencai, 5))
}

// firstRows 截取表格前 n 行便于打印。
func firstRows(table *thsdk.Table, n int) *thsdk.Table {
	rows := table.Rows
	if len(rows) > n {
		rows = rows[:n]
	}
	return &thsdk.Table{Columns: table.Columns, Rows: rows, Index: table.Index, IndexName: table.IndexName, Attrs: table.Attrs}
}
