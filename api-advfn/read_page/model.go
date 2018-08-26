package read_page

type stock struct {
	Price float64
}

type option struct {
	Kind       string // Call or Put
	Style      string //American or European
	Expiration int
	Url        string
	Price      float64
	Strike     float64
	Stock      stock
}
