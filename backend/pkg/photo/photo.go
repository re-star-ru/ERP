package photo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"io"
	"sync"
)

const (
	original  = iota
	thumbnail // 200px
	small     // 400px
	medium    // 800px
	large     // 1200px
)

var ErrSizeNotSupported = fmt.Errorf("size not supported")

type Photo struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	// Фото в разных размерах
	Sizes [5]string `json:"sizes"`
}

func (p Photo) Get() (string, error) {
	return p.Sizes[0], nil
}

func (p Photo) GetBySize(size int) (string, error) {
	if size > len(p.Sizes) || size < 0 {
		return "", fmt.Errorf("%w: %d", ErrSizeNotSupported, size)
	}

	return p.Sizes[size], nil
}

type Storer interface {
	StoreSize(fpath, contentType string, r io.Reader, size int64) (string, error)
}

func NewPhotoUsecase(store Storer) *Usecase {
	return &Usecase{store}
}

type Usecase struct {
	store Storer
}

func (uc *Usecase) NewPhoto(ctx context.Context, file io.Reader, fileSize int64, name string) (Photo, error) {
	// gen uuid
	// store file with uuid/0, uuid/1 ...

	// 1 генерируем uuid
	// 2 - пралельно запускаемт обработку фото и отправку в s3, ждем пока получим все ссылки

	buf, _ := io.ReadAll(file)

	//buf := bufio.NewScanner()
	//
	//io.Copy(w, file)

	var (
		wg    sync.WaitGroup
		photo = Photo{
			ID:   uuid.NewString(),
			Name: name,
		}
	)

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go uc.processFile(ctx, &wg, &photo.Sizes[i], i)
	}

	wg.Wait()

	return photo, nil
}

func (uc *Usecase) processFile(ctx context.Context, wg *sync.WaitGroup, returnPath *string, i int) {
	defer wg.Done()
	// uc.store.StoreSize("", "", file, fileSize)
	log.Debug().Msgf("process file %v, i: %d", *returnPath, i)
}
