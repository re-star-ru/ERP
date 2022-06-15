package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/spf13/viper"

	"github.com/asdine/storm"
)

// Group это структура с артикулами
type Group struct {
	ID  int    `storm:"increment"`
	SKU string `json:"sku" storm:"index"`

	MainProductName string `json:"mainProductName"`
	MainProductGUID string `json:"mainProductGUID"`

	MainProductTypeName string `json:"mainProductTypeName"`
	MainProductTypeGUID string `json:"mainProductTypeGUID"`

	MainProductProperties []Property `json:"mainProductProperties"`

	Products []Product `json:"catalog"`
	//
	//SkuConfig struct {
	//} `json:"skuConfig"`
}

// Product это продукт
type Product struct {
	Name         string `json:"name"`
	GUID         string `json:"guid" storm:"id"`
	SKU          string `json:"sku" storm:"index"`
	Description  string `json:"description"`
	Manufacturer string `json:"manufacturer"`
	TypeGUID     string `json:"typeGUID"`
	TypeName     string `json:"typeName"`

	Characteristics []Characteristic `json:"characteristics"`

	Properties []Property `json:"properties"`
}

type Characteristic struct {
	GUID  string `json:"characteristic" storm:"id"`
	Name  string `json:"characteristicName" storm:"index"`
	Owner string `json:"characteristicOwner"`
}

type Property struct {
	Name  string `json:"propertyName" storm:"index"`
	GUID  string `json:"property"`
	Value string `json:"value"`
	Unit  string `json:"unit"`
}

// ProductType это вид продукта который содержит все реквизиты для фильтрации для
// этого типа
type ProductType struct {
	GUID       string `json:"guid" storm:"id"`
	Name       string `json:"name" storm:"index"`
	Parent     string `json:"parent"`
	IsGroup    bool   `json:"group"`
	Properties []struct {
		PropertyName string `json:"propertyName"`
		Property     string `json:"property"`
	} `json:"properties"`
}

// ProductStore это хранилище для этого пакета
type ProductStore struct {
	Products     []Product     `json:"products"`
	ProductTypes []ProductType `json:"productTypes"`
	DB           *storm.DB     `json:"-"`
	mongoDB      *mongo.Database

	productsMap map[string]*Product //map for fast catalog
	skuMap      map[string][]string // map for fast search by skus
}

var productStore ProductStore

