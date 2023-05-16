package stores

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type selverResponse struct {
	Hits struct {
		Hits []struct {
			Source struct {
				Price       float64 `json:"price_incl_tax"`
				Name        string  `json:"name"`
				Description string  `json:"description"`
				Prices      []struct {
					Price           float64 `json:"price"`
					CustomerGroupId int     `json:"customer_group_id"`
					IsDiscount      bool    `json:"is_discount"`
				}
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func GetSelverPrice(product Product) *float64 {
	query := url.QueryEscape(fmt.Sprintf(`{"query":{"bool":{"filter":{"terms":{"product_main_ean":["%s"]}}}}}`, product.Ean))
	reqUrl := "https://www.selver.ee/api/catalog/vue_storefront_catalog_et/product/_search?request=" + query
	r, _ := http.Get(reqUrl)
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	var res selverResponse
	err := json.Unmarshal(body, &res)
	if err != nil {
		log.Error(err)
		return nil
	}
	if len(res.Hits.Hits) == 0 {
		return nil
	}
	prod := res.Hits.Hits[0].Source
	for _, price := range prod.Prices {
		if price.CustomerGroupId == 17 {
			return &price.Price
		}
	}
	return &prod.Price
}
