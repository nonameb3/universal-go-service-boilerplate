package request

type GetItem struct {
	Id string `json:"id"`
}

type AddItem struct {
	Name   string `json:"name"`
	Amount uint   `json:"amount"`
}
