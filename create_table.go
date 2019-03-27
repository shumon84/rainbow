package rainbow

import (
	"crypto/md5"
	"fmt"
	"os"
	"strings"
	"sync"
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
