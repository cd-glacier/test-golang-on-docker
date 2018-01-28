今更ながらDockerの勉強をしています。
色々わからないところが多く、詰まったりしたところが多かったので自分はこのようにしたというのを書きます。
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

 * localhost:8080/pingにアクセスするとjsonを返すだけのコードをかく
 * vendoringする

## localhost:8080/pingにアクセスするとjsonを返すだけのコードをかく

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



