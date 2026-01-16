package transactionsvc

type TransactionItemCreateInput struct {
	Title      string
	Price      float64
	CategoryID string
}
type TransactionItemUpdateInput struct {
	Title      string
	Price      float64
	CategoryID string
}
