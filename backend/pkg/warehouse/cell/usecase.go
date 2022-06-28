package cell

type Cell struct {
	ID    string
	Items []Item `json:"body"`
}

type Item struct {
	Index              int      `json:"index"`
	SKU                string   `json:"sku"`
	CharacteristicGUID string   `json:"characteristicGUID"`
	Name               string   `json:"name"`
	Characteristic     string   `json:"characteristic"`
	Amount             int      `json:"amount"`
	Images             []string `json:"images"`
}

type ICellRepo interface {
	GetCellByID(cellID string) (*Cell, error)
}

type Usecase struct {
	repo ICellRepo
}

func NewCellUsecase(repo ICellRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (c *Usecase) ByID(cellID string) (*Cell, error) {
	return c.repo.GetCellByID(cellID)
}
