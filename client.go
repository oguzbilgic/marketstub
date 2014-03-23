package marketstub

import (
	"errors"
	"github.com/oguzbilgic/market"
	"time"
)

type Client struct {
}

func New() *Client {
	return &Client{}
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

	go func() {
		for {
			time.Sleep(2 * time.Second)
			tradeChan <- stubTrade()
		}
	}()

	return tradeChan
}
