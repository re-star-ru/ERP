package pricelist

import (
	"backend/pkg/item"
	"backend/pkg/renderer"
	"errors"
	"fmt"
	"io"
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
	Path() string
}

type Usecase struct {
	S3Directory string
	Consumers   map[string]consumerMeta
	it          Itemer
	store       Storer
}

func NewPricerUsecase(store Storer, i Itemer) *Usecase {
	ucase := &Usecase{
		"preicelists",
		make(map[string]consumerMeta),
		i, store,
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

	return fmt.Sprintf("%s/%s/%s", s.store.Path(), s.S3Directory, v.Path), nil
}

// GetPricelists returns list with price consumers and paths
func (s *Usecase) GetPricelists() (map[string]string, error) {
	m := make(map[string]string)

	for k, v := range s.Consumers {
		m[k] = fmt.Sprintf("%s/%s/%s", s.store.Path(), s.S3Directory, v.Path)
	}

	return m, nil
}

// Update download info from onec and update pricelists
func (s *Usecase) Update() (err error) {
	items, err := s.it.Items()
	if err != nil {
		return err
	}

	var (
		document    io.Reader
		contentType string
	)

	for _, consumer := range s.Consumers {
		document, contentType, err = consumer.Render(items)
		if err != nil {
			return fmt.Errorf("cant render consumer %v, got error: %w", consumer.Name, err)
		}

		if _, err = s.store.Store(path.Join(s.S3Directory, consumer.Path), contentType, document); err != nil {
			return fmt.Errorf("cant store consumer %v, got error: %w", consumer.Name, err)
		}
	}

	return nil
}
