package main

//go:generate go-bindata assets/

import (
	"bytes"
	"image"
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

	outImg := image.NewRGBA(baseImg.Bounds())

	draw.Draw(outImg, baseImg.Bounds(), baseImg, baseImg.Bounds().Min, draw.Over)

	f, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := freetype.NewContext()
	ctx.SetClip(baseImg.Bounds())
	ctx.SetDst(outImg)
	ctx.SetSrc(image.White)
	ctx.SetFont(f)

	p, s := getDrawParams(text, baseImg.Bounds(), f, ctx.PointToFixed(48))

	ctx.SetFontSize(s)
	_, err = ctx.DrawString(text, p)
	if err != nil {
		log.Fatalln(err)
	}

	png.Encode(os.Stdout, outImg)
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

	return fixed.P(ipt.X, b.Max.Y-70), size
}
