package pac

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/axelrindle/proxyguy/config"
	"github.com/darren/gpac"
	"github.com/sirupsen/logrus"
	"gobn.github.io/coalesce"
)

type Pac struct {
	Logger *logrus.Logger
	Config *config.Structure

	script *string
}

func (p Pac) CheckConnectivity() bool {
	theUrl, err := url.Parse(p.Config.PacUrl)
	if err != nil {
		p.Logger.Fatalln(err)
	}

	hostWithoutPort := strings.Split(theUrl.Host, ":")[0]
	isIpAddress := net.ParseIP(hostWithoutPort) != nil
	if !isIpAddress {
		timeout := time.Duration(p.Config.Timeout) * time.Millisecond
		_, err = p.lookupHost(hostWithoutPort, timeout)
	} else {
		var conn net.Conn
		conn, err = net.Dial("tcp", theUrl.Host)
		if conn != nil {
			conn.Close()
		}
	}

	if err != nil {
		return false
	}

	return true
}

func (p *Pac) LoadPacScript() error {
	p.Logger.Debugln("Loading pac script from " + p.Config.PacUrl + "…")

	resp, err := http.Get(p.Config.PacUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	script := strings.ToValidUTF8(string(body), "")
	p.script = &script

	return nil
}

func (p Pac) DetermineProxies(determineOverride *string) []string {
	if len(p.Config.Proxy.Override) != 0 {
		return []string{p.Config.Proxy.Override}
	}

	pac, err := gpac.New(*p.script)
	if err != nil {
		p.Logger.Fatalln(err)
	}

	p.Logger.Debugln("Determining proxy endpoint…")

	url := coalesce.String(determineOverride, &p.Config.Proxy.Determine)
	proxies, err := pac.FindProxyForURL(*url)
	if err != nil {
		p.Logger.Fatalln(err)
	}

	return strings.Split(proxies, ";")
}

func TrimProxy(part string) string {
	return strings.Trim(strings.Replace(part, "PROXY", "", 1), " ")
}

func (p Pac) lookupHost(hostname string, timeout time.Duration) ([]string, error) {
	c1 := make(chan []string)
	c2 := make(chan error)

	var ipaddr []string
	var err error

	go func() {
		var ipaddr []string
		ipaddr, err := net.LookupHost(hostname)
		if err != nil {
			c2 <- err
		}

		c1 <- ipaddr
	}()

	select {
	case ipaddr = <-c1:
	case err = <-c2:
	case <-time.After(timeout):
		return ipaddr, errors.New("timeout")
	}

	if err != nil {
		return nil, err
	}

	return ipaddr, nil
}
