# delete-esa-posts

esaの記事を一括削除するスクリプト

## インストール

```bash
# depで依存パッケージを取得するため、`go get`は使わない
$ git clone git@github.com:animateLAB/delete-esa-posts $GOPATH/src/github.com/animateLAB/delete-esa-posts
$ cd $GOPATH/src/github.com/animateLAB/delete-esa-posts
$ dep ensure -update # 依存パッケージの取得
$ go build
```

## 使い方

```bash
# `go build`でコンパイルした場合
$ ESA_TOKEN=esaのtoken ESA_TEAM=team名 ESA_SEARCH_QUERY="in:Archived/削除予定" ./delete-esa-posts
# コンパイルしない場合
$ ESA_TOKEN=esaのtoken ESA_TEAM=team名 ESA_SEARCH_QUERY="in:Archived/削除予定" go run main.go
```

- `ESA_TOKEN`
    - esaのAPIを用いるためのtoken
- `ESA_TEAM`
    - 削除する記事があるteam名
- `ESA_SEARCH_QUERY`
    - 検索クエリ
    - この検索クエリにマッチする記事が削除される
