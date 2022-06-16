package qr

import (
	"fmt"
	"image"
	"image/draw"
	"net/http"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/rs/zerolog/log"
	"github.com/signintech/gopdf"
	"github.com/skip2/go-qrcode"
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

func (qr *HTTPDelivery) NewQRCode(w http.ResponseWriter, r *http.Request) {
	dpm := 1 * 8 // dots per mm
	widthMM := 58
	heightMM := 30
	border := 1

	qrimage, err := newQRImage("https://re-star.ru", heightMM*dpm, border*dpm)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "")

		return
	}

	img := image.NewRGBA(image.Rect(0, 0, widthMM*dpm, heightMM*dpm))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src) // draw white background
	draw.Draw(img, img.Bounds(), qrimage, image.Point{}, draw.Src)

	if err = qr.addLabel(img, "WTF АБВГ WTF ЯЫК WTF asd asd aasdf asdf23412313123 фывфы", (heightMM*dpm)+8, 10); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "")

		return
	}

	if err = qr.addLabel(img, "asdf23412313123 фывфы", (heightMM*dpm)+8, 8*5); err != nil {
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

func (qr *HTTPDelivery) addLabel(img *image.RGBA, label string, x, y int) error {
	c := freetype.NewContext()
	c.SetFont(qr.fnt)
	c.SetDPI(203.0)

	size := 12.0
	c.SetFontSize(size)

	c.SetDst(img)
	c.SetClip(img.Bounds())
	c.SetSrc(image.Black)
	c.SetHinting(font.HintingNone)

	pt := freetype.Pt(x, y+int(c.PointToFixed(size)>>6))

	if _, err := c.DrawString(label, pt); err != nil {
		return fmt.Errorf("cant draw label: %w", err)
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
