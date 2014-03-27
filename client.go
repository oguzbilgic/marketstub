package marketstub

import (
	"errors"
	"github.com/oguzbilgic/fpd"
	"github.com/oguzbilgic/market"
	"math/rand"
	"time"
)

type Client struct {
	tradeChans []chan *market.Trade
	tickChans  []chan *market.Tick
}

func New() *Client {
	client := &Client{}

	ticks := tickEngine(2*time.Second, client.NewTradeChan())
	go func() {
		for {
			tick := <-ticks
			for _, tickChan := range client.tickChans {
				tickChan <- tick
			}
		}
	}()

	trades := tradeEngine()
	go func() {
		for {
			trade := <-trades
			for _, tradeChan := range client.tradeChans {
				tradeChan <- trade
			}
		}
	}()

	return client
}

func (c *Client) Ticker() (*market.Tick, error) {
	return nil, errors.New("not implemented")
}

func (c *Client) OrderBook() ([]*market.Depth, error) {
	return nil, errors.New("not implemented")
}

func (c *Client) NewTickChan() chan *market.Tick {
	tickChan := make(chan *market.Tick)
	c.tickChans = append(c.tickChans, tickChan)
	return tickChan
}

func (c *Client) NewDepthChan() chan *market.Depth {
	return nil
}

func (c *Client) NewTradeChan() chan *market.Trade {
	tradeChan := make(chan *market.Trade)
	c.tradeChans = append(c.tradeChans, tradeChan)
	return tradeChan
}

func tradeEngine() chan *market.Trade {
	rand.Seed(time.Now().UnixNano())
	tradePrice := fpd.NewFromFloat(rand.Float64()*100, -3)
	tradeChan := make(chan *market.Trade)

	go func() {
		for {
			duration := time.Duration(rand.Float32() * 2000)
			time.Sleep(duration * time.Millisecond)

			priceChange := rand.Float64() * rand.Float64()

			if rand.Intn(10) < 5 {
				tradePrice = tradePrice.Add(fpd.NewFromFloat(priceChange, -3))
			} else {
				tradePrice = tradePrice.Sub(fpd.NewFromFloat(priceChange, -3))
			}

			tradeVolume := fpd.NewFromFloat(rand.Float64()*10, -2)

			trade := &market.Trade{
				Symbol:   "market:stubUSD",
				Currency: market.USD,
				Volume:   tradeVolume,
				Price:    tradePrice,
				Time:     time.Now(),
			}

			tradeChan <- trade
		}
	}()

	return tradeChan
}

func tickEngine(duration time.Duration, tradeChan chan *market.Trade) chan *market.Tick {
	tickChan := make(chan *market.Tick)
	tickDuration := time.Tick(duration)

	go func() {
		var tick *market.Tick

		for {
			select {
			case trade := <-tradeChan:
				if tick == nil {
					tick = &market.Tick{
						Symbol:   "market:stubUSD",
						Currency: market.USD,
						Time:     time.Now(),
						Volume:   fpd.New(0, -3),
						High:     trade.Price,
						Low:      trade.Price,
						Last:     trade.Price,
					}
				}

				if trade.Price.Cmp(tick.High) == 1 {
					tick.High = trade.Price
				}

				if trade.Price.Cmp(tick.Low) == -1 {
					tick.Low = trade.Price
				}

				tick.Volume = tick.Volume.Add(trade.Volume)
			case <-tickDuration:
				tickChan <- tick
				tick = nil
			}
		}

	}()

	return tickChan
}
