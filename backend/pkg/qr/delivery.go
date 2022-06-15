package qr

import (
	"backend/pkg"
	"bytes"
	"fmt"
	pdf "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/skip2/go-qrcode"
	"image"
	"image/draw"
	"net/http"
)

type HTTPDelivery struct {
}

func NewHTTPDelivery() *HTTPDelivery {
	return &HTTPDelivery{}
}

func (qr *HTTPDelivery) NewQRCode(w http.ResponseWriter, r *http.Request) {
	dpm := 1 * 8 // dots per mm
	widthMM := 58 * dpm
	heightMM := 30 * dpm
	border := 1 * dpm

	qrimage, err := newQRImage("https://www.unidoc.com", heightMM, border)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	widthMM, heightMM = heightMM, widthMM // swap width and height

	img := image.NewRGBA(image.Rect(0, 0, widthMM, heightMM))
	draw.Draw(img, img.Bounds(), qrimage, image.Point{}, draw.Src)

	newPdf := new(bytes.Buffer)

	for i := 3; i > 0; i-- {

		//c.NewPage()
		//
		//if err = c.RotateDeg(-90); err != nil {
		//	pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "")
		//
		//	return
		//}
		//
		//pimg.SetPos(0, 0)
		//
		//if err = c.Draw(pimg); err != nil {
		//	pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "")
		//
		//	return
		//}
	}

	if err = pdf.ImportImages(nil, newPdf, nil, nil, nil); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "")

		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	if _, err = newPdf.WriteTo(w); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "")

		return
	}
	//
	//png.Encode(buf)
	////png.Encode(w, qrimage)
	//
	//png
	//
	//pdf.ImportImages(nil, newPdf, buf, nil)
	//
	//img.
	//	pdf.ImportImages(buf, wr)
	//
	//pdf.
	//	pdf.Usage()
	//
	//c := creator.New()
	//
	//pimg, err := c.NewImageFromGoImage(img)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//
	//	return
	//}
	//
	//c.SetPageSize(creator.PageSize{float64(widthMM), float64(heightMM)})
	//
	//optimize.Op
	//c.SetOptimizer(model.Optimizer(model.OptimizeTypeNone))
	//
	//for i := 3; i > 0; i-- {
	//	c.NewPage()
	//
	//	if err = c.RotateDeg(-90); err != nil {
	//		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "")
	//
	//		return
	//	}
	//
	//	pimg.SetPos(0, 0)
	//
	//	if err = c.Draw(pimg); err != nil {
	//		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "")
	//
	//		return
	//	}
	//}
	//
	////w.Header().Set("Content-Type", "image/png")
	//
	//err = c.Write(w)
	//if err != nil {
	//	pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "")
	//
	//	return
	//}
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
