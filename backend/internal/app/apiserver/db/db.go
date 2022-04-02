package db

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/asdine/storm"
)

type OffersType struct {
	Offers     map[string]Offer `json:"ignore"`
	DB         *storm.DB
	Production bool
	Host       string
	S3path     string
}

var Offers = OffersType{map[string]Offer{}, nil, false, "", ""}

const (
	prefixMainImage = "mainImage_"
	prefixImages    = "images_"
)

func GetOffersFrom1C() ([]Offer, error) {
	offers1c := []Offer(nil)
	c := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", "https://1c.restar26.site/sm1/hs/offers/", nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("API", "6O7EHDWS0Sk$yZ%i80p5")

	start := time.Now()
	res, err := c.Do(req)
	log.Println(time.Since(start))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		log.Println("Ошибка в ответе", res.Status)
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if err := json.Unmarshal(body, &offers1c); err != nil {
		log.Fatalln(err)
	}

	for i, v := range offers1c {
		offer, err := GetOffer(v.GUID)
		// Если заказ не найдет в базе данных, пропускаем цикл
		if err == ErrOfferNotFound {
			offers1c[i].Images = []string{}
			Offers.Offers[v.GUID] = v
			continue
		}
		offers1c[i].UpdateOfferWithData(offer)
	}

	return offers1c, nil
}

//func ImageAdd(guid, imgPath string) error {
//	offer, err := GetOffer(guid)
//	if err != nil {
//		return err
//	}
//
//	if offer.MainImage == "" {
//		offer.MainImage = imgPath
//	}
//	offer.Images = append(offer.Images, imgPath)
//	if err := Offers.DB.Save(&offer); err != nil {
//		log.Fatalln(err)
//	}
//	return nil
//}

//func GetImages(guid string) (images []string, err error) {
//	offer, err := GetOffer(guid)
//	if err != nil {
//		return nil, err
//	}
//	return offer.Images, nil
//}

//func ImageDelete(guid, link string) {
//	img, err := url.Parse(link)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	offer, err := GetOffer(guid)
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	oimg, err := url.Parse(offer.MainImage)
//	if err != nil {
//		log.Println(err)
//	}
//	if oimg.Path == img.Path {
//		offer.MainImage = ""
//	}
//
//	newimgs := []string(nil)
//	for _, v := range offer.Images {
//
//		log.Println(img.Path, v)
//		vimg, err := url.Parse(v)
//		if err != nil {
//			log.Println(err)
//		}
//		if vimg.Path != img.Path {
//			if offer.MainImage == "" {
//				offer.MainImage = vimg.Path
//			}
//			newimgs = append(newimgs, vimg.Path)
//		}
//	}
//	offer.Images = newimgs
//	Offers.Offers[guid] = offer
//	if err := Offers.DB.Save(&offer); err != nil {
//		log.Fatalln(err)
//	}
//
//	log.Println(offer)
//	log.Println(link)
//}

func (o *Offer) ImageSetMain() {

}

func (ofs *OffersType) SaveOffers() {

}

func (o *Offer) GetMainImage() {

}

func (o *Offer) DeleteOffer() {

}

func (o *Offer) SetPrice() {

}
