package main

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/op/go-logging"
	"github.com/rasmus-saks/price-tracker/stores"
	"net/http"
	"time"
)

var log = logging.MustGetLogger("price-tracker")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} %{level:.4s} %{color:reset} %{message}`,
)

type TrackedProduct struct {
	Name      string
	Ean       string
	LastPrice map[string]*float64
}

func main() {
	logging.SetFormatter(format)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WritePrometheus(w, false)
	})
	for i := range products {
		product := &products[i]
		product.LastPrice = make(map[string]*float64, len(stores.Stores))
		for _, store := range stores.Stores {
			registerGauge(product, store.Name)
		}
	}

	autoRefresh()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func autoRefresh() {
	refreshPrices()
	go func() {
		for {
			time.Sleep(60 * time.Minute)
			refreshPrices()
		}
	}()
}

func refreshPrices() {
	log.Info("Refreshing prices")
	for i := range products {
		product := &products[i]
		log.Info(product.Name)
		prod := stores.Product{
			Name: product.Name,
			Ean:  product.Ean,
		}
		for _, store := range stores.Stores {
			price := store.GetPrice(prod)
			if price != nil {
				log.Infof("\t%s: %f", store.Name, *price)
			} else {
				log.Infof("\t%s: (not found)", store.Name)
			}
			product.LastPrice[store.Name] = price
			time.Sleep(500 * time.Millisecond)
		}
	}
	log.Info("Done refreshing prices")
}

func registerGauge(product *TrackedProduct, storeName string) {
	name := fmt.Sprintf(`product_price{name="%s", ean="%s", store_name="%s"}`, product.Name, product.Ean, storeName)
	metrics.NewGauge(name, func() float64 {
		price := product.LastPrice[storeName]
		if price == nil {
			return -1
		}
		return *price
	})
}
