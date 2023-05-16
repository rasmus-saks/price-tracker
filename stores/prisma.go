package stores

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"reflect"
	"strconv"
)

func GetPrismaProduct(product Product) *float64 {
	c := colly.NewCollector()
	var prod *float64
	c.OnHTML(".pricebox", func(e *colly.HTMLElement) {
		whole, _ := strconv.Atoi(e.ChildText(".whole-number"))
		dec, _ := strconv.Atoi(e.ChildText(".decimal"))
		num := float64(whole) + (float64(dec) / 100)
		prod = &num
	})

	err := c.Visit(fmt.Sprintf("https://www.prismamarket.ee/entry/%s", product.Ean))

	c.Wait()

	if err != nil && err.Error() != "Not Found" {
		log.Info(reflect.TypeOf(err))
		log.Error(err)
		return nil
	}
	return prod
}
