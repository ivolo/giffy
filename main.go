
package main

import (
  "fmt"
  //"log"
  "image/gif"
  "golang.org/x/crypto/ssh/terminal"
  "github.com/nfnt/resize"
  "github.com/ivolo/go-image-to-ascii"
  "github.com/ivolo/go-giphy"
  "errors"
  "net/http"
)

func check(err error) {
  if err != nil {
    panic(err)
  }
}

func main() {
    query := "simpsons ralph"

    c := giphy.New("dc6zaTOxFJmzC")
    gifs, err := c.Search(query)
    check(err)

    fmt.Printf("Found %d gifs for '%s'.\n", len(gifs), query)

    if len(gifs) == 0 {
      return
    }

    for _, g := range(gifs) {
      gif, err := download(g.Images["original"].URL)
      check(err)

      for _, img := range(gif.Image) {
        //size := img.Bounds().Max;
        //log.Printf("Parsed png image [x: %d y: %d]", size.X, size.Y)

        ttyWidth, ttyHeight, err := terminal.GetSize(1)
        check(err)
        //log.Printf("Terminal size [x: %d y: %d]", ttyWidth, ttyHeight)

        resized := resize.Resize(uint(ttyWidth), uint(ttyHeight), img, resize.Lanczos3)

        str := ascii.Convert(resized)

        fmt.Print(str)
      }
    }
}

func download(url string) (*gif.GIF, error) {
  //log.Printf("GET %s ..", url)
  res, err := http.Get(url)
  if err != nil {
    return nil, err
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
    return nil, errors.New(fmt.Sprintf("error response '%d'", res.StatusCode))
  }
  gif, err :=  gif.DecodeAll(res.Body)
  return gif, err
}
