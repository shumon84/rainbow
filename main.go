package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
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
	ASCIIChars  = NumberChars + LowerChars + UpperChars
)

const (
	MessageChars  = ASCIIChars
	MessageLength = 8
)

type HashFunc func([]byte) []byte
type ReductionFunc func(int, []byte) []byte
type ChainLength int
type NumOfChains int

func CreateTable(H HashFunc, R ReductionFunc, t NumOfChains, m ChainLength, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	//w := make(chan string, 0)
	//go func() {
	//	for str := range w {
	//		n, err := io.WriteString(file, str)
	//		if err != nil {
	//			log.Print(n, err)
	//		}
	//	}
	//}()

	plain := strings.Repeat(MessageChars[0:1], MessageLength)
	for i := 0; i < int(t); i++ {
		beforePlain := plain
		go func() {
			Rx := []byte(beforePlain)
			var Hx []byte
			for j := 0; j < int(m); j++ {
				Hx = H(Rx)
				Rx = R(j, Hx)
			}
			fmt.Fprintln(file, beforePlain, string(Rx))
		}()
		plain = NextPermutation(plain)
	}
	return nil
}

func main() {
	t, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	m, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	fileName := fmt.Sprintf("rainbow_table_%d_%d.txt", t, m)
	if err := CreateTable(H, R, NumOfChains(t), ChainLength(m), fileName); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HTime : %15d\n", HTime.Nanoseconds())
	fmt.Printf("RTime : %15d\n", RTime.Nanoseconds())
	fmt.Printf("NTime : %15d\n", NextPermutationTime.Nanoseconds())

	for runtime.NumGoroutine() > 1 {
		time.Sleep(time.Second / 2)
	}
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
	sum := uint64(0)
	for i, v := range digest {
		sum += uint64(v) << uint64(8*i)
	}
	sum += uint64(times)
	src := rand.NewSource(int64(sum))
	availableByte := []byte(MessageChars)
	availableByteLength := int64(len(availableByte))

	str := make([]byte, MessageLength)
	for i := 0; i < len(str); i += 8 {
		r := src.Int63()
		rx := []int64{
			(r >> 0) & 0xff,
			(r >> 8) & 0xff,
			(r >> 16) & 0xff,
			(r >> 24) & 0xff,
			(r >> 32) & 0xff,
			(r >> 40) & 0xff,
			(r >> 48) & 0xff,
			(r >> 56) & 0xff,
		}
		str[i+0] = availableByte[rx[0]%availableByteLength]
		str[i+1] = availableByte[rx[1]%availableByteLength]
		str[i+2] = availableByte[rx[2]%availableByteLength]
		str[i+3] = availableByte[rx[3]%availableByteLength]
		str[i+4] = availableByte[rx[4]%availableByteLength]
		str[i+5] = availableByte[rx[5]%availableByteLength]
		str[i+6] = availableByte[rx[6]%availableByteLength]
		str[i+7] = availableByte[rx[7]%availableByteLength]
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
		return str[:index] + MessageChars[i+1:i+2] + str[index+1:]
	}

	return nextPermutation(str, index+1)
}
