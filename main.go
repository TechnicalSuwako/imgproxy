package main

import (
  "fmt"
  "net/http"
  "io"
  "log"
  "net/url"
)

var version = "1.0.0"
var client = http.Client{}

func imgproxy (w http.ResponseWriter, r *http.Request) {
  // URLを確認
  uri, err := url.Parse("https:/" + r.URL.Path)
  if err != nil {
    fmt.Println(err)
    return
  }

  // HTTPリク
  req, err := http.NewRequest("GET", "https:/" + r.URL.Path, nil)
  if err != nil {
    fmt.Println(err)
    return
  }

  // Pixivかどうかの確認
  if uri.Host == "i.pximg.net" || uri.Host == "s.pximg.net" {
    req.Header.Set("Referer", "https://www.pixiv.net/")
  }

  // r.URL.Pathは「/」で始まるから、「https://」じゃなくて、「https:/」となります。
  img, err := client.Do(req)
  if err != nil {
    fmt.Fprintf(w, "Error %d", err)
    return
  }

  // 画像だけをプロクシーする様に
  ctype := img.Header.Get("Content-Type")
  if ctype != "image/jpeg" && ctype != "image/png" && ctype != "image/gif" && ctype != "image/tga" && ctype != "image/targa" && ctype != "image/webp" {
    return
  }

  // ヘッダー
  w.Header().Set("Content-Length", fmt.Sprint(img.ContentLength))
  w.Header().Set("Content-Type", img.Header.Get("Content-Type"))

  // 表示
  if _, err = io.Copy(w, img.Body); err != nil {
    log.Fatalf("ioエラー：%v", err)
  }

  // もう要らない
  img.Body.Close()
}

func main () {
  http.HandleFunc("/", imgproxy)
  http.ListenAndServe(":9810", nil)
}
