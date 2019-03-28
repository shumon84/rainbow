package rainbow

import (
	"crypto/md5"
	"encoding/gob"
	"go/token"
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

	plain := strings.Repeat(MessageChars[0:1], MessageLength)
	chainChan := make(chan *Chain, 4096)
	go AddTable(H, R, chainChan)
	for i := 0; i < int(numOfChains); i++ {
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
		plain = nextPermutation(plain)
	}

	close(chainChan)
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(rainbowTable); err != nil {
		return err
	}
	return nil
}

func AddTable(chainChan <-chan *Chain) {
	for chain := range chainChan {
		if _, ok := rainbowTable[chain.Tail]; !ok {
			rainbowTable[chain.Tail] = *chain
		}
	}
}

func Hash(message []byte) []byte {
	digest := md5.New()
	digest.Write(message)
	return digest.Sum(nil)
}

func Reduction(times int, digest []byte) []byte {
	seed := uint64(0)
	for i := 0; i < 8; i++ {
		seed |= uint64(digest[i%len(digest)]) << uint64(i*8)
	}
	seed += uint64(times)
	seed = xorshift(seed)
	random := uint64(0)

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

func xorshift(x uint64) uint64 {
	x = x ^ (x << 7)
	return x ^ (x >> 9)
}
