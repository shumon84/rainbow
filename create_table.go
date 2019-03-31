package rainbow

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
)

// 文字種のエイリアス
const (
	NumberChars = "0123456789"
	UpperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowerChars  = "abcdefghijklmnopqrstuvwxyz"
)

// レインボーテーブルの設定
const (
	NumOfChains   = 20000
	ChainLength   = 5000
	MessageLength = 4
	MessageChars  = "-" + NumberChars + UpperChars + "_" + LowerChars // 記号2文字 + 英数字 = 64文字(2^6)
)

type HashFunc func([]byte) []byte
type ReductionFunc func(int, []byte) []byte

func CreateTable(H HashFunc, R ReductionFunc, fileName string){
	// 初期文字列を生成
	minPlain := strings.Repeat(MessageChars[0:1], MessageLength)
	plain := minPlain

	// ファイルにレインボーテーブルを書き込むのに使うチャネルを生成し
	// ゴルーチンを起動させておく
	chainChan := make(chan string)
	isDoneChan := make(chan bool)
	go WriteTable(fileName, NumOfChains,chainChan,isDoneChan)

	// レインボーテーブルの生成
	for i:=0;i< NumOfChains;i++ {
		// チェインを生成するゴルーチンを起動
		go func(beforePlain string) {
			Rx := []byte(beforePlain)
			var Hx []byte
			for j := 0; j < int(ChainLength); j++ {
				Hx = H(Rx)
				Rx = R(j, Hx)
			}
			chainChan<-fmt.Sprintln(string(Rx),beforePlain)
		}(plain)

		// 次のHeadになる文字列を生成する
		// すでにHeadが取りうる全ての文字列を使用している場合は終了
		plain = nextPermutation(plain)
		if plain == minPlain {
			break
		}
	}

	// WriteTableの終了を待ってチャネルをcloseする
	<-isDoneChan
	close(chainChan)
	close(isDoneChan)
}

func WriteTable(fileName string,numOfLines int,lineChan <- chan string,isDoneChan chan<- bool){
	defer func(){isDoneChan<-true}()
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for i := 0;i<numOfLines;i++{
		fmt.Fprint(file,<-lineChan)
	}
}

func Hash(message []byte) []byte {
	digest := sha256.New()
	digest.Write(message)
	return digest.Sum(nil)
}

func Reduction(times int, digest []byte) []byte {
	seed := uint64(0)
	for i := 0; i < 8; i++ {
		seed |= uint64(digest[i%len(digest)]) << uint64(i*8)
	}
	seed += uint64(times)
	random := seed

	str := make([]byte, MessageLength)
	messageBytes := []byte(MessageChars)

	for i := 0; i < MessageLength; i++ {
		if random == 0 {
			random = seed
		}
		index := random & 0x3f
		str[i] = messageBytes[int(index)]
		random = random >> 6
	}
	return str
}
