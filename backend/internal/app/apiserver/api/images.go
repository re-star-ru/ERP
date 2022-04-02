package api

import (
	"backend/internal/app/apiserver/db"
	"backend/internal/app/apiserver/img"
	"backend/internal/app/apiserver/s3"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gabriel-vasile/mimetype"

	"github.com/go-chi/chi"
)

type Imaging interface {
	Add(string, []byte) error
	Delete(string, int) error
	SetCover(string, int) error
}

func imgInit() {
	imaging = img.Init()
}

var imaging Imaging

func checkMime(mime *mimetype.MIME) (err error) {
	log.Println(mime.String())
	err = errors.New("неправильный тип файла")
	switch mime.String() {
	case "image/jpeg", "image/png", "image/svg+xml", "image/gif":
		log.Println("нет ошибок")
		return nil
	default:
		return err
	}
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(0); err != nil && err.Error() != "multipart: NextPart: EOF" {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	guid := chi.URLParam(r, "GUID")
	if guid == "" {
		err := errors.New("нет guid")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offer, err := db.GetOffer(guid)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	sliceFiles, ok := r.MultipartForm.File["file"]
	if !ok {
		err := errors.New("нет данных картинки")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	formFile := sliceFiles[0]
	tempFile, err := formFile.Open()
	if err != nil {
		log.Fatalln(err)
	}
	f, err := ioutil.ReadAll(tempFile)
	if err != nil {
		log.Fatalln(err)
	}
	mime := mimetype.Detect(f)
	if err := checkMime(mime); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	//imaging := img.ProductImaging{
	//	GUID: offer.GUID,
	//	Name: offer.Name,
	//}
	imaging.Add(offer.GUID, f)
	processedImg := img.Prep(f)
	path := s3.UploadFile(processedImg, mime)
	log.Println(path)

	if err := offer.ImageAdd(path); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(200)

	// закрытие файлов и удаление данных формы
	defer func() {
		if err := r.MultipartForm.RemoveAll(); err != nil {
			log.Println("Не удалось удалить или хз ", err)
		}
	}()
	defer func() {
		if err := tempFile.Close(); err != nil {
			log.Println("Не удаляется насяльника ", err)
		}
	}()
}

func deleteImage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	link := struct {
		Link string `json:"link"`
		El   int    `json:"el"`
	}{}
	if err := json.Unmarshal(body, &link); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	guid := chi.URLParam(r, "GUID")
	if guid == "" {
		err := errors.New("пустой guid")
		http.Error(w, err.Error(), 400)
		return
	}

	o, err := db.GetOffer(guid)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	log.Println("Link", link)
	imaging.Delete(o.GUID, link.El)
	if err := o.ImageDelete(link.Link); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.WriteHeader(http.StatusOK)
}
