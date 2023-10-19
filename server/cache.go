package server

import (
	"errors"
	"net/url"

	"github.com/axelrindle/proxyguy/pac"
)

type Cache struct {
	Pac *pac.Pac

	proxies map[string][]string
}

func (c *Cache) Init() {
	c.proxies = make(map[string][]string)
}

func (c *Cache) Update() error {
	if c.Pac.CheckConnectivity() {
		err := c.Pac.LoadPacScript()
		if err != nil {
			return err
		}
	}

	return errors.New("update failed")
}

func (c *Cache) FindProxies(u string) []string {
	url, err := url.Parse(u)
	if err != nil {
		return nil
	}

	result, found := c.proxies[url.Host]

	if !found {
		result = c.Pac.DetermineProxies(&u)
		c.proxies[url.Host] = result
	}

	return result
}
