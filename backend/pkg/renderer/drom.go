package renderer

import (
	"backend/pkg/item"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

type ProductMap map[string]item.Item

// <?xml version="1.0" encoding="UTF-8"?>
// <offers>
//     <offer>
//         <name>Амортизатор</name>
//         <oem_number>334388</oem_number>
//         <analog_numbers>4853033281,4853033291,4853039885,DS2130GS </analog_numbers>
//         <manufacturer>KYB</manufacturer>
//         <ordercode>577</ordercode>
//         <condition>Б/у</condition>
//         <brandcars>Honda</brandcars>
//         <modelcars>Airwave</modelcars>
//         <bodycars>GJ1</bodycars>
//         <engine>1MZFE</engine>
//         <year>2002</year>
//         <lr>лево</lr>
//         <fr>перед</fr>
//         <ud>низ</ud>
//         <color></color>
//     </offer>
// </offers>

type Offer struct {
	Name          string `xml:"name"`           // Наименование что это : Датчик температуры охлаждающей жидкости
	OemNumber     string `xml:"oem_number"`     // номер агрегата
	AnalogNumbers string `xml:"analog_numbers"` // with comma delim артикул,все номера артикула
	Manufacturer  string `xml:"manufacturer"`   // производитель
	Ordercode     string `xml:"ordercode"`

	Condition string `xml:"condition"` // Б/у или новый

	Brandcars string `xml:"brandcars"` // Марка
	Modelcars string `xml:"modelcars"` // Модель
	Engine    string `xml:"engine"`    // двигатель
	Year      string `xml:"year"`      // год

	Photos string `xml:"photos"`

	Amount int `xml:"amount"`
	Price  int `xml:"price"`
}

func photosString(imgs []item.Image) (str string) {
	for i, v := range imgs {
		str += v.Path
		if i < len(imgs)-1 {
			str += ", "
		}
	}

	return str
}

func condition(cond string) string {
	if cond == "RG" {
		return "Б/У"
	}

	return "Новый"
}

func DromRender(items []item.Item) (io.Reader, string, error) {
	offers := make([]*Offer, len(items))
	for idx := 0; len(items) > idx; idx++ {
		if items[idx].Price == 0 {
			continue // пропускаем если нет цены
		}

		offers[idx] = &Offer{
			Ordercode: items[idx].ID,

			Name:          items[idx].Type,
			Manufacturer:  items[idx].Char,
			OemNumber:     items[idx].Name,
			AnalogNumbers: items[idx].SKU,

			Brandcars: items[idx].Brandcars,
			Modelcars: items[idx].Modelcars,
			Engine:    items[idx].Engine,
			Year:      items[idx].Year,

			Condition: condition(items[idx].Char),
			Photos:    photosString(items[idx].Images),

			Amount: items[idx].Amount,
			Price:  items[idx].Price / 100, // убираем копейки
		}
	}

	type Offers struct {
		XMLName xml.Name `xml:"offers"`
		Offers  []*Offer `xml:"offer"`
	}

	nesting := &Offers{}
	nesting.Offers = offers

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(xml.Header)

	enc := xml.NewEncoder(buf)
	enc.Indent("", "  ")

	if err := enc.Encode(nesting); err != nil {
		return nil, "", fmt.Errorf("cant encode pricelist %w", err)
	}

	// b = []byte(xml.Header + string(b))
	// re := bytes.NewReader(b)

	return buf, txml, nil
}
