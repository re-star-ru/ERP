package restaritem

import (
	"backend/pkg/defect"
	"backend/pkg/employee"
	"backend/pkg/photo"
	"backend/pkg/work"
)

// создаем restoreitem
// отдаем qr код на запись в базе данных
// есть родительский итем из базы который означает приоритетный номер

type RestarItem struct {
	ID   int    `json:"id"` // этот ид идет в qr код
	Name string `json:"name"`
	SKU  string `json:"sku"`
	Char string `json:"char"`

	Description string            `json:"description"` // текстовое описание дефектов (можно набрать голосом)
	Inspector   employee.Employee `json:"inspector"`   // тот кто проводил инспекцию (дефектовку) товара
	Inspection  []defect.Defect   `json:"inspection"`  // список обнаруженых дефектов

	Works []struct {
		Employee employee.Employee `json:"employee"`
		Work     work.Work         `json:"work"`
		Price    int               `json:"price"`
	}

	Photos []photo.Photo `json:"photos"`
}
