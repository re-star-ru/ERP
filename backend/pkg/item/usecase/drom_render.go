package usecase

import (
	"backend/pkg/item"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
)

func NewDromRender(bucket string, m *minio.Client) *DromRender {
	return &DromRender{
		c:        m,
		bucket:   bucket,
		filePath: "drom.xml",
	}
}

type DromRender struct {
	c                *minio.Client
	bucket, filePath string
}

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

func (r *DromRender) Render(price map[string]item.Item) error {
	offers := make([]*Offer, len(price))
	i := 0
	for _, v := range price {
		offers[i] = &Offer{
			Name:      v.SKU,
			OemNumber: v.Name,
			Ordercode: v.ID,
			Amount:    v.Amount,
			Price:     v.Price,
		}

		i++
	}

	type Offers struct {
		XMLName xml.Name `xml:"offers"`
		Offers  []*Offer `xml:"offer"`
	}

	nesting := &Offers{}
	nesting.Offers = offers

	b, err := xml.MarshalIndent(nesting, " ", " ")
	if err != nil {
		return fmt.Errorf("cant encode pricelist %w", err)
	}
	b = []byte(xml.Header + string(b))
	re := bytes.NewReader(b)

	info, err := r.c.PutObject(
		context.Background(),
		r.bucket,
		r.filePath,
		re,
		-1,
		minio.PutObjectOptions{ContentType: "text/xml"},
	)

	if err != nil {
		return fmt.Errorf("ошибка загрузки DromRender %s/%s: %w", r.bucket, r.filePath, err)
	}

	log.Info().Str("путь к файлу выгрузки", path.Join(info.Bucket, info.Key)).Send()

	return nil
}
