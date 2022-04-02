package catalog

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func searchSKUIn1c(searchQuery string) []string {
	//skus := make([]string, 0)
	//normalize input string
	searchQuery = strings.ReplaceAll(searchQuery, " ", "")
	log.Println(searchQuery)
	c := &http.Client{}
	req, err := http.NewRequest("GET", viper.GetString("srv1sv8.path")+"/products/text-search/"+searchQuery, nil)
	if err != nil {
		log.Println(err)
		return []string{}
	}
	req.SetBasicAuth(viper.GetString("srv1sv8.login"), viper.GetString("srv1sv8.password"))

	start := time.Now()
	res, err := c.Do(req)
	log.Println(time.Since(start))
	if err != nil {
		log.Println(err)
		return []string{}
	}
	if res.StatusCode != http.StatusOK {
		err := errors.New("Ошибка в ответе: " + res.Status)
		log.Println(err, res.Status)
		return []string{}
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return []string{}
	}
	log.Println(len(body))

	resp := struct {
		FoundSKU []string `json:"foundSKU"`
	}{}
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Println(err)
		return []string{}
	}

	return resp.FoundSKU
}
