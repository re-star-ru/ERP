package db

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/asdine/storm"
)

var (
	ErrOfferNotFound = errors.New("заказ не найден в базе данных и памяти")
)

type Offer struct {
	GUID        string   `json:"guid" storm:"id"`
	Name        string   `json:"name"`
	Spec        string   `json:"spec"`
	Amount      int      `json:"amount"`
	SKU         string   `json:"sku"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	MainImage   string   `json:"mainImage"` // path in s3
	Images      []string `json:"images"`    // path to all images in offer
}

func (o *Offer) Save() {
	if err := Offers.DB.Save(o); err != nil {
		log.Fatalln(err)
	}
}

func (o *Offer) ImagesGet() (img []string, err error) {
	return
}

func (o *Offer) ImageAdd(imgPath string) error {
	log.Println(o.MainImage)
	if o.MainImage == "" {
		o.MainImage = imgPath
	}
	o.Images = append(o.Images, imgPath)

	o.Save()
	return nil
}

func (o *Offer) ImageDelete(link string) error {
	log.Println(o.MainImage)
	imgToDel, err := url.Parse(link)
	if err != nil {
		log.Fatalln(err)
	}
	if o.MainImage == imgToDel.Path {
		o.MainImage = ""
	}

	newimgs := []string(nil)

	for _, v := range o.Images {
		log.Println(imgToDel.Path, v)
		if v != imgToDel.Path {
			if o.MainImage == "" {
				o.MainImage = v
			}
			newimgs = append(newimgs, v)
		}
	}
	o.Images = newimgs

	o.Save()
	return nil
}

func GetOffer(guid string) (offer Offer, err error) {
	if guid == "" {
		return offer, errors.New("пустой guid")
	}
	err = Offers.DB.One("GUID", guid, &offer)
	if err != nil && err != storm.ErrNotFound {
		log.Fatalln(err)
		return
	}
	if err == storm.ErrNotFound {
		if _, ok := Offers.Offers[guid]; !ok {
			return offer, ErrOfferNotFound
		}
		offer = Offers.Offers[guid]
		return offer, nil
	}
	return offer, nil
}

func (o *Offer) UpdateOfferWithData(data Offer) {
	if data.MainImage != "" {
		o.MainImage = fmt.Sprintf("%v%v", Offers.S3path, data.MainImage)
	}
	if len(data.Images) == 0 {
		o.Images = []string{}
		return
	}
	for i, _ := range data.Images {
		o.Images = append(o.Images, fmt.Sprintf("%v%v", Offers.S3path, data.Images[i]))
	}
}
