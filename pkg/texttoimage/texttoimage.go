package texttoimage

import (
	"image"
	"image/draw"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

// TextToImageOptions are options struct for TextToImage function.
// If max width is 0, will use 400 instead.
type TextToImageOptions struct {
	Font     *truetype.Font
	FontSize float64
	DPI      float64

	Padding  int
	MaxWidth int
}

// TextToImage converts text to a png image.
func TextToImage(msg string, options *TextToImageOptions) *image.RGBA {
	ctx := freetype.NewContext()
	ctx.SetFontSize(options.FontSize)
	ctx.SetFont(options.Font)
	ctx.SetDPI(options.DPI)
	ctx.SetSrc(image.White)

	// measure text size and split to lines.
	face := truetype.NewFace(options.Font, &truetype.Options{
		Size: options.FontSize,
		DPI:  options.DPI,
	})

	line := ""
	splittedText := []string{}
	x := options.Padding
	height := ctx.PointToFixed(options.FontSize).Ceil() + options.Padding*3
	for _, c := range msg {
		width, _ := face.GlyphAdvance(c)

		// new line
		if x+width.Round() > options.MaxWidth-options.Padding*2 {
			splittedText = append(splittedText, line)
			line = string(c)
			height += ctx.PointToFixed(options.FontSize).Ceil()
			x = width.Ceil() + options.Padding
			continue
		}
		// add rune to line
		line += string(c)
		x += width.Ceil()
	}
	// add last line
	splittedText = append(splittedText, line)
	height += ctx.PointToFixed(options.FontSize).Ceil()

	println("text size: ", options.MaxWidth, height)
	// create a empty image
	dist := image.NewRGBA(image.Rect(0, 0, options.MaxWidth, height))
	draw.Draw(dist, dist.Bounds(), image.Black, image.Point{}, draw.Over)
	// draw line
	ctx.SetDst(dist)
	ctx.SetClip(dist.Bounds())
	ctx.SetSrc(image.White)

	for i, line := range splittedText {
		ctx.DrawString(line, fixed.Point26_6{
			X: ctx.PointToFixed(float64(options.Padding)),
			Y: ctx.PointToFixed(float64(i+1)*options.FontSize + float64(options.Padding)),
		})
		println("drawing: ", line, i, i*int(options.FontSize))
	}

	return dist
}
