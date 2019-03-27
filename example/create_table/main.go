package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 文字種のエイリアス
const (
	NumberChars = "0123456789"
	LowerChars  = "abcdefghijklmnopqrstuvwxyz"
	UpperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// 平文に使える文字種と文字数の設定
const (
	MessageChars  = NumberChars + LowerChars + UpperChars + "-_" // 英数字 + 記号2文字 = 64文字(2^6)
	MessageLength = 6
)

type HashFunc func([]byte) []byte
type ReductionFunc func(int, []byte) []byte

// レインボーテーブルのチェーン
type Chain struct {
	Length int
	Head   string
	Tail   string
}

var RainbowTable = map[string]string{}

func CreateTable(H HashFunc, R ReductionFunc, numOfChains int, chainLength int, fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return nil, err
	}

	plain := strings.Repeat(MessageChars[0:1], MessageLength)
	isExist := map[string]bool{}
	mutex := sync.Mutex{}
	for i := 0; i < int(numOfChains); i++ {
		go func(beforePlain string) {
			Rx := []byte(beforePlain)
			var Hx []byte
			for j := 0; j < int(chainLength); j++ {
				Hx = H(Rx)
				Rx = R(j, Hx)
			}
			mutex.Lock()
			if !isExist[string(Rx)] {
				isExist[string(Rx)] = true
				fmt.Fprintf(file, "	\"%s\": \"%s\",\n", string(Rx), beforePlain)
			}
			mutex.Unlock()
		}(plain)
		plain = nextPermutation(plain)
	}
	return file, nil
}

func Hash(message []byte) []byte {
	digest := md5.New()
	digest.Write(message)
	return digest.Sum(nil)
}

func Reduction(times int, digest []byte) []byte {
	seed := uint64(0)
	for i, v := range digest {
		seed |= uint64(v) << uint64(i*8)
	}
	seed += uint64(times)

	str := make([]byte, MessageLength)
	messageBytes := []byte(MessageChars)

	for i := 0; i < MessageLength; i++ {
		index := seed & 0x3f
		str[i] = messageBytes[int(index)]
		seed = seed >> 6
	}
	return str
}

// MessageCharsに含まれる文字だけを使った文字列のうち、
// strの次に辞書順最小の文字列を返す
func nextPermutation(str string) string {
	revStr := reverseString(str)
	nextRevStr := next(revStr, 0)
	nextStr := reverseString(nextRevStr)
	return nextStr
}

func reverseString(str string) string {
	revStr := make([]rune, len(str))
	for i, v := range str {
		revStr[len(revStr)+^i] = v
	}
	return string(revStr)
}

func next(str string, index int) string {
	if index >= MessageLength {
		return strings.Repeat(MessageChars[0:1], MessageLength)
	}

	i := strings.Index(MessageChars, str[index:index+1])
	if len(MessageChars)-1 > i {
		// 現在の桁を1増やす
		return str[:index] + MessageChars[i+1:i+2] + str[index+1:]
	}

	// 現在の桁を0に戻す
	str = str[:index] + MessageChars[:1] + str[index+1:]
	return next(str, index+1)
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
	fileName := fmt.Sprintf("rainbow_table_%d_%d.go", t, m)
	file, err := CreateTable(Hash, Reduction, t, m, fileName)
	defer func() {
		fmt.Fprintln(file, "}")
		file.Close()
	}()
	if err != nil {
		log.Fatal(err)
	}

	for runtime.NumGoroutine() > 1 {
		time.Sleep(time.Second / 2)
	}
}
