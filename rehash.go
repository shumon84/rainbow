package rainbow

import (
	"bufio"
	"encoding/hex"
	"io"
	"strings"
)

var (
	RainbowTable = map[string]string{}
)

const NotFound = "Not Found"

func ReadRainbowTable(rainbowTableReader io.Reader) error {
	// レインボーテーブルの読み込み
	scanner := bufio.NewScanner(rainbowTableReader)
	for scanner.Scan() {
		line := scanner.Text()
		chain := strings.Split(line, " ")
		if len(chain) != 2 {
			break
		}
		RainbowTable[chain[0]] = chain[1]
	}
	return nil
}

func ReHash(hash []byte, H HashFunc, R ReductionFunc) []byte {
	// チェーンの復元
	tailChan := make(chan string, 4096)
	for i := 0; i < ChainLength; i++ {
		go func(i int) {
			rx := R(ChainLength-i-1, hash)
			for j := 0; j < i; j++ {
				beforeHash := H(rx)
				rx = R(ChainLength-i+j, beforeHash)
			}
			tailChan <- string(rx)
		}(i)
	}

	// テーブルとチェーンの照合
	answer := make(chan []byte)
	startGoroutine := make(chan int, 1024)
	endGoroutine := make(chan int, 1024)
	egCount := 0
	sgCount := 0

	for {
		select {
		case tail := <-tailChan:
			go func(tail string) {
				startGoroutine <- 1
				defer func() { endGoroutine <- 1 }()
				head, ok := RainbowTable[tail]
				if !ok {
					return
				}
				Rx := []byte(head)
				hexHash := hex.EncodeToString(hash)
				var Hx []byte
				for j := 0; j < ChainLength; j++ {
					Hx = H(Rx)
					if hex.EncodeToString(Hx) == hexHash {
						answer <- Rx
					}
					Rx = R(j, Hx)
				}
			}(tail)
		case ans := <-answer:
			return ans
		case n := <-startGoroutine:
			sgCount += n
		case n := <-endGoroutine:
			egCount += n
		default:
			if sgCount == ChainLength && egCount == ChainLength {
				return []byte(NotFound)
			}
		}
	}
}
