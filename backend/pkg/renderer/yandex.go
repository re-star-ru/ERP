package renderer

import (
	"backend/pkg/item"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

// Заполните важное поле – «Описание товара»
// Заполните важное поле – «Штрихкод»

type YandexOffer struct {
	Type       string `xml:"type,attr"` //*
	Model      string `xml:"model"`     //*
	Vendor     string `xml:"vendor"`    //*
	TypePrefix string `xml:"typePrefix"`

	ID         string `xml:"id,attr"` //*
	VendorCode string `xml:"vendorCode"`
	Url        string `xml:"url"`        //*
	Price      int    `xml:"price"`      //*
	CurrencyId string `xml:"currencyId"` //*
	CategoryId string `xml:"categoryId"` //*
	Count      int    `xml:"count"`      //*

	Store       bool   `xml:"store"`
	Weight      string `xml:"weight"`      //* kilo
	Dimensions  string `xml:"dimensions"`  //* sm (len/width/height)
	Description string `xml:"description"` //*
	Barсode     string `xml:"barcode"`     //*

	Condition *Condition `xml:"condition,omitempty"` // Б/у или новый
	Pictures  []string   `xml:"picture"`             //*
	Params    []Param    `xml:"param"`
}

type Param struct {
	Name  string `xml:"name,attr"`
	Unit  string `xml:"unit,attr,omitempty"`
	Value string `xml:",chardata"`
}

type Condition struct {
	Type   string `xml:"type,attr"`
	Reason string `xml:"reason"`
}

type Category struct {
	ID       string `xml:"id,attr"`
	ParentID string `xml:"parentId,attr"`
	Value    string `xml:",chardata"`
}

type Shop struct {
	Name    string `xml:"name"`
	Company string `xml:"company"`
	Url     string `xml:"url"`
	Email   string `xml:"email"`

	Categories []*Category    `xml:"categories>category"`
	Offers     []*YandexOffer `xml:"offers>offer"`
}

type YandexCatalog struct {
	XMLName xml.Name `xml:"yml_catalog"`
	Date    string   `xml:"date,attr"`

	Shop *Shop `xml:"shop"`
}

func YandexRender(items []item.Item) (io.Reader, string, error) {
	offers := make([]*YandexOffer, len(items))
	for idx := 0; len(items) > idx; idx++ {
		if items[idx].Price == 0 {
			continue
		}

		if items[idx].Char == "RG" { // пропускать бу товары
			continue
		}

		offers[idx] = &YandexOffer{
			ID:         items[idx].ID,
			Type:       "vendor.model",
			TypePrefix: items[idx].Type,
			Vendor:     items[idx].Char,
			Model:      items[idx].SKU,
			VendorCode: items[idx].Name,
			CurrencyId: "RUR",
			CategoryId: "100",
			Price:      items[idx].Price / 100,
			Count:      items[idx].Amount,

			Store:       true,
			Weight:      "5.000",
			Dimensions:  "40.000/20.000/20.000",
			Description: "описание товара",
			Barсode:     "9876543210",

			Params: []Param{
				{Name: "Марка автомобиля", Value: "Audi"},
			},
		}

		if items[idx].Char == "RG" {
			offers[idx].Condition = &Condition{
				Type:   "used",
				Reason: "Следы потертостей на корпусе, проведена предпродажная подготовка", // todo вынести это в 1с
			}
		}

		imgs := make([]string, 0, len(items[idx].Images))
		for _, v := range items[idx].Images {
			imgs = append(imgs, v.Path)
		}
		if len(imgs) == 0 {
			imgs = append(imgs, "https://s3.re-star.ru/srv1c/images/b6cc4653-8694-4b13-83d5-5747ef6507c4.jpeg")
		}

		offers[idx].Pictures = imgs
	}

	y := YandexCatalog{
		Date: time.Now().Format(time.RFC3339),
		Shop: &Shop{
			Name:    "RESTAR",
			Company: "Стартеры Генераторы Пятигорск",
			Url:     "https://re-star.ru",
			Email:   "re-star@mail.ru",
			Categories: []*Category{
				{ID: "1", Value: "Авто"},
				{ID: "10", ParentID: "1", Value: "Запчасти"},
				{ID: "100", ParentID: "10", Value: "Электрика"},
			},

			Offers: offers,
		},
	}

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(xml.Header)

	enc := xml.NewEncoder(buf)
	enc.Indent("", "  ")

	if err := enc.Encode(y); err != nil {
		return nil, "", fmt.Errorf("cant encode pricelist %w", err)
	}

	return buf, txml, nil
}
