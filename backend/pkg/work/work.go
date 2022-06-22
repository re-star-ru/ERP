package work

type Work struct {
	EmployeeGUID string `json:"employee"`
	WorkGUID     string `json:"work"`
	Price        int    `json:"price"`
}
