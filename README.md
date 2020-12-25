# GORM & go-sqlmock Example

## 書いてある処理
- 全取得
- 一部取得
- 作成
- 削除
- 更新

## sqlite3 mock
```
記述なし
```

## mysql mock
> mysqlの場合は `SkipInitializeWithVersion: true`が必須です。

```go
func GetNewDbMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, mock, err
	}

	gormDB, err := gorm.Open(mysql.Dialector{Config: &mysql.Config{DriverName: "mysql", Conn: db, SkipInitializeWithVersion: true}}, &gorm.Config{})

	if err != nil {
		return gormDB, mock, err
	}

	return gormDB, mock, err
}
```

## postgresql mock
```
記述なし
```

## 今までの問題をクリア、SELECT, INSERTの問題解決
注意: SQLクエリの正規表現は重要


## 参考記事
[Goでデータベースを簡単にモック化する【sqlmock】](https://qiita.com/gold-kou/items/cb174690397f651e2d7f)

[gorm Transactions official guide](https://gorm.io/ja_JP/docs/transactions.html)

[DATA-DOG/go-sqlmockを使ってGormDBをmockする](https://tech.fusic.co.jp/posts/2020-12-02-mock-gormdb-using-go-sqlmock/)

[go-sqlmockを使ってGORMで書いたコードをテストする](https://qiita.com/otanu/items/761de2bfc38468e9d353)

[Gorm2.0でsqlMockに対応する](https://qiita.com/hosakak/items/a20af188846ef48f2e03)

[gorm v2.0 unit testing with sqlmock #3565](https://github.com/go-gorm/gorm/issues/3565)

[testifyによるMockテストの参考例](https://tutorialedge.net/golang/improving-your-tests-with-testify-go/)

[testifyのMockメソッドに関する詳細説明](https://qiita.com/muroon/items/f8beec802c29e66d1918#mockon-%E3%83%A1%E3%82%BD%E3%83%83%E3%83%89)

[testifyによるMockテストのControllerとModelの例に役立った](https://qiita.com/takeshi_miyajim/items/d2fe1ed3c2e85b014b02)