package categorysvc

type CategoryCreateInput struct {
	CategoryName string
	Budget       float64
	ColorCode    string
	Icon         string
}

type CategoryUpdateInput struct {
	CategoryName string
	Budget       float64
	ColorCode    string
	Icon         string
}
