package renderer

import (
	"backend/pkg/item"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

type YandexOffer struct {
	Name          string `xml:"name"`           // Наименование что это : Датчик температуры охлаждающей жидкости
	OemNumber     string `xml:"oem_number"`     // номер агрегата
	AnalogNumbers string `xml:"analog_numbers"` // with comma delim артикул,все номера артикула
	Manufacturer  string `xml:"manufacturer"`   // производитель
	Ordercode     string `xml:"ordercode"`

	Condition string `xml:"condition"` // Б/у или новый

	Brandcars string `xml:"brandcars"` // Марка
	Modelcars string `xml:"modelcars"` // Модель
	Engine    string `xml:"engine"`    // двигатель

	Photos string `xml:"photos"`

	Year   int `xml:"year"` // год
	Amount int `xml:"amount"`
	Price  int `xml:"price"`
}

type Category struct {
	ID    string `xml:"id,attr"`
	Value string `xml:",chardata"`
}

type Shop struct {
	Name    string `xml:"name"`
	Company string `xml:"company"`

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

		offers[idx] = &YandexOffer{
			Ordercode: items[idx].ID,
			Name:      items[idx].Type,
		}
	}

	y := YandexCatalog{
		Date: time.Now().Format(time.RFC3339),
		Shop: &Shop{
			Name:    "RESTAR",
			Company: "Стартеры Генераторы Пятигорск",
			Categories: []*Category{
				{ID: "some id", Value: "value"},
				{ID: "some id 2", Value: "value 2"},
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
