package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-chi/chi"
)

type TestSlicer struct {
	ID          string
	IndexString string `storm:"index"`
	Owner       string
}

// UpdateProductStore хендлер для обновления продуктов
func UpdateProductStore(w http.ResponseWriter, _ *http.Request) {
	log.Println("Update product store")

	err := productStore.update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(http.StatusText(http.StatusOK))); err != nil {
		log.Println(err)
	}
}

func FilterOwners(vs []TestSlicer) []string {
	guids := make([]string, 0)
	for _, v := range vs {
		guids = append(guids, v.Owner)
	}
	return guids
}

// GetInfo получение инфо о текущих продуктах
func GetInfo(w http.ResponseWriter, _ *http.Request) {
	log.Println("GetInfo")
	err := productStore.update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(productStore)
	if err != nil {
		log.Println(err)
		return
	}

	if _, err := w.Write(data); err != nil {
		log.Println(err)
	}
}

// GetDefaultGroupList
func GetDefaultGroupList(w http.ResponseWriter, r *http.Request) {
	log.Println("get catalog default")

	log.Println(r.URL.Query())
	page := int64(0)
	if r.URL.Query().Get("p") != "" {
		page, _ = strconv.ParseInt(r.URL.Query().Get("p"), 10, 64)
		if page > 0 {
			page--
		}
	}

	productType := ""
	if r.URL.Query().Get("t") != "" {
		productType = r.URL.Query().Get("t")
	}

	searchQuery := r.URL.Query().Get("q")
	skuSlice := make([]string, 0)
	if searchQuery != "" {
		log.Println("search:", searchQuery)
		skuSlice = searchSKUIn1c(searchQuery)
		log.Println(skuSlice)
		if len(skuSlice) == 0 {
			err := errors.New("nothing found for this request")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	}

	log.Println("request for page:", page)
	defPerPage := int64(48)
	defSkip := defPerPage * page
	//groups := make([]Group, 0)
	//if err := productStore.DB.All(&groups); err != nil {
	//	log.Println(err)
	//}
	//var respGroups []interface{}
	//for _, v := range groups {
	//	respGroups = append(respGroups, v)
	//}

	///// TODO MONGO
	//t := time.Now()
	catalog := productStore.mongoDB.Collection(viper.GetString("mongo.catalogCollection"))
	//products := db.Collection(viper.GetString("productsCollection"))
	//orders := db.Collection(viper.GetString("ordersCollection"))

	//t := time.Now()
	//_, err := catalog.InsertMany(context.TODO(), respGroups)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//log.Println(time.Since(t))
	//mod := mongo.IndexModel{
	//	Keys:    bson.M{"sku": 1},
	//	Options: nil,
	//}
	//ind, err := catalog.Indexes().CreateOne(context.TODO(), mod)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//log.Println("CreateOne() index:", ind)
	//log.Println("CreateOne() type:", reflect.TypeOf(ind))
	t := time.Now()
	filter := bson.M{}
	if productType != "" {
		filter["mainproducttypeguid"] = productType
	}
	if len(skuSlice) > 0 {
		skuA := bson.A{}
		for _, v := range skuSlice {
			skuA = append(skuA, v)
		}
		filter["sku"] = bson.M{"$in": skuA}
	}

	cursor, err := catalog.Find(
		context.TODO(),
		filter,
		options.Find().SetSort(bson.M{"sku": 1}),
		options.Find().SetProjection(bson.D{{Key: "_id", Value: 0}}),
		options.Find().SetLimit(defPerPage),
		options.Find().SetSkip(defSkip),
	)

	if err != nil {
		log.Println(err)
	}

	count, err := catalog.CountDocuments(
		context.TODO(),
		filter,
	)
	if err != nil {
		log.Println(err)
	}
	var docs []bson.M
	if err = cursor.All(context.TODO(), &docs); err != nil {
		log.Println(err)
	}
	log.Println(time.Since(t))

	//// TODO MONGO

	respData := struct {
		Count   int64    `json:"count"`
		Groups  []bson.M `json:"groups"`
		PerPage int64    `json:"perPage"`
	}{
		count,
		docs,
		defPerPage,
	}

	data, err := json.Marshal(respData)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, _ = w.Write(data)
}

// GetSKUList получение списка артикулов
func GetAllSKUList(w http.ResponseWriter, _ *http.Request) {
	log.Println("GetSKUList")
	skuList := productStore.getSKUProductList()
	//if err != nil {
	//	log.Print(err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	data, err := json.Marshal(&skuList)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(data); err != nil {
		log.Println(err)
	}
}

// GetSKUList получение списка артикулов по группе  товара и оффсету
func GetSKUListByTypeAndOffset(w http.ResponseWriter, r *http.Request) {
	log.Println("GetSKUList")
	log.Println(chi.URLParam(r, "productTypeGuid"))

	offset := 48
	if chi.URLParam(r, "offset") != "" {
		offset, _ = strconv.Atoi(chi.URLParam(r, "offset"))
	}
	log.Println(offset)
	log.Println(chi.URLParam(r, "offset"))
	skuList := productStore.getSKUProductList()
	//if err != nil {
	//	log.Print(err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	data, err := json.Marshal(&skuList)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(data); err != nil {
		log.Println(err)
	}
}

// GetProductTypes получение списка артикулов
func GetProductTypes(w http.ResponseWriter, _ *http.Request) {
	log.Println("GetProductTypes")

	productTypeTree := bson.M{}
	if err := productStore.mongoDB.Collection("productTypesTree").FindOne(context.TODO(), bson.M{}, options.FindOne().SetProjection(bson.M{"_id": 0})).Decode(&productTypeTree); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(&productTypeTree)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(data); err != nil {
		log.Println(err)
	}
}

// GetProductListParams получение списка артикулов по параметрам
func GetProductListParams(w http.ResponseWriter, r *http.Request) {
	log.Println("Get catalog by params")
	type sku string
	skuNode := productStore.DB.From("sku")
	sku1 := sku("cs150")
	sku2 := sku("js150")
	if err := skuNode.Save(sku1); err != nil {
		log.Println(err)
	}
	if err := skuNode.Save(sku2); err != nil {
		log.Println(err)
	}
	skus := make([]sku, 0)
	if err := skuNode.Find("sku", "cs150", &skus); err != nil {
		log.Println(err)
	}
	log.Println(skus)

	//log.Println(r.URL.Query().Get("q"))

	log.Println(chi.URLParam(r, "*"))
	query := strings.Split(chi.URLParam(r, "*"), "/")
	log.Println(query)

	var skuList []Product
	//if haveProductType(query) {
	//	log.Println("type")
	//	pType := getPType(query)
	//	pTypes := "Габариты"
	//	log.Println(pType)
	//	if err := productStore.DB.Find("PropertyName", pTypes, &skuList, storm.Limit(10)); err != nil {
	//		log.Println(err)
	//	}
	//}

	//skuList := productStore.getSKUProductList()
	//if err != nil {
	//	log.Print(err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	data, err := json.Marshal(skuList)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(data); err != nil {
		log.Println(err)
	}
}

//func haveProductType(a []string) bool {
//	x := "t"
//	for _, n := range a {
//		first := n[0:1]
//		if x == first {
//			return true
//		}
//	}
//	return false
//}

//func getParams(a []string) error {
//	return nil
//}

//func getPType(q []string) int {
//	return 0
//}
