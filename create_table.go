package rainbow

import (
	"crypto/sha256"
	"encoding/gob"
	"log"
	"os"
	"strings"
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

var rainbowTable = map[string]Chain{}

func CreateTable(H HashFunc, R ReductionFunc, numOfChains int, chainLength int, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	minPlain := strings.Repeat(MessageChars[0:1], MessageLength)
	plain := minPlain
	chainChan := make(chan *Chain, numOfChains)
	isDoneChan := make(chan bool)
	isDone := false
	go AddTable(chainChan, isDoneChan, numOfChains)
	for {
		go func(beforePlain string) {
			Rx := []byte(beforePlain)
			var Hx []byte
			for j := 0; j < int(chainLength); j++ {
				Hx = H(Rx)
				Rx = R(j, Hx)
			}
			chainChan <- &Chain{
				Length: chainLength,
				Head:   beforePlain,
				Tail:   string(Rx),
			}
		}(plain)
		select {
		case isDone = <-isDoneChan:
		default:
		}
		plain = nextPermutation(plain)
		if isDone || plain == minPlain {
			break
		}
	}

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(rainbowTable); err != nil {
		return err
	}
	return nil
}

func AddTable(chainChan <-chan *Chain, isDoneChan chan<- bool, numOfChains int) {
	for chain := range chainChan {
		if _, ok := rainbowTable[chain.Tail]; !ok {
			rainbowTable[chain.Tail] = *chain
		}
		rainbowTableLength := len(rainbowTable)
		if rainbowTableLength > numOfChains {
			isDoneChan <- true
			break
		}
		if rainbowTableLength%1000 == 0 {
			log.Println(rainbowTableLength, "まで完了")
		}
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
