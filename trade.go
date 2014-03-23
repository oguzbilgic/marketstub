package marketstub

import (
	"github.com/oguzbilgic/fpd"
	"github.com/oguzbilgic/market"
	"time"
)

func stubTrade() *market.Trade {
	return &market.Trade{
		Symbol:   "market:stubUSD",
		Currency: market.USD,
		Volume:   fpd.New(1324, -2),
		Price:    fpd.New(13024, -2),
		Time:     time.Now(),
	}
}
