package pricelist

import (
	"backend/pkg/item"
	"backend/pkg/renderer"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"
)

type consumerMeta struct {
	Name   string
	Path   string
	Render func(items []item.Item) (r io.Reader, contentType string, err error)
}

type Itemer interface {
	Items() ([]item.Item, error)
}

type Storer interface {
	Store(key string, contentType string, r io.Reader) (path string, err error)
}

type Usecase struct {
	S3path    string
	Consumers map[string]consumerMeta
	Itemer
	store Storer
}

func NewPricerUsecase(store Storer, i Itemer) *Usecase {
	ucase := &Usecase{
		"https://s3.re-star.ru/pricelists",
		make(map[string]consumerMeta),
		i,
		store,
	}

	// setup consumers
	ucase.Consumers["drom"] = consumerMeta{
		Name:   "drom",
		Path:   "drom.xml",
		Render: renderer.DromRender,
	}

	return ucase
}

var ErrNoSuchConsumer = errors.New("no such consumer")

func (s *Usecase) GetPricelistByConsumerName(name string) (string, error) {
	v, ok := s.Consumers[name]
	if !ok {
		return "", fmt.Errorf("%w: %v", ErrNoSuchConsumer, name)
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

// Update download info from onec and update pricelists
func (s *Usecase) Update() (err error) {
	items, err := s.Items()
	if err != nil {
		return err
	}

	var document io.Reader

	var contentType string

	for _, consumer := range s.Consumers {
		document, contentType, err = consumer.Render(items)
		if err != nil {
			return fmt.Errorf("cant render consumer %v, got error: %w", consumer.Name, err)
		}

		if _, err = s.store.Store(path.Join(s.S3path, consumer.Path), contentType, document); err != nil {
			return fmt.Errorf("cant store consumer %v, got error: %w", consumer.Name, err)
		}
	}

	return nil
}
