[MySQLのdocker hubページ](https://hub.docker.com/_/mysql/)
見れば乗っているのですが、

```sh
docker run --name some-mysql -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:tag --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
```

```
--character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
```

のように書けば文字コードを変更してくれるようです。
docker-compose.ymlでは以下のように書けば良いようです。

```yml
...

  command: mysqld --character-set-server=utf8 --collation-server=utf8_unicode_ci

...
```

