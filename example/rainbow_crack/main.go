package main

import (
	"encoding/hex"
	"fmt"
	"github.com/shumon84/rainbow"
	"log"
	"os"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	rainbow.ReadRainbowTable(file)
	hash, err := hex.DecodeString(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	plain := rainbow.ReHash(hash, rainbow.Hash, rainbow.Reduction)
	fmt.Println(string(plain))
}
