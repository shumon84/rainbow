package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const debug = false

const (
	LowerAlphabet = "abcdefghijklmnopqrstuvwxyz"
	UpperAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Alphabet      = LowerAlphabet + UpperAlphabet
	Numeric       = "0123456789"
	AlphaNumeric  = Alphabet + Numeric
)

const (
	AC = AlphaNumeric + "!#$%&'*+-/=?^_`{|}~"
	BC = AC + "(),.:;<>@[]"
)

var (
	RegAC = regexp.MustCompile("[" + AC + "]")
)

func IsContains(s []byte, b byte) bool {
	for _, v := range s {
		if v == b {
			return true
		}
	}
	return false
}

type StateMachine struct {
	io.Reader
}

func NewStateMachine(reader io.Reader) *StateMachine {
	return &StateMachine{
		Reader: reader,
	}
}

func (stateMachine *StateMachine) Next() (byte, error) {
	state := make([]byte, 1)
	_, err := stateMachine.Read(state)
	if err != nil {
		return 0, err
	}
	return state[0], nil
}

func (stateMachine *StateMachine) IsAcceptance() (bool, error) {
	return stateMachine.Q3()
}

func (stateMachine *StateMachine) Q1() (bool, error) {
	if debug {
		fmt.Print("-> Q1 ")
	}
	const isAcceptance = true
	current, err := stateMachine.Next()
	if err != nil && err != io.EOF {
		return false, err
	}
	if err == io.EOF {
		return isAcceptance, nil
	}
	if current == '.' {
		return stateMachine.Q4()
	} else if IsContains([]byte(AC), current) {
		return stateMachine.Q1()
	} else {
		return stateMachine.Q2()
	}
}

func (stateMachine *StateMachine) Q2() (bool, error) {
	if debug {
		fmt.Print("-> Q2 ")
	}
	const isAcceptance = false
	_, err := stateMachine.Next()
	if err != nil && err != io.EOF {
		return false, err
	}
	if err == io.EOF {
		return isAcceptance, nil
	}

	return stateMachine.Q2()
}

func (stateMachine *StateMachine) Q3() (bool, error) {
	if debug {
		fmt.Print("-> Q3 ")
	}
	const isAcceptance = false

	current, err := stateMachine.Next()
	if err != nil && err != io.EOF {
		return false, err
	}
	if err == io.EOF {
		return isAcceptance, nil
	}

	if current == '.' {
		return stateMachine.Q2()
	} else if IsContains([]byte(AC), current) {
		return stateMachine.Q1()
	} else {
		return stateMachine.Q2()
	}
}

func (stateMachine *StateMachine) Q4() (bool, error) {
	if debug {
		fmt.Print("-> Q4 ")
	}
	const isAcceptance = false

	current, err := stateMachine.Next()
	if err != nil && err != io.EOF {
		return false, err
	}
	if err == io.EOF {
		return isAcceptance, nil
	}

	if current == '.' {
		return stateMachine.Q2()
	} else {
		return stateMachine.Q1()
	}
}

func DomainVerify(str string) bool {
	stateMachine := NewStateMachine(strings.NewReader(str))
	ok, err := stateMachine.IsAcceptance()
	if err != nil {
		panic(err)
	}
	return ok
}

func main() {
	cin := bufio.NewScanner(os.Stdin)
	for cin.Scan() {
		fmt.Println(DomainVerify(cin.Text()))
	}
}
