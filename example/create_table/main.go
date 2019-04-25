package main

import (
	"fmt"
	"github.com/shumon84/rainbow"
)

func main() {
	fileName := fmt.Sprintf("rainbow_table_%d_%d_%d.txt", rainbow.MessageLength, rainbow.NumOfChains, rainbow.ChainLength)
	fmt.Println(fileName + "にレインボーテーブルを書き込みます")
	rainbow.CreateTable(rainbow.Hash, rainbow.Reduction, fileName)
}
