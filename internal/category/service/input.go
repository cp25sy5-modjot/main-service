package categorysvc

type CategoryCreateInput struct {
	CategoryName string
	Budget       float64
	ColorCode    string
}

type CategoryUpdateInput struct {
	CategoryName string
	Budget       float64
	ColorCode    string
}
