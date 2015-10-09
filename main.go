
package main

import (
  "fmt"
  "strings"
  "os"
  "image"
  "image/gif"
  "image/draw"
  "golang.org/x/crypto/ssh/terminal"
  "github.com/nfnt/resize"
  "github.com/ivolo/go-image-to-ascii"
  "image/color/palette"
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
    query := strings.Join(os.Args[1:], " ")
    if len(query) == 0 {
      fmt.Println("usage: giffy <query>")
      os.Exit(1)
    }

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

      ttyWidth, ttyHeight, err := terminal.GetSize(1)
      check(err)

      // fix inconsistent frame sizing with dealising
      dealias(gif, uint(ttyWidth), uint(ttyHeight))

      for _, img := range(gif.Image) {
        resized := resize.Resize(uint(ttyWidth), uint(ttyHeight), img, resize.NearestNeighbor)
        str := ascii.Convert(resized)
        fmt.Print(str)
        
        // fmt.Printf("\x1b[%dA", ttyHeight) // move cursor up
        // fmt.Printf("\x1b[%dD", ttyWidth) // move cursor left
        // fmt.Printf("\x1b[%dF", ttyHeight) // move cursor to the beginning of the line
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
  return gif.DecodeAll(res.Body)
}

func dealias(gif *gif.GIF, width uint, height uint) {
  // TODO: add better dealiasing algorithm: http://stackoverflow.com/questions/9988517/resize-gif-animation-pil-imagemagick-python
  
  // credit: https://github.com/dpup/go-scratch/blob/master/gif-resize/gif-resize.go#L3
  // This demonstrates a solution to resizing animated gifs.
  //
  // Frames in an animated gif aren't necessarily the same size, subsequent
  // frames are overlayed on previous frames. Therefore, resizing the frames
  // individually may cause problems due to aliasing of transparent pixels. This
  // example tries to avoid this by building frames from all previous frames and
  // resizing the frames as RGB.

  // Create a new RGBA image to hold the incremental frames.
  firstFrame := gif.Image[0].Bounds()
  b := image.Rect(0, 0, firstFrame.Dx(), firstFrame.Dy())
  img := image.NewRGBA(b)

  for index, frame := range gif.Image {
    bounds := frame.Bounds()
    draw.Draw(img, bounds, frame, bounds.Min, draw.Over)
    resized := resize.Resize(width, height, img, resize.NearestNeighbor)
    b2 := resized.Bounds()
    pm := image.NewPaletted(b2, palette.Plan9)
    draw.FloydSteinberg.Draw(pm, b2, resized, image.ZP)
    gif.Image[index] = pm
  }
}