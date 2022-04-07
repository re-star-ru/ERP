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
	Name          string // Наименование что это
	OemNumber     string // номер агрегата
	AnalogNumbers string // with comma delim артикул
	Manufacturer  string // производитель
	Ordercode     string
	Condition     string // Б/у или новый
	Brandcars     string // Марка
	Modelcars     string // Модель
	Engine        string // двигатель
	Year          string // год

	Amount int
	Price  int
}

func DromRender(items []item.Item) (io.Reader, string, error) {
	offers := make([]*Offer, len(items))
	for idx := 0; len(items) > idx; idx++ {
		offers[idx] = &Offer{
			Name:      items[idx].SKU,
			OemNumber: items[idx].Name,
			Ordercode: items[idx].ID,
			Amount:    items[idx].Amount,
			Price:     items[idx].Price,
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

	if err := xml.NewEncoder(buf).Encode(nesting); err != nil {
		return nil, "", fmt.Errorf("cant encode pricelist %w", err)
	}

	// b = []byte(xml.Header + string(b))
	// re := bytes.NewReader(b)

	return buf, txml, nil
}
