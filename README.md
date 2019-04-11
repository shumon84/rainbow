# rainbow
レインボーテーブルの作成と、作成したテーブルを使ってハッシュ値からハッシュ化される前の情報を取得するライブラリです。

# Example
## レインボーテーブルの作成
```
$ cd example/create_table
$ go build
$ ./create_table
```

`example/create_table/collision_check.sh`を実行してもレインボーテーブルが作成できます

## レインボーテーブルを使ったハッシュ値の複合
事前に`create_table`によってレインボーテーブルが作成されており、
作成したテーブルが`example/rainbow_crack/rainbow_table_4_20000_5000.txt`に存在する場合、
以下の手順でハッシュ値からハッシュ化される前の情報を取得できます。
```
$ cd example/rainbow_crack
$ go build
$ ./rainbow_crack rainbow_table_4_20000_5000.txt 2f9acb02faa121bb2a3621951f57b4c690655337edee2e5ac350be2b3be88ea8
```


ちなみに

`2f9acb02faa121bb2a3621951f57b4c690655337edee2e5ac350be2b3be88ea8`

は `PASS` をsha256でハッシュ化したものです。

あらかじめ`example/rainbow_crack/convert.sh`を使ってレインボーテーブルから重複を削除しておくと、ファイルサイズが節約でき、起動も若干早くなります。

# To do
現在は末尾が衝突しているチェーンはレインボーテーブルから捨てているので、末尾が衝突しないように生成する方法を探るか、衝突していても活用するようにする