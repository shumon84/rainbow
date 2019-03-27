package rainbow

type plainGenerator struct {
	DictOrder []rune
	Length    int
}

func NewPlainGenerator(dictOrder []rune, length int) *plainGenerator {
	return &plainGenerator{
		DictOrder: dictOrder,
		Length:    length,
	}
}

func (plainGenerator *plainGenerator) Generate() *plain {
	plain := &plain{
		DictOrder: plainGenerator.DictOrder,
		length:    plainGenerator.Length,
	}
	plain.text = plain.MinDictOrderText()
	return plain
}

type plain struct {
	text      string
	DictOrder []rune
	length    int
}

func (plain *plain) MinDictOrderText() string {
	textRunes := make([]rune, plain.length)
	for i := range textRunes {
		textRunes[i] = plain.MinDictOrderRune()
	}
	return string(textRunes)
}

func (plain *plain) DictIndex(r rune) (int, bool) {
	for i, v := range []rune(plain.text) {
		if v == r {
			return i, true
		}
	}
	return -1, false
}

func (plain *plain) MinDictOrderRune() rune {
	return plain.DictOrder[0]
}

func (plain *plain) MaxDictOrderRune() rune {
	return plain.DictOrder[len(plain.DictOrder)-1]
}

func (plain *plain) Text() string {
	return plain.text
}

// MessageCharsに含まれる文字だけを使った文字列のうち、
// strの次に辞書順最小の文字列を返す
func (plain *plain) nextPermutation() {
	reverseText := reverseString(plain.text)
	nextReverseText := next(reverseText, 0)
	plain.text = reverseString(nextReverseText)
}

func (plain *plain) reverse() {
	revStr := make([]rune, len(plain.text))
	for i, v := range plain.text {
		revStr[len(revStr)+^i] = v
	}
	plain.text = string(revStr)
}

func (plain *plain) next(index int) {
	// indexが大きすぎる場合は辞書順最小の文字列を返す
	if index >= plain.length {
		plain.text = plain.MinDictOrderText()
		return
	}

	plainRunes := []rune(plain.text)

	// 辞書順でplainRunes[index]が何番目かを取得
	dictIndex, ok := plain.DictIndex(plainRunes[index])
	if !ok {
		panic(nil)
	}

	plainRunes[index] = plain.DictOrder[(dictIndex+1)%len(plain.DictOrder)]
	plain.text = string(plainRunes)

	// 現在の桁が既に最大の場合は繰り上がり
	if plainRunes[index] == plain.MaxDictOrderRune() {
		plain.next(index + 1)
	}
}
