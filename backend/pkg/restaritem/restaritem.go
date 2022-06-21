package restaritem

import (
	"backend/pkg/photo"
	"errors"
)

// создаем restoreitem
// отдаем qr код на запись в базе данных
// есть родительский итем из базы который означает приоритетный номер

//type RestarItem struct {
//	ID       int    `json:"id"`     // этот ид идет в qr код
//	OnecID   int    `json:"onecId"` // ид в базе данных 1с, какой талон ремонта или еще что то создал
//	Name     string `json:"name"`
//	SKU      string `json:"sku"`
//	ItemGUID string `json:"itemGUID"` // ид товара в базе 1с
//	CharGUID string `json:"char"`     // ид характеристики в базе 1с
//
//	Description string            `json:"description"` // текстовое описание дефектов (можно набрать голосом)
//	Inspector   employee.Employee `json:"inspector"`   // тот кто проводил инспекцию (дефектовку) товара
//	Inspection  []defect.Defect   `json:"inspection"`  // список обнаруженых дефектов
//
//	Works []struct {
//		Employee employee.Employee `json:"employee"`
//		Work     work.Work         `json:"work"`
//		Price    int               `json:"price"`
//	}
//
//	Photos []photo.Photo `json:"photos"`
//}

type RestarItem struct {
	ID       string `json:"id"`       // этот ид идет в qr код
	OnecGUID string `json:"onceGUID"` // ид в базе данных 1с, какой талон ремонта или еще что то создал
	Name     string `json:"name"`
	SKU      string `json:"sku"`
	ItemGUID string `json:"itemGUID"` // ид товара в базе 1с
	CharGUID string `json:"charGUID"` // ид характеристики в базе 1с

	Description string   `json:"description"` // текстовое описание дефектов (можно набрать голосом) GUID из 1с
	Inspector   string   `json:"inspector"`   // тот кто проводил инспекцию (дефектовку) товара GUID из 1с
	Inspection  []string `json:"inspection"`  // список обнаруженых дефектов GUID из 1с

	Works  []Work        `json:"works"`
	Photos []photo.Photo `json:"photos"`
}

type Work struct {
	EmployeeGUID string `json:"employee"`
	WorkGUID     string `json:"work"`
	Price        int    `json:"price"`
}

var ErrNotFound = errors.New("not found")
var ErrValidation = errors.New("validation error")
