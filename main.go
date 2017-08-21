package main

//go:generate go-bindata assets/

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"

	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	sb "github.com/yzguy/sPoNgEbOb"
	// "golang.org/x/image/font/gofont/goregular"
	goregular "golang.org/x/image/font/gofont/gobold"
)

func main() {
	testArgs := os.Args[1:]

	text := sb.Mock(strings.Join(testArgs, " "))
	os.Stderr.Write([]byte(text))

	bimgb, _ := assetsSpongebobiconJpgBytes()
	buf := bytes.NewBuffer(bimgb)
	baseImg, _, _ := image.Decode(buf)

	bb := baseImg.Bounds()

	bufImg := image.NewRGBA(bb)

	draw.Draw(bufImg, bb, baseImg, bb.Min, draw.Over)

	f, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatalln(err)
	}

	textImg := image.NewRGBA(bb)
	ctx := freetype.NewContext()
	ctx.SetClip(bb)
	ctx.SetDst(textImg)
	ctx.SetSrc(image.White)
	ctx.SetFont(f)

	p, s := getDrawParams(text, bb, f, ctx.PointToFixed(48))

	ctx.SetFontSize(s)
	_, err = ctx.DrawString(text, p)
	if err != nil {
		log.Fatalln(err)
	}

	draw.Draw(bufImg, bb, getStroke(textImg), bb.Min, draw.Over)
	ctx.SetDst(bufImg)
	_, err = ctx.DrawString(text, p)
	if err != nil {
		log.Fatalln(err)
	}

	png.Encode(os.Stdout, bufImg)
}

func getDrawParams(text string, b image.Rectangle, f *truetype.Font, scale fixed.Int26_6) (p fixed.Point26_6, size float64) {
	rs := []rune(text)

	size = 48

	mid := b.Max.Div(2)

	textWidth := 0

	for _, r := range rs {
		idx := f.Index(r)
		textWidth += f.HMetric(scale, idx).AdvanceWidth.Round()
	}

	mpt := image.Pt(textWidth, 0).Div(2)
	ipt := mid.Sub(mpt)

	return fixed.P(ipt.X, 100), size
}

func getStroke(i image.Image) image.Image {
	bounds := i.Bounds()
	out := image.NewRGBA(bounds)

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			strPix := getNearby(i, image.Pt(x, y), 1)
			if strPix > 0 {
				out.Set(x, y, color.RGBA{
					R: 0,
					G: 0,
					B: 0,
					A: strPix,
				})
			}
		}
	}

	return out
}

func getNearby(i image.Image, p image.Point, prox int) uint8 {
	for dy := p.Y - prox; dy < p.Y+prox+1; dy++ {
		for dx := p.X - prox; dx < p.X+prox+1; dx++ {

			if p.X == dx && p.Y == dx {
				continue
			}

			_, _, _, a := i.At(dx, dy).RGBA()
			if a > 0 {
				return 255
			}
		}
	}

	return 0
}
