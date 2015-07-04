package price

import (
	"net/http"

	"hd.qiniu.com/services/price"
)

func PriceService(host string, tr http.RoundTripper) price.PriceService {
	return price.NewPriceService(host, tr)
}
