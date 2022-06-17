package qr

import (
	"encoding/json"
	"fmt"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"net/http"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/rs/zerolog/log"
	"github.com/signintech/gopdf"
	"golang.org/x/image/font"

	"backend/pkg"
)

type HTTPDelivery struct {
	fnt *truetype.Font
}

func NewHTTPDelivery() *HTTPDelivery {
	fontBytes, err := os.ReadFile("./assets/fonts/RobotoMono-Medium.ttf")
	if err != nil {
		log.Fatal().Err(err).Msg("cant read font")
	}

	fnt, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal().Err(err).Msg("cant parse font")
	}

	return &HTTPDelivery{
		fnt: fnt,
	}
}

type qrrequest struct {
	W    int      `json:"width"`  // > 0, in mm
	H    int      `json:"height"` // > 0, in mm
	Code string   `json:"code"`
	Text []string `json:"text"`
}

func (qr *HTTPDelivery) NewQRCode(w http.ResponseWriter, r *http.Request) {
	const dpm = 1 * 8 // dots per mm

	var req qrrequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "")

		return
	}

	if req.W <= 0 || req.H <= 0 || len(req.Code) < 1 {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("invalid width or height"), "")

		return
	}

	// widthMM := 58
	// heightMM := 30
	// widthMM := 43
	// heightMM := 25

	widthMM := req.W
	heightMM := req.H
	code := req.Code
	labels := req.Text

	border := 1

	qrimage, err := newQRImage(code, heightMM*dpm, border*dpm)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "")

		return
	}

	img := image.NewRGBA(image.Rect(0, 0, widthMM*dpm, heightMM*dpm))

	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src) // draw white background
	draw.Draw(img, img.Bounds(), qrimage, image.Point{}, draw.Src)

	labeler := newLabel(heightMM*dpm, 0, qr.fnt, img)

	if err = labeler.print(labels); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "")

		return
	}

	rect := gopdf.Rect{W: float64(widthMM), H: float64(heightMM)}
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{Unit: gopdf.UnitMM, PageSize: *gopdf.PageSizeA4})

	pdf.AddPage()

	if err = pdf.ImageFrom(img, 0, 0, &rect); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "")

		return
	}

	w.Header().Set("Content-Type", "application/pdf")

	if err = pdf.Write(w); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "")

		return
	}
}

type label struct {
	prefix, middle, suffix string

	c *freetype.Context

	maxChars int
	pt       fixed.Point26_6
	line     fixed.Int26_6
}

func newLabel(x, y int, fnt *truetype.Font, img *image.RGBA) *label {
	const fontSize = 10.0

	lbl := &label{}

	ctx := freetype.NewContext()
	ctx.SetFont(fnt)
	ctx.SetDPI(203.0) //default dpi
	ctx.SetFontSize(fontSize)
	ctx.SetHinting(font.HintingNone)
	ctx.SetSrc(image.Black)

	lbl.pt = freetype.Pt(x, y+int(ctx.PointToFixed(fontSize)>>6))
	lbl.c = ctx
	lbl.line = ctx.PointToFixed(fontSize * 1.5)

	// calc max charaters in line TODO
	// work only on monospaced fonts

	lbl.maxChars = lbl.calcMaxChars(img.Bounds().Max.X - x)
	log.Printf("max chars: %v", lbl.maxChars)

	ctx.SetDst(img)
	ctx.SetClip(img.Bounds())

	return lbl
}

func (lbl *label) calcMaxChars(lineWidth int) int {
	log.Printf("line width: %v", lineWidth)

	lbl.c.SetDst(image.NewRGBA(image.Rect(0, 0, lineWidth, 100)))

	var (
		left     = lbl.pt
		maxChars int
	)

	for {
		newPt, err := lbl.c.DrawString("X", left)
		if err != nil {
			log.Fatal().Err(err).Msg("cant draw label")
		}

		log.Printf("%+v, %+v", left, newPt)

		if (newPt.X - lbl.pt.X).Round() > lineWidth {
			break
		}

		left = newPt

		maxChars++
	}

	log.Printf("%v", maxChars)

	return maxChars
}

func (lbl *label) print(texts []string) error {
	log.Printf("%+v", lbl.pt)

	for _, text := range texts {
		_, err := lbl.c.DrawString(text, lbl.pt)
		if err != nil {
			return fmt.Errorf("cant draw label: %w", err)
		}

		lbl.pt.Y += lbl.line
	}

	return nil
}

func newQRImage(data string, width, border int) (image.Image, error) {
	qrImage, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("cant encode qrcode: %w", err)
	}

	qrImage.DisableBorder = true

	img := image.NewRGBA(image.Rect(0, 0, width, width))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src) // draw white background

	draw.Draw(
		img, image.Rect(border, border, width+border, width+border),
		qrImage.Image(width-(border*2)), image.Point{}, draw.Src)

	return img, nil
}
