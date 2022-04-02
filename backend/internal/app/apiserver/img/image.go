package img

import (
	"backend/internal/app/apiserver/s3"
	"bytes"
	"context"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gabriel-vasile/mimetype"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/spf13/viper"

	"github.com/disintegration/imaging"
)

// Images prepare image to upload on s3: watermark and resize

var spriteImage image.Image
var collection *mongo.Collection
var ctx context.Context

const (
	watermarkWidth  = 400
	watermarkHeight = 400
	image1800       = 1800
	image900        = 900
	opacity         = 0.15
	quality         = 88
)

type Imaging struct {
	spriteImage image.Image
	collection  *mongo.Collection
	ctx         context.Context
}

type ProductImages struct {
	GUID   string // guid of product
	Cover  int    // position of cover image in array
	Images []Image
}

type Image struct {
	Raw      string // Raw image
	Webp1800 string // webp1800px image with wathermark
	Webp900  string // Webp900px image with wathermark
	Jpg1800  string // Jpg1800px image with wathermark
	Jpg900   string // Jpg900px image with wathermark
}

type ImageBin struct {
	Raw      []byte // Raw image
	Webp1800 []byte // webp1800px image with wathermark
	Webp900  []byte // Webp900px image with wathermark
	Jpg1800  []byte // Jpg1800px image with wathermark
	Jpg900   []byte // Jpg900px image with wathermark
}

func Init() *Imaging {
	log.Print("INIT IMG RUN")
	log.Println(viper.GetString("mongo.uri"))
	var err error
	spriteImage, err = imaging.Open("assets/watermark.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	spriteImage = imaging.Resize(spriteImage, watermarkWidth, watermarkHeight, imaging.Lanczos)

	ctx = context.TODO()
	clientOptions := options.Client().ApplyURI(viper.GetString("mongo.uri"))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database(viper.GetString("mongo.dbName")).Collection(viper.GetString("mongo.imageCollection"))
	log.Print("INIT IMG END")

	return &Imaging{
		collection:  collection,
		spriteImage: spriteImage,
		ctx:         ctx,
	}
}

func (i Imaging) Add(guid string, img []byte) error {
	product := ProductImages{}
	if err := i.collection.FindOne(i.ctx, bson.D{{"guid", guid}}).Decode(&product); err != nil && err != mongo.ErrNoDocuments {
		log.Fatal(err)
	}
	product.GUID = guid

	imgBin := prep(img)

	imageStruct := Image{}

	////////////////////// upload files ///////////////////////
	mime := mimetype.Detect(imgBin.Raw)
	path := s3.UploadFile(imgBin.Raw, mime)
	imageStruct.Raw = path

	mime = mimetype.Detect(imgBin.Webp1800)
	path = s3.UploadFile(imgBin.Webp1800, mime)
	imageStruct.Webp1800 = path

	mime = mimetype.Detect(imgBin.Webp900)
	path = s3.UploadFile(imgBin.Webp900, mime)
	imageStruct.Webp900 = path

	mime = mimetype.Detect(imgBin.Jpg1800)
	path = s3.UploadFile(imgBin.Jpg1800, mime)
	imageStruct.Jpg1800 = path

	mime = mimetype.Detect(imgBin.Jpg900)
	path = s3.UploadFile(imgBin.Jpg900, mime)
	imageStruct.Jpg900 = path
	////////////////////////// end upload ///////////////////////

	log.Print(viper.GetString("mongo.imageCollection"))
	product.Images = append(product.Images, imageStruct)
	replaceResult, err := collection.ReplaceOne(ctx, bson.D{{"guid", guid}}, &product)
	if replaceResult.MatchedCount == 0 {
		result, err := collection.InsertOne(ctx, &product)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(result)
	}
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(replaceResult)

	log.Println(product)

	return nil
}

func (i Imaging) Delete(guid string, el int) error {
	product := ProductImages{}
	if err := collection.FindOne(ctx, bson.D{{"guid", guid}}).Decode(&product); err != nil && err != mongo.ErrNoDocuments {
		log.Fatal(err)
	}
	product.Images = append(product.Images[:el], product.Images[el+1:]...)

	// изменить позицию cover если она выше позиции удаленного изображения
	if product.Cover > el {
		product.Cover--
	}
	if product.Cover == el {
		product.Cover = 0
	}

	replaceResult, err := collection.ReplaceOne(ctx, bson.D{{"guid", guid}}, &product)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(replaceResult)
	return nil
}

func (i Imaging) SetCover(guid string, el int) error {
	product := ProductImages{}
	if err := i.collection.FindOne(i.ctx, bson.D{{"guid", guid}}).Decode(&product); err != nil && err != mongo.ErrNoDocuments {
		log.Fatal(err)
	}
	return nil
} // TODO

func Prep(i []byte) []byte {
	img, format, err := image.Decode(bytes.NewReader(i))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Формат ", format)
	dstImage128 := imaging.Fill(img, image900, image900, imaging.Center, imaging.Lanczos)
	dstImage128 = imaging.Overlay(dstImage128, spriteImage, image.Pt(watermarkWidth, watermarkHeight), opacity)
	buf := new(bytes.Buffer)
	if err := imaging.Encode(buf, dstImage128, imaging.JPEG, imaging.JPEGQuality(quality)); err != nil {
		log.Fatalln(err)
	}

	return buf.Bytes()
}

func prep(i []byte) ImageBin {
	imgBin := ImageBin{Raw: i}

	pathln := fmt.Sprintf("%v/resize?width=%d&height=%d&nocrop=false",
		viper.GetString("imageService.path"), image1800, image1800)
	//resize
	resp, err := http.Post(pathln, "image/*", bytes.NewBuffer(i))
	if err != nil {
		log.Fatal(err)
	}

	temp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	pathln = fmt.Sprintf("%v/watermarkimage?image=%v&top=%d&left=%d&opacity=%f",
		viper.GetString("imageService.path"), viper.GetString("imageService.watermarkURL"),
		image1800-watermarkHeight, image1800-watermarkWidth, opacity)
	//watermark
	resp, err = http.Post(pathln, "image/*", bytes.NewBuffer(temp))
	if err != nil {
		log.Fatal(err)
	}

	imgBin.Jpg1800, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//resize
	pathln = fmt.Sprintf("%v/resize?width=%d&height=%d",
		viper.GetString("imageService.path"), image900, image900)
	resp, err = http.Post(pathln, "image/*", bytes.NewBuffer(imgBin.Jpg1800))
	if err != nil {
		log.Fatal(err)
	}
	imgBin.Jpg900, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//convert
	pathln = fmt.Sprintf("%v/convert?type=webp",
		viper.GetString("imageService.path"))
	resp, err = http.Post(pathln, "image/*", bytes.NewBuffer(imgBin.Jpg1800))
	if err != nil {
		log.Fatal(err)
	}
	imgBin.Webp1800, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//convert
	pathln = fmt.Sprintf("%v/convert?type=webp",
		viper.GetString("imageService.path"))
	resp, err = http.Post(pathln, "image/*", bytes.NewBuffer(imgBin.Jpg900))
	if err != nil {
		log.Fatal(err)
	}
	imgBin.Webp900, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("./test.webp", imgBin.Webp1800, 0644); err != nil {
		log.Fatal(err)
	}

	return imgBin
}

func ConvertImage(i []byte) []byte {
	return nil
}

func SetWaterMark() {

}
