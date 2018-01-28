今更ながらDockerの勉強をしています。
色々わからないところが多く、詰まったりしたところが多かったので自分はこのようにしたというのをまとめて見ます。output大事かなって思ったので

他にやり方があるよとかあれば、教えていただけると嬉しいです。

一応Dockerとか関係なく、
ローカルでさえ設定ができませえんって人が読めるくらいのレベルで書きたいです。

コードは全てgithubにおいておきます。
記事では関係ないところは飛ばしていたりします。
疑問があればご確認ください.
[この記事のコード](https://github.com/g-hyoga/test-golang-on-docker)

# この記事の目標

コマンドを一つ叩くと、
(データがなければ)初期化されたmysqlコンテナ、goのwebコンテナが立ち上がり、
今回はmysqlに入ったデータを取得し、jsonで返すくらいにしておきます。

この際、もちろんですが、データは永続化されていて、
ついでにmysqlのログがローカルで見れるようにします。

docker-composeが動く環境にだけしておいてください。


# 下準備:mysqlコンテナ

mysqlコンテナの下準備をします。
具体的には,

 * 文字コードの変更
 * logをローカルで見れるようにする
 * 初期データを作る

## 文字コードの変更

公式のmysql imageを使いますが、
そのままだと日本語を入れると文字化けしてしまいます.

解決策は調べてみた感じで複数あるようです。

 * mysqlの設定ファイルをVOLUMEを使って書き換える(これを選択)
 * mysqlの公式のコマンドを使う(試してない)

### mysqlの設定ファイルをVOLUMEを使って書き換える

僕はこれを使いました。
query logを見たくなった時とかにどうせ書き換える必要があるし

という訳で設定ファイルを用意します.
mysql系の設定ファイルはmysqlというでディレクトリを作ってそこに入れることにします。

##### ./mysql/my.cnf
```
...
[client]

...

default-character-set=utf8

...

[mysqld]

...

character-set-server=utf8

```

文字コードをutf8に変更しておきましょう。
寿司ビール問題とかはとりあえず無視でします。

### mysqlの公式のコマンドを使う(試してない)

見にくくなるの[こちら](https://github.com/g-hyoga/test-golang-on-docker/blob/master/blog/mysql-charaset.md)に書きました.


## logをローカルで見れるようにする

ここでは、とりあえずerr-logをローカルで見れるようにします。

そもそもエラーログを出力させるために
my.cnfを編集しましょう.

##### ./mysql/my.cnf
```
[mysqld]

...

log-error      = /var/log/mysql/error.log
```

次にローカルにログを出力するためのディレクトリを作っておきます。

```sh
mkdir mysql/log
```

ログの下準備は以上です。


## 初期データを作る

やり方は色々あるでしょうが、
ここでは、mysqlの公式コンテナでは,
コンテナ内の/docker-entrypoint-initdb.dというディレクトリにある
.sqlファイルと.shファイルをコンテナ作成時に読み込んでくれるというものを利用します。

ローカルに初期データ用のディレクトリを作っておき,
以下のようなsqlファイルを作っておきます.

##### ./mysql/init/0_init.sql
```sql
CREATE DATABASE test_db;

CREATE TABLE test_db.test_table(id int, name varchar(256));

INSERT INTO test_db.test_table VALUES(1, "ひょーが");
```

/docker-entrypoint-initdb.d以下のファイルはアルファベット順に読み込まれるようなので、
0_hoge.sql, 1_foo.shのようにしておくと楽かもしれません。


## mysql下準備まとめ

ここでは以下の三つのことをするための下準備をしました。

 * 文字コードの変更
 * logをローカルで見れるようにする
 * 初期データを作る

ここまでの操作は全てローカルで行われています。
ディレクトリ構成は以下のようになりました。
詳しくはgithubで確認して見てください.
[この記事のコード](https://github.com/g-hyoga/test-golang-on-docker)

```
.
|--mysql
   |--init
   |  |--0_init.sql
   |--log
   |--my.cnf

```


# 下準備:Golangコンテナ

特になんてことはしません.
ここでは以下のことをします.

 * jsonを返すだけのコードをかく
 * vendoringする

## jsonを返すだけのコードをかく

##### ./src/cmd/main.go
```go
package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int
	Name string
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		db, err := sql.Open("mysql", "root:password@tcp(db-server:3306)/test_db")
		defer db.Close()
		if err != nil {
			fmt.Println(err.Error())
		}

		rows, err := db.Query("SELECT * FROM test_table")
		defer rows.Close()
		if err != nil {
			fmt.Println(err)
		}

		user := User{}
		for rows.Next() {
			err = rows.Scan(&user.ID, &user.Name)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(user)
		}

		c.JSON(200, gin.H{
			"hello": user.Name,
		})
	})
	r.Run()
}
```

コードが汚い？ハードコーディング？
なんの事やら

## vendoringする

これもやり方は複数あると思います。
何がベストプラクティスかわかりません、教えていただけると幸いです。

ここでは、[dep](https://golang.github.io/dep/)というものを使いました。

```sh
cd src
dep init
```

とするとsrc以下にvendorディレクトリが作成され、
依存パッケージが入ってくれます。
詳しい使い方は調べて見てください。

### 余談

このvendoringはgolangで用いている
依存パッケージをコンテナ内で使うために行なっています。

ここでは、依存パッケージをVOLUMEでコンテナ内に追加する事で解決しています。
ですが、コンテナを立ち上げ時に依存パッケージを全部ダウンロードということもできます。
他にも方法はあるでしょう。

何が良いのかぶっちゃけわからないので、教えてください。

## Golang下準備まとめ

ここでは、以下のことをしました。

 * jsonを返すだけのコードをかく
 * vendoringする

ここまでのディレクトリ構成は以下のようになっています。

```
.
|--src
|  |--Gopkg.lock
|  |--Gopkg.toml
|  |--cmd
|  |  |--main.go
|  |--vendor
|     |--github.com
|        |--gin-contrib
|        |--.....
|        |--.....
|        |--.....
|
|--mysql
   |--init
   |  |--0_init.sql
   |--log
   |--my.cnf
```

# docker-compose.ymlをかく

ここからが本番です。
実際にmysqlのコンテナとgoのコンテナを立ち上げ、
動かします。

先にdocker-compose.ymlが結局こうなったというのを示します。

##### ./docker-compose.yml
```yaml
version: '3.3'

services:
  db-server:
    image: mysql
    ports:
      - "3306:3306"
    volumes:
      - "data:/var/lib/mysql"
      - "data:/var/log/mysql"
      - "./mysql:/etc/mysql/conf.d"
      - "./mysql/init:/docker-entrypoint-initdb.d"
      - "./mysql/log:/var/log/mysql"
    environment:
      mysql_root_password: password
      
  app:
    image: golang:1.8-onbuild
    ports:
      - "8080:8080"
    volumes:
      - ".src/vendor:/go/src/github.com"
      - "./src:/go/src/app"
    depends_on:
      - db-server
    command: go run cmd/main.go

volumes:
  data:
```

## data volume

Dockerではコンテナを起動し、停止すると基本的にコンテナ内のデータは消えます。
永続的なデータを扱いたい時はData用のVOLUMEを用意します。

```yaml
volumes:
  data:
```

これだけで用意してくれるようです.

今回はmysqlのデータを永続化したいです。

```yaml
services:
  db-server:
   volumes:
      - "data:/var/lib/mysql"
      - "data:/var/log/mysql"
```

このようにdb-serverのvolumesにmysqlのデータ保存先を追加します.
今回はlogも永続化しておきます。


## mysqlの下準備を適用する

前半に書いたmysqlの下準備を適用していきます。
dockerのVolumeを用いて、コンテナ内の設定ファイルをローカルのものと置き換えます。

```yaml
services:
  db-server:
   volumes:
      - "./mysql:/etc/mysql/conf.d" #設定ファイルを置き換える
      - "./mysql/init:/docker-entrypoint-initdb.d" #ローカルのinitとコンテナのdocker-entrypoint-initdb.dを繋げる
      - "./mysql/log:/var/log/mysql" #ローカルの./mysql/logとコンテナと/var/log/mysqlを繋げる
```

```
ローカルのパス:コンテナ内のパス
```
という風に書きます。

設定ファイルを置き換えたり、ローカルのディレクトリとコンテナのディレクトリをつなげたりできます。
これによってmy.cnfは置き換えられ、mysqlで出力されたlogはローカルの./mysql/logでも確認できるようになります。

あとは3306portを使用するので、あけておきます

## goの下準備を適用する

mysqlと同じです。
volumeを使ってコードと依存パッケージをコンテナ内と共有します。

```yaml
services:
  app:
    image: golang:1.8-onbuild
    volumes:
      - ".src/vendor:/go/src/github.com"
      - "./src:/go/src/app"
```

気をつけることはそんなにないですが、
goの公式imageでは
GOPATHが/goになっているくらいでしょうか

特に意味はないですが、
```sh
docker-compose exec app bash
```
とかで、直接コンテナ内でビルドできたりするのでbuild済みのものを使いました.
(一回もしませんでしたが)

あとはこちらも8080ポートを開けておきます

## コンテナ同士の接続

今回は, GoコンテナからMySQLコンテナに接続する必要があります。
link(非推奨)というものを使うか、networkというものを使う方法があります.
docker-composeのversion3を使うと、
services以下に定義したコンテナ同士はデフォルトで、
同じnetworkによしなに属してくれます。

名前は,

```yaml
services:
  app:  #ここ
    ....
  db-server:  #ここ
    ....
```

のものが使われます。

なので、片方のコンテナから(例えば、今だとappのコンテナから)

```sh
ping db-server
```

するとちゃんとpingが通ります.


## 動かす

一思いに,
```sh
docker-compose up
```
します。

docker-compose.ymlにdepends_onと書いていると思うのですが、
起動する順番が決まるだけで、mysqlが起動するまで待ってくれないようです。

少し待ってから,
[http://localhost:8080/ping](http://localhost:8080/ping)にアクセスすれば、
mysqlから取得したデータが入ったjsonが返ってきているはずです.

