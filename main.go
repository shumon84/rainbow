package main

import (
	"crypto/md5"
	"fmt"
	"io"
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

// ハッシュ関数
var HTime = time.Duration(0)

func H(message []byte) []byte {
	startTime := time.Now()
	defer func() {
		HTime += time.Since(startTime)
	}()
	digest := md5.New()
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

func main2() {
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

func ReHash(rainbow_table io.Reader,hash []byte)[]byte {
	const numOfChain= 10000 // チェーン数
	const chainLength= 3000 // チェーン長
	//candidateList := [chainLength]string{}
	candidateList := make(chan string, 4096)
	//endList := [numOfChain]string{}

	// チェーンの復元
	fmt.Println("チェーン複合中")
	numOfEndGoroutine :=0
	for i := 0; i < chainLength; i++ {
		go func(i int) {
			rx := R(chainLength-i-1, hash)
			for j := 0; j < i; j++ {
				beforeHash := H(rx)
				rx = R(chainLength-i+j, beforeHash)
			}
			candidateList <- string(rx)
			numOfEndGoroutine++
		}(i)
	}
	fmt.Println("チェーン複合終了")

	// テーブルの読み込み
	//fmt.Println("テーブル読み込み中")
	//scanner := bufio.NewScanner(rainbow_table)
	//for i := 0; scanner.Scan(); i++ {
	//	line := scanner.Text()
	//	endList[i] = line[9:]
	//	fmt.Println(i,endList[i])
	//}
	//fmt.Println("テーブル読み込み終了")

	fmt.Println("待ち合わせ中")
	for numOfEndGoroutine != chainLength {}
	fmt.Println("待ち合わせ終了")

	// テーブルとチェーンの照合
	fmt.Println("テーブルとチェーンの照合中")
	answer := ""
	//for i := 0; i < numOfChain; i++ {
	//	for j := 0; j < chainLength; j++ {
	//		if RainbowTable[i] == candidateList[j] {
	//			answer = candidateList[j]
	//		}
	//	}
	//}
	for candidate := range candidateList{
		go func(candidate string){
			for _,rainbow := range RainbowTable{
				if rainbow == candidate{
					fmt.Println(candidate)
				}
			}
		}(candidate)
	}

	return []byte(answer)
}

func main() {
	//fileName := os.Args[1]
	//file,err := os.Open(fileName)
	//if err != nil{
	//	log.Fatal(err)
	//}
	hash := H([]byte("90000000"))

	plain := ReHash(nil,hash)
	fmt.Println(string(plain))
}