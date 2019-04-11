#!/bin/bash
NUM_OF_CHAINS=20000
CHAIN_LENGTH=5000
MESSAGE_CHARS_LENGTH=64
MESSAGE_LENGTH=4

go build -o create_table &&
./create_table &&

LINE=$(cat rainbow_table_"$MESSAGE_LENGTH"_"$NUM_OF_CHAINS"_"$CHAIN_LENGTH".txt | cut -d " " -f 1 | sort | uniq | wc -l)

echo 重複しなかったのは $NUM_OF_CHAINS 個中 $LINE 個のチェーンです
echo ユニークなハッシュ数は $(($LINE*$CHAIN_LENGTH)) です
python -c "print(\"独立チェーン生成成功率は{0:.2f}%でした\".format($LINE/$NUM_OF_CHAINS*100))"
python -c "print(\"網羅率は{0:.10f}%でした\".format($LINE*$CHAIN_LENGTH*100/$MESSAGE_CHARS_LENGTH**$MESSAGE_LENGTH))"