func InitDB() {
	var err error

	if productStore.DB, err = storm.Open(viper.GetString("DBCatalogPath")); err != nil {
		log.Fatalln(err)
	}

	clientOptions := options.Client().ApplyURI(viper.GetString("mongo.uri"))

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalln(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("==== MongoDB connected ====")

	productStore.mongoDB = client.Database(viper.GetString("mongo.dbName"))
	test := productStore.mongoDB.Collection("test")

	startuptest := struct {
		Timestamp time.Time
		Date      string
	}{
		time.Now(),
		time.Now().Format("2006.01.02 15:04:05.000"),
	}
	_, err = test.InsertOne(context.TODO(), startuptest)
	if err != nil {
		log.Fatalln(err)
	}

}

func (p *ProductStore) update() error {

	// Дропаем все учавствующие коллекции
	catalog := productStore.mongoDB.Collection(viper.GetString("mongo.catalogCollection"))
	if err := catalog.Drop(context.TODO()); err != nil {
		log.Println(err)
		return err
	}
	productsCollection := productStore.mongoDB.Collection(viper.GetString("mongo.productsCollection"))
	if err := productsCollection.Drop(context.TODO()); err != nil {
		log.Println(err)
		return err
	}
	productTypesCollection := productStore.mongoDB.Collection(viper.GetString("mongo.productTypesCollection"))
	if err := productStore.mongoDB.Collection("productTypesTree").Drop(context.TODO()); err != nil {
		log.Println(err)
		return err
	}
	//// Конец очистки коллекций

	c := &http.Client{}
	req, err := http.NewRequest("GET", viper.GetString("srv1sv8.path")+"/products", nil)
	if err != nil {
		log.Println(err)
		return err
	}
	req.SetBasicAuth(viper.GetString("srv1sv8.login"), viper.GetString("srv1sv8.password"))

	start := time.Now()
	res, err := c.Do(req)
	log.Println(time.Since(start))
	if err != nil {
		log.Println(err)
		return err
	}
	if res.StatusCode != http.StatusOK {
		err := errors.New("Ошибка в ответе: " + res.Status)
		log.Println(err, res.Status)
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(len(body))
	if err := json.Unmarshal(body, p); err != nil {
		log.Println(err)
	}
	// TODO получить сразу слайс интферфейсов вместо анмаршалинга потом
	productsSlice := make([]interface{}, len(p.Products))
	for i, d := range p.Products {
		productsSlice[i] = d
	}

	productTypesSlice := make([]interface{}, len(p.ProductTypes))
	for i, d := range p.ProductTypes {
		productTypesSlice[i] = d
	}

	////// TODO преобразовать список продуктов в аггрегат Group
	// TODO Сохранить products в mongodb
	t := time.Now()
	_, err = productsCollection.InsertMany(context.TODO(), productsSlice)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(time.Since(t))

	// TODO Сохранить productTypes в mongodb
	t = time.Now()
	_, err = productTypesCollection.InsertMany(context.TODO(), productTypesSlice)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(time.Since(t))
	// TODO преобразовать productTypes  в дерево

	cursor, _ := productStore.mongoDB.Collection(viper.GetString("mongo.productTypesCollection")).Find(context.TODO(), bson.D{}, options.Find().SetProjection(bson.M{"_id": 0}))

	productTypesSlice1 := make([]interface{}, 0)
	for cursor.Next(context.TODO()) {
		result := bson.M{}
		if err := cursor.Decode(&result); err != nil {
			log.Println("cursor.Next() error:", err)
			return err
		}
		productTypesSlice1 = append(productTypesSlice1, result)
	}

	productTypeTree := makeProductTypesTree(productTypesSlice1)

	if _, err := productStore.mongoDB.Collection("productTypesTree").InsertOne(context.TODO(), &productTypeTree); err != nil {
		log.Println(err)
		return err
	}

	// создаем мапу по артикулам
	p.createSKUMap()
	p.createProductsMap()

	groups := productStore.getSKUProductList()

	respGroups := make([]interface{}, 0)
	for _, v := range groups {
		respGroups = append(respGroups, v)
	}

	_, err = catalog.InsertMany(context.TODO(), respGroups)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(time.Since(t))
	mod := mongo.IndexModel{
		Keys:    bson.M{"sku": 1},
		Options: nil,
	}
	ind, err := catalog.Indexes().CreateOne(context.TODO(), mod)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(ind)

	return nil
}

func makeProductTypesTree(ts []interface{}) map[string]interface{} {

	tree := make(map[string]interface{})

	rootParent := "00000000-0000-0000-0000-000000000000"

	obj := make(map[string]interface{})
	for _, v := range ts {
		if v.(bson.M)["guid"] != nil && v.(bson.M)["guid"].(string) != "" {
			obj[v.(bson.M)["guid"].(string)] = v
		}
	}

	for _, v := range obj {
		if v.(bson.M)["parent"] == rootParent {
			tree[v.(bson.M)["guid"].(string)] = v
			continue
		}
		t := obj[v.(bson.M)["parent"].(string)] // мб вернуть указатель
		if t.(bson.M)["children"] == nil {
			t.(bson.M)["children"] = make(bson.M)
		}
		t.(bson.M)["children"].(bson.M)[v.(bson.M)["guid"].(string)] = v.(bson.M)
	}

	return tree
}

//func createAggregationGroup(ps []Product) []Group {
//	aggGroups := make([]Group, 0)
//	groupedBySKU := groupBySKU(ps)
//	log.Println(len(groupedBySKU))
//
//	for key, element := range groupedBySKU {
//		group := Group{}
//		group.SKU = key
//		group.Products = element
//
//		// поиск подходящего главного продукта для артикула
//		for _, v := range group.Products {
//			if v.Name == group.SKU {
//				group.MainProductName = v.Name
//				group.MainProductGUID = v.GUID
//				group.MainProductTypeName = v.TypeName
//				group.MainProductTypeGUID = v.TypeGUID
//				break
//			}
//			if group.MainProductName == "" {
//				group.MainProductName = v.Name
//				group.MainProductGUID = v.GUID
//				group.MainProductTypeName = v.TypeName
//				group.MainProductTypeGUID = v.TypeGUID
//			}
//		}
//
//		aggGroups = append(aggGroups, group)
//	}
//
//	return aggGroups
//}

func (p *ProductStore) createProductsMap() {
	log.Println("create productsMap")
	p.productsMap = make(map[string]*Product)

	for i := range p.Products {
		p.productsMap[p.Products[i].GUID] = &p.Products[i]
	}
}

func (p *ProductStore) createSKUMap() {
	log.Println("create createSKUMap")

	p.skuMap = make(map[string][]string)

	for _, v := range p.Products {
		arr := append(p.skuMap[v.SKU], v.GUID)
		p.skuMap[v.SKU] = arr

	}
}

//func groupBySKU(ps []Product) map[string][]Product {
//	uniqSKUMap := make(map[string][]Product)
//	for _, p := range ps {
//		skuPs := append(uniqSKUMap[p.SKU], p)
//		uniqSKUMap[p.SKU] = skuPs
//	}
//	return uniqSKUMap
//}

func (p *ProductStore) getSKUProductList() (skuList []Group) {
	for key, element := range p.skuMap {
		listItem := Group{}
		listItem.SKU = key

		for _, v := range element {
			listItem.Products = append(listItem.Products, *p.productsMap[v])
		}

		// TODO: переделать чтобы было красиво
		// поиск подходящего главного продукта для артикула
		for _, v := range listItem.Products {
			if v.Name == listItem.SKU {
				listItem.MainProductName = v.Name
				listItem.MainProductGUID = v.GUID
				listItem.MainProductTypeName = v.TypeName
				listItem.MainProductTypeGUID = v.TypeGUID
				listItem.MainProductProperties = v.Properties
				break
			}
			if listItem.MainProductName == "" {
				listItem.MainProductName = v.Name
				listItem.MainProductGUID = v.GUID
				listItem.MainProductTypeName = v.TypeName
				listItem.MainProductTypeGUID = v.TypeGUID
				listItem.MainProductProperties = v.Properties
			}
		}

		skuList = append(skuList, listItem)
	}

	return skuList
}
