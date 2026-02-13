package categorysvc

type CategoryConfig struct {
	Color string
	Icon  string
}

var defaultCategories = map[string]CategoryConfig{
	"อาหาร": {
		Color: "#E86A26",
		Icon:  "pizza",
	},
	"การเดินทาง": {
		Color: "#AB47BC",
		Icon:  "subway",
	},
	"ความบันเทิง": {
		Color: "#EC407A",
		Icon:  "musical-notes",
	},
	"ชอปปิ้ง": {
		Color: "#547BE9",
		Icon:  "bag-handle",
	},
	"อื่นๆ": {
		Color: "#A1887F",
		Icon:  "ellipsis-horizontal-circle-sharp",
	},
}
