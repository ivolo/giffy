
package main

import (
  "fmt"
  //"log"
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
    query := "simpsons ralph"

    c := giphy.New("dc6zaTOxFJmzC")
    gifs, err := c.Search(query)
    check(err)

    fmt.Printf("Found %d gifs for '%s'.\n", len(gifs), query)

    if len(gifs) == 0 {
      return
    }

    for i := 3; i < 6; i += 1 {
      g := gifs[i]
      gif, err := download(g.Images["original"].URL)
      check(err)

      ttyWidth, ttyHeight, err := terminal.GetSize(1)
      check(err)

      // https://github.com/dpup/go-scratch/blob/master/gif-resize/gif-resize.go#L3
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

      // Resize each frame.
      for index, frame := range gif.Image {
        bounds := frame.Bounds()
        draw.Draw(img, bounds, frame, bounds.Min, draw.Over)
        resized := resize.Resize(uint(ttyWidth), uint(ttyHeight), img, resize.NearestNeighbor)
        b2 := resized.Bounds()
        pm := image.NewPaletted(b2, palette.Plan9)
        draw.FloydSteinberg.Draw(pm, b2, resized, image.ZP)
        gif.Image[index] = resized
      }

      for _, img := range(gif.Image) {
        //size := img.Bounds().Max;
        //log.Printf("Parsed png image [x: %d y: %d]", size.X, size.Y)

      
        //log.Printf("Terminal size [x: %d y: %d]", ttyWidth, ttyHeight)

        //resized := resize.Resize(uint(ttyWidth), uint(ttyHeight), img, resize.NearestNeighbor)

        str := ascii.Convert(img)

        fmt.Print(str)
        fmt.Println(g.Images["original"].URL)
        x := 0
        y := 0
        r, g, b, a := img.At(x, y).RGBA()
        r2, g2, b2, a2 := uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)
        fmt.Println(img.At(x,y))
        fmt.Println(gif.Config.ColorModel.Convert(img.At(x,y)))
        fmt.Println(r, g, b, a)
        fmt.Println(r2, g2, b2, a2)
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