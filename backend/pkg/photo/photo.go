package photo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"strconv"
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
	ID string `json:"id"`

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
	phuc := &Usecase{store: store}

	var (
		bpreset []byte
		err     error
	)
	for i, v := range presets {

		bpreset, err = json.Marshal(append(v, convertJpeg))
		if err != nil {
			log.Fatal().Err(err).Msg("can't marshal preset")
		}
		phuc.presets[i] = append(phuc.presets[i], url.Values{"operations": {string(bpreset[:])}}.Encode())

		bpreset, err = json.Marshal(append(v, convertWebp))
		if err != nil {
			log.Fatal().Err(err).Msg("can't marshal preset")
		}
		phuc.presets[i] = append(phuc.presets[i], url.Values{"operations": {string(bpreset[:])}}.Encode())
	}

	return phuc
}

type Usecase struct {
	store   Storer
	presets [4][]string
}

func (uc *Usecase) NewPhoto(ctx context.Context, dir string, photo io.ReadCloser) (Photo, error) {
	// gen uuid
	// store file with uuid/0, uuid/1 ...

	// 1 генерируем uuid
	// 2 - пралельно запускаемт обработку фото и отправку в s3, ждем пока получим все ссылки
	buf, err := io.ReadAll(photo)
	if err != nil {
		photo.Close()
		return Photo{}, fmt.Errorf("cant read file: %w", err)
	}
	photo.Close()

	var (
		wg     sync.WaitGroup
		nPhoto = Photo{
			ID: uuid.NewString(),
		}
	)

	path, err := uc.store.StoreSize(dir+"/"+nPhoto.ID+"/0.jpg", "image/jpg", bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		return Photo{}, fmt.Errorf("cant store file: %w", err)
	}
	nPhoto.Sizes[original] = path

	for i := 1; i < 5; i++ {
		wg.Add(2)

		storePath := dir + "/" + nPhoto.ID + "/" + strconv.Itoa(i)
		jpgPath := storePath + ".jpg"
		webpPath := storePath + ".webp"

		go uc.processFile(ctx, &wg, &buf, &nPhoto.Sizes[i], uc.presets[i-1][0], jpgPath)
		go uc.processFile(ctx, &wg, &buf, &nPhoto.Sizes[i], uc.presets[i-1][1], webpPath)
	}

	wg.Wait()

	return nPhoto, nil
}

func (uc *Usecase) processFile(ctx context.Context, wg *sync.WaitGroup, file *[]byte, returnPath *string, preset string, storePath string) {
	defer wg.Done()

	cmdPath := "http://localhost:5556/pipeline"
	cmdPath += "?" + preset

	req, err := http.NewRequestWithContext(ctx, "POST", cmdPath, bytes.NewReader(*file))
	if err != nil {
		log.Error().Err(err).Msg("can't create request")
		return
	}

	req.Header.Add("Content-Type", "image/*")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("can't do request")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("can't read response")

			return
		}

		log.Log().Msgf("resp: %v", string(body))

		return
	}

	path, err := uc.store.StoreSize(storePath, "image/webp", resp.Body, resp.ContentLength)
	if err != nil {
		log.Error().Err(err).Msg("can't store file")

		return
	}
	*returnPath = path

	log.Debug().Msgf("process file %v", *returnPath)
}

//////////////// presets configs

var presets = [4][]Operation{
	presetThumbnail,
	presetSmall,
	presetMedium,
	presetLarge,
}

type Operation struct {
	Operation     string                 `json:"operation"`
	IgnoreFailure bool                   `json:"ignore_failure"`
	Params        map[string]interface{} `json:"params"`
}

var watermark = Operation{
	Operation:     "watermarkImage",
	IgnoreFailure: false,
	Params: map[string]interface{}{
		"image":   "https://s3.re-star.ru/oprox/watermark.png",
		"top":     "0",
		"left":    "0",
		"opacity": "0.15",
	},
}

var convertWebp = Operation{
	Operation: "convert",
	Params: map[string]interface{}{
		"type":        "webp",
		"aspectratio": "4:3",
	},
}

var convertJpeg = Operation{
	Operation: "convert",
	Params: map[string]interface{}{
		"type": "webp",
	},
}

var presetThumbnail = []Operation{
	{
		Operation:     "resize",
		IgnoreFailure: false,
		Params: map[string]interface{}{
			"width":       "200",
			"aspectRatio": "4:3",
		},
	},
	convertWebp,
}

var presetSmall = []Operation{
	watermark,
	{
		Operation:     "resize",
		IgnoreFailure: false,
		Params: map[string]interface{}{
			"width":       "400",
			"aspectRatio": "4:3",
		},
	},
	convertWebp,
}

var presetMedium = []Operation{
	{
		Operation:     "resize",
		IgnoreFailure: false,
		Params: map[string]interface{}{
			"width":       "800",
			"aspectratio": "4:3",
			"nocrop":      false,
		},
	},
	watermark,
	convertWebp,
}

var presetLarge = []Operation{
	watermark,
	{
		Operation:     "resize",
		IgnoreFailure: false,
		Params: map[string]interface{}{
			"width":       "1200",
			"aspectRatio": "4:3",
		},
	},
	convertWebp,
}
