package main

import (
  "fmt"
  "net/http"
  "io"
  "log"
)

var client = http.Client{}

func imgproxy (w http.ResponseWriter, r *http.Request) {
  // r.URL.Pathは「/」で始まるから、「https://」じゃなくて、「https:/」となります。
  img, err := client.Get("https:/" + r.URL.Path)
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
