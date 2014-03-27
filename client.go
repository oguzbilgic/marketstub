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
}

func New() *Client {
	client := &Client{}

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
	return nil
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

			tradePrice := tradePrice.Add(fpd.NewFromFloat(rand.Float64(), -2))
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
