package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	NumberChars = "0123456789"
	LowerChars  = "abcdefghijklmnopqrstuvwxyz"
	UpperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

const (
	MessageChars  = NumberChars+LowerChars+UpperChars+"-_" // 英数字 + 記号2文字 = 64文字(2^6)
	MessageLength = 8
)

type HashFunc func([]byte) []byte
type ReductionFunc func(int, []byte) []byte
type ChainLength int
type NumOfChains int

func CreateTable(H HashFunc, R ReductionFunc, t NumOfChains, m ChainLength, fileName string) (*os.File,error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return nil,err
	}

	plain := strings.Repeat(MessageChars[0:1], MessageLength)
	for i := 0; i < int(t); i++ {
		go func(beforePlain string) {
			Rx := []byte(beforePlain)
			var Hx []byte
			for j := 0; j < int(m); j++ {
				Hx = H(Rx)
				Rx = R(j, Hx)
			}
			fmt.Fprintln(file, beforePlain, string(Rx))
		}(plain)
		plain = NextPermutation(plain)
	}
	return file,nil
}

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
	file,err := CreateTable(H, R, NumOfChains(t), ChainLength(m), fileName)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	for runtime.NumGoroutine() > 1 {
		time.Sleep(time.Second / 2)
	}
	fmt.Println("ハッシュ関数 :",HTime)
	fmt.Println("　　還元関数 :",RTime)
	fmt.Println("順列生成関数 :",NextPermutationTime)
}

// ハッシュ関数
var HTime = time.Duration(0)

func H(message []byte) []byte {
	startTime := time.Now()
	defer func() {
		HTime += time.Since(startTime)
	}()
	digest := sha256.New()
	digest.Write(message)
	return digest.Sum(nil)
}

// 還元関数
var RTime = time.Duration(0)

func R(times int, digest []byte) []byte {
	startTime := time.Now()
	defer func() {
		RTime += time.Since(startTime)
	}()
	seed := uint64(0)
	for i, v := range digest {
		seed |= uint64(v) << uint64(i*8)
	}
	seed += uint64(times)
	availableByte := []byte(MessageChars)

	str := make([]byte, MessageLength)

	for i:=0;i<8;i++ {
		index := seed & 0x3f
		str[i] = availableByte[int(index)]
		seed = seed >> 8
	}
	return str
}

var NextPermutationTime = time.Duration(0)

func NextPermutation(str string) string {
	startTime := time.Now()
	defer func() {
		NextPermutationTime += time.Since(startTime)
	}()
	return nextPermutation(str, 0)
}

func nextPermutation(str string, index int) string {
	if index == MessageLength {
		return strings.Repeat(MessageChars[0:1], MessageLength)
	}

	i := strings.Index(MessageChars, str[index:index+1])
	if len(MessageChars)-1 > i {
		// 現在の桁を1増やす
		return str[:index] + MessageChars[i+1:i+2] + str[index+1:]
	}

	// 現在の桁を0に戻す
	str = str[:index] + MessageChars[:1] + str[index+1:]
	return nextPermutation(str, index+1)
}
