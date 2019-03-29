package main

import (
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/shumon84/rainbow-table"
	"log"
	"os"
	"time"
)

func ReHash(RainbowTable map[string]*rainbow.Chain, hash []byte, H rainbow.HashFunc, R rainbow.ReductionFunc) []byte {
	startTime := time.Now()
	chainLength := 100
	candidateList := make(chan string, 4096)

	// チェーンの復元
	fmt.Println("チェーン複合中")
	for i := 0; i < chainLength; i++ {
		go func(i int) {
			rx := R(chainLength-i-1, hash)
			for j := 0; j < i; j++ {
				beforeHash := H(rx)
				rx = R(chainLength-i+j, beforeHash)
			}
			candidateList <- string(rx)
		}(i)
	}
	fmt.Println("チェーン複合終了")

	// テーブルとチェーンの照合
	fmt.Println("テーブルとチェーンの照合中")
	answer := make(chan []byte)
	startGoroutine := make(chan int, 1024)
	endGoroutine := make(chan int, 1024)
	egCount := 0
	sgCount := 0

	////
	//RainbowTable := map[string]string{} // これは消す
	////
	for {
		select {
		case candidate := <-candidateList:
			go func(candidate string) {
				startGoroutine <- 1
				defer func() { endGoroutine <- 1 }()
				chain, ok := RainbowTable[candidate]
				if !ok {
					return
				}
				//fmt.Println(head, "にヒットしました")
				Rx := []byte(chain.Head)
				hexHash := hex.EncodeToString(hash)
				var Hx []byte
				for j := 0; j < chainLength; j++ {
					Hx = H(Rx)
					if hex.EncodeToString(Hx) == hexHash {
						answer <- Rx
					}
					Rx = R(j, Hx)
				}
			}(candidate)
		case ans := <-answer:
			return ans
		case n := <-startGoroutine:
			sgCount += n
		case n := <-endGoroutine:
			egCount += n
			//fmt.Println("起動済み:実行済み = ", sgCount, ":", egCount)
			if sgCount == chainLength && egCount == chainLength {
				return []byte("Not Found This Hash")
			}
		default:
			if time.Since(startTime) > 6*time.Second {
				return []byte("Not Found This Hash")
			}
		}
	}
}

func main() {
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// レインボーテーブルの読み込み
	var rainbowTable map[string]*rainbow.Chain
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&rainbowTable); err != nil {
		log.Fatal(err)
	}

	// ハッシュ値を16進数でデコード
	hash, err := hex.DecodeString(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	plain := ReHash(rainbowTable, hash, rainbow.Hash, rainbow.Reduction)
	fmt.Print(string(plain))
}
