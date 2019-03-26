package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	NumberChars = "0123456789"
	LowerChars  = "abcdefghijklmnopqrstuvwxyz"
	UpperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

const (
	MessageChars  = NumberChars + LowerChars + UpperChars + "-_" // 英数字 + 記号2文字 = 64文字(2^6)
	MessageLength = 4
)

var availableByte = []byte(MessageChars)

type HashFunc func([]byte) []byte
type ReductionFunc func(int, []byte) []byte
type ChainLength int
type NumOfChains int

func CreateTable(H HashFunc, R ReductionFunc, t NumOfChains, m ChainLength, fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(file, "package main\n")
	fmt.Fprintln(file, "var RainbowTable = map[string]string{")

	plain := strings.Repeat(MessageChars[0:1], MessageLength)
	isExist := map[string]bool{}
	mutex := sync.Mutex{}
	for i := 0; i < int(t); i++ {
		go func(beforePlain string) {
			Rx := []byte(beforePlain)
			var Hx []byte
			for j := 0; j < int(m); j++ {
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
		plain = NextPermutation(plain)
	}
	return file, nil
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

	str := make([]byte, MessageLength)

	for i := 0; i < MessageLength; i++ {
		index := seed & 0x3f
		str[i] = availableByte[int(index)]
		seed = seed >> 6
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
	fileName := fmt.Sprintf("rainbow_table_%d_%d.go", t, m)
	file, err := CreateTable(H, R, NumOfChains(t), ChainLength(m), fileName)
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
	fmt.Println("ハッシュ関数 :", HTime)
	fmt.Println("　　還元関数 :", RTime)
	fmt.Println("順列生成関数 :", NextPermutationTime)
}

func ReHash(hash []byte) []byte {
	startTime := time.Now()
	const numOfChain = 10000 // チェーン数
	const chainLength = 3000 // チェーン長
	//candidateList := [chainLength]string{}
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
				head, ok := RainbowTable[candidate]
				if !ok {
					return
				}
				//fmt.Println(head, "にヒットしました")
				Rx := []byte(head)
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
	//fileName := os.Args[1]
	//file,err := os.Open(fileName)
	//if err != nil{
	//	log.Fatal(err)
	//}
	//hash := H([]byte("a1000000"))

	//fmt.Println(string(R(0, H([]byte("20000000")))))

	//Gen := RandomStringGenerator(MessageChars, MessageLength)
	//x := 0
	//for {
	//	x++
	//	os.Args[1] = hex.EncodeToString(H([]byte(Gen())))
	//
	hash, err := hex.DecodeString(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	plain := ReHash(hash)
	fmt.Print(string(plain))
	//	if string(plain) != "Not Found This Hash" {
	//		fmt.Println(string(plain), os.Args[1])
	//		break
	//	} else {
	//		fmt.Println(x, "回目の挑戦失敗 :", os.Args[1])
	//	}
	//}

}

// availableStringに含まれる文字のみを使って、
// 長さlengthのランダムな文字列を生成するジェネレータを返す高階関数
func RandomStringGenerator(availableString string, length int) func() string {
	src := rand.NewSource(time.Now().UnixNano())
	availableRune := []rune(availableString)
	availableRuneLength := int64(len(availableRune))
	return func() string {
		str := make([]rune, length)
		for i := range str {
			str[i] = availableRune[src.Int63()%availableRuneLength]
		}
		return string(str)
	}
}
