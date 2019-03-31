package main

import (
	"fmt"
	"github.com/shumon84/rainbow-table"
	"log"
	"os"
	"strconv"
)

func main() {
	// 第一引数がチェーン数
	t, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	// 第二引数がチェーン長
	m, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	fileName := fmt.Sprintf("rainbow_table_%d_%d.txt", t, m)
	fmt.Println(fileName + "にレインボーテーブルを書き込みます")
	rainbow.CreateTable(rainbow.Hash, rainbow.Reduction, t, m, fileName)
}
