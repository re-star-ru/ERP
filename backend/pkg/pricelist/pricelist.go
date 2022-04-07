package pricelist

import (
	"backend/pkg/item"
	"backend/pkg/renderer"
	"fmt"
	"io"
	"net/url"
	"path"
)

type consumerMeta struct {
	Name     string
	Path     string
	Consumer func([]item.Item) io.Reader
}

type Itemer interface {
	Items() ([]item.Item, error)
}

type Storer interface {
	Store(key, contentType string , io.Reader) (string, error)
}

type Usecase struct {
	S3path    string
	Consumers map[string]consumerMeta
	Itemer
	store Storer
}

func NewPricerUsecase(store Storer, i Itemer) *Usecase {

	u := &Usecase{
		"https://s3.re-star.ru/pricelists",
		make(map[string]consumerMeta),
		i,
		store,
	}

	// setup consumers
	u.Consumers["drom"] = consumerMeta{
		Name:     "drom",
		Path:     "drom.xml",
		Consumer: renderer.DromRender,
	}

	return u
}

// GetPricelistByServiceName(string) string // return path to s3 pricelist by consumer name
// GetPricelists() map[string]string        // return map [consumer: pricelist]
// Update()

func (s *Usecase) GetPricelistByConsumerName(name string) (string, error) {
	v, ok := s.Consumers[name]
	if !ok {
		return "", fmt.Errorf("no consumer with name %v", name)
	}

	return getPath(s.S3path, v.Name)

}

// GetPricelists returns list with price consumers and paths
func (s *Usecase) GetPricelists() (map[string]string, error) {
	m := make(map[string]string)

	for k, v := range s.Consumers {
		pth, err := getPath(s.S3path, v.Path)
		if err != nil {
			return nil, err
		}

		m[k] = pth
	}

	return m, nil
}

func getPath(basepath, pth string) (string, error) {
	u, err := url.Parse(basepath)
	if err != nil {
		return "", fmt.Errorf("wrong s3path: %w", err)
	}

	u.Path = path.Join(u.Path, pth)

	return u.String(), nil
}

// Upadte download info from onec and update pricelists
func (s *Usecase) Update() error {
	items, err := s.Items()
	if err != nil {
		return err
	}

	for k, v := range s.Consumers {
		document := v.Consumer(items)

		s.store.Store(s)

	}

	panic("update todo")
}
