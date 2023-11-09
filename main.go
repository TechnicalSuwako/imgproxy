package main

import (
  "fmt"
  "io"
  "log"
  "net/http"
  "strings"
  "time"
)

var sofname = "imgproxy"
var version = "1.1.0"
var client = http.Client{
  Timeout: time.Second * 30,
}

func fetchImage(curl string, nurl string) (*http.Response, error) {
  req, err := http.NewRequest("GET", curl, nil)
  if err != nil {
    return nil, err
  }

  resp, err := client.Do(req)
  if err != nil {
    return nil, err
  }

  // 画像だけをプロクシーする様に
  ctype := resp.Header.Get("Content-Type")
  if ctype != "image/jpeg" && ctype != "image/png" && ctype != "image/gif" && ctype != "image/tga" && ctype != "image/targa" && ctype != "image/webp" {
    return nil, fmt.Errorf("画像ではありません: %s", ctype)
  }

  return resp, nil
}

func serveImage(w http.ResponseWriter, img *http.Response) {
  w.Header().Set("Content-Type", img.Header.Get("Content-Type"))
  w.Header().Set("Content-Length", fmt.Sprint(img.ContentLength))
  if _, err := io.Copy(w, img.Body); err != nil {
    log.Printf("画像を出力に失敗: %v", err)
  }
}

func pixivImg(w http.ResponseWriter, req *http.Request) {
  resp, err := client.Do(req)
  if err != nil {
    log.Printf("Pixivから画像を受取に失敗: %v", err)
    http.Error(w, "Pixivから画像を受取に失敗。", http.StatusInternalServerError)
    return
  }
  defer resp.Body.Close()
  serveImage(w, resp)
}

func imgproxy(w http.ResponseWriter, r *http.Request) {
  // URLを修正
  nurl := r.URL.Path[1:]
  if nurl == "" {
    fmt.Fprintf(
      w,
      `<style>body { background: #1b0326; color: #f545f5; } a { color: #ff88ff; }</style><p style="text-align: center;"><a href="https://technicalsuwako.moe/blog/%s-%s">%s-%s</a> | <a href="https://gitler.moe/suwako/%s">Git</a> | <a href="https://076.moe/">０７６</a></p>`,
      sofname,
      strings.ReplaceAll(version, ".", ""),
      sofname,
      version,
      sofname,
    )
    return
  }

  // Pixivの場合
  if strings.Contains(nurl, "s.pximg.net") || strings.Contains(nurl, "i.pximg.net") {
    req, _ := http.NewRequest("GET", "https://" + nurl, nil)
    req.Header.Set("Referer", "https://www.pixiv.net/")
    pixivImg(w, req)
    return
  }

  var img *http.Response
  var err error

  // HTTPリク
  img, err = fetchImage("https://" + nurl, nurl)
  if err != nil {
    img, err = fetchImage("http://" + nurl, nurl)
  }

  if err != nil {
    log.Printf("画像を受け取りに失敗: %v", err)
    http.Error(w, "画像を受け取りに失敗。", http.StatusInternalServerError)
    return
  }

  defer img.Body.Close()
  serveImage(w, img)
}

func main() {
  http.HandleFunc("/", imgproxy)
  fmt.Println("http://0.0.0.0:9810 でサーバーを実行中。終了するには、CTRL+Cを押して下さい。")
  http.ListenAndServe("0.0.0.0:9810", nil)
}
