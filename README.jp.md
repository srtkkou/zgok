# Zgok

[![GoDoc](https://godoc.org/github.com/srtkkou/zgok?status.svg)](https://godoc.org/github.com/srtkkou/zgok) [![Build Status](https://travis-ci.org/srtkkou/zgok.svg?branch=master)](https://travis-ci.org/srtkkou/zgok)

Zgokは静的ファイルを実行可能なバイナリに埋め込む
ための[Go](http://golang.org/)製のツールです。

English: README.md

## 特徴

* 全ての静的ファイルを埋め込んだ実行可能なバイナリを作成することが出来ます。デプロイが簡単です。
* 他のライブラリに依存しません。Windows/Linux/Mac/ARMで実行可能です。

## インストール方法

以下の手順でライブラリとコマンドがインストール出来ます。

	go get -u github.com/srtkkou/zgok/...

## 使い方

以下のコマンドでzgokの実行可能バイナリを作成できます。

	$GOPATH/bin/zgok build -e exePath -z zipPath1 -z zipPath2 -o outPath

Goのプログラム内で埋め込んだバイナリを読みたい場合、以下のようなコードで
読むことが出来ます。

```go
package main

import (
	"fmt"
	"github.com/srtkkou/zgok"
	"io/ioutil"
	"os"
)

func main() {
	var content []byte
	path := "test.txt"
	zfs, _ := zgok.RestoreFileSystem(os.Args[0])
	if zfs != nil {
		// For release.
		content, _ = zfs.ReadFile(path)
	} else {
		// For development.
		content, _ = ioutil.ReadFile(path)
	}
	fmt.Println(string(content))
}
```

Webアプリケーションにてzgokで埋め込んだ静的ファイルを
以下のようなコードで公開することが出来ます。
注意: 静的ファイルが [./web/public/*] に保存されていることを前提にしています。

1. 以下がWebアプリケーションでzgokを使う際の最小のコード例になります。

```go
package main

import (
	"net/http"
	"github.com/srtkkou/zgok"
	"os"
)

func main() {
	zfs, err := zgok.RestoreFileSystem(os.Args[0])
	if err != nil {
		panic(err)
	}
	assetServer := zfs.FileServer("web/public")
	http.Handle("/assets/", http.StripPrefix("/assets/", assetServer))
	http.ListenAndServe(":8080", nil)
}
```

2. zgok実行可能バイナリのビルド

	go build -o web web.go
	$GOPATH/bin/zgok build -e web -z web/public -o web_all

3. ブラウザにて [http://localhost:8080/assets/css/sample.css] のようなURLを参照する。

## 説明

Linuxの実行可能ファイル形式(ELF)とWindowsの実行可能ファイル形式(PE)は
実行可能ファイル内の各セグメントのサイズを格納したヘッダーセクションから
始まります。つまり実行可能ファイル形式は実際のファイルのサイズをヘッダーの
情報として持っています。なので、実行可能ファイル形式の末尾に追加のデータを
追加しても、そのまま実行可能ファイルは通常通り実行出来ます。

Zgokで作られた実行可能ファイルは以下のような形式です。

| ヘッダー　　　 |
| -------------- |
| セクション1　  |
| セクション2　  |
| ...            |
| セクションn　  |
| ZIP圧縮データ  |
| Zgok情報　　　 |

ZgokはZIP圧縮データを展開し、そのパスをキーとしたmapでファイル内容に
アクセス可能にします。

## ライセンス

Apache License Version 2.0. 詳細は LICENSE ファイルを参照して下さい。

