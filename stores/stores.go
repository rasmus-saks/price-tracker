package stores

import "github.com/op/go-logging"

type ProductPrice struct {
	Price float64
	Shop  string
}

type Product struct {
	Name string
	Ean  string
}

var log = logging.MustGetLogger("stores")

var Stores = [...]struct {
	Name     string
	GetPrice func(product Product) *float64
}{
	{
		Name:     "Selver",
		GetPrice: GetSelverPrice,
	},
	{
		Name:     "Prisma",
		GetPrice: GetPrismaProduct,
	},
}
