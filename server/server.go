package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/axelrindle/proxyguy/config"
	"github.com/axelrindle/proxyguy/logger"
	"github.com/axelrindle/proxyguy/pac"
	"github.com/sirupsen/logrus"
	"gopkg.in/elazarl/goproxy.v1"
)

var log *logrus.Entry = logger.ForComponent("server")

type Server struct {
	proxy *goproxy.ProxyHttpServer
	p     *pac.Pac
	cache *Cache
}

func (s *Server) Start() {
	log.Logger.ExitFunc = func(code int) {
		// ignore
	}

	s.p = &pac.Pac{}

	s.cache = &Cache{Pac: s.p}
	s.cache.Init()
	s.cache.Update()

	s.proxy = goproxy.NewProxyHttpServer()
	s.proxy.Verbose = log.Logger.IsLevelEnabled(logrus.DebugLevel)

	s.proxy.OnRequest().DoFunc(s.http)
	s.proxy.ConnectDial = s.connectDial

	cfg := config.Data()
	address := fmt.Sprintf("%s:%v", cfg.Server.Address, cfg.Server.Port)
	log.Println("Starting proxy server on " + address + "…")

	server := http.Server{
		Addr:    address,
		Handler: s.proxy,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Infoln("Server closed.")
			log.Fatalln("Error starting proxy server: ", err)
		}
	}()

	<-GracefulShutdown(context.Background(), 5*time.Second, log.Logger, map[string]ShutdownHook{
		"Server": func(ctx context.Context) error {
			server.Shutdown(ctx)
			return nil
		},
	})

	os.Exit(0)
}

func (s Server) connectDial(network, addr string) (net.Conn, error) {
	if s.p.CheckConnectivity() {
		targets := s.cache.FindProxies(addr)

		var conn net.Conn
		var err error

		for _, target := range targets {
			if target != "DIRECT" && target != "" {
				stripped := pac.TrimProxy(target)
				log.Traceln("Trying \"" + stripped + "\"")

				fun := s.proxy.NewConnectDialToProxy("http://" + stripped)
				if fun == nil {
					continue
				}

				conn, err = fun(network, addr)

				if err != nil {
					continue
				}

				if conn != nil {
					return conn, nil
				}
			}
		}

		if err != nil {
			return nil, err
		}

		return nil, errors.New("unknown error")
	}

	return net.Dial(network, addr)
}

func (s Server) http(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	if s.p.CheckConnectivity() {
		ctx.RoundTripper = goproxy.RoundTripperFunc(s.roundTrip)
	} else {
		log.Debugln("Proxy inactive.")
	}

	return req, nil
}

func (s Server) roundTrip(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Response, error) {
	targetUrl := req.URL.String()
	targets := s.p.DetermineProxies(&targetUrl)

	if len(targets) == 0 {
		return s.proxy.Tr.RoundTrip(req)
	}

	var resp *http.Response
	var err error

	for _, target := range targets {
		var transport *http.Transport = http.DefaultTransport.(*http.Transport)

		if target != "DIRECT" && target != "" {
			stripped := strings.Trim(strings.Replace(target, "PROXY", "", 1), " ")
			log.Traceln("Trying \"" + stripped + "\"")
			parsedUrl, err := url.Parse("http://" + stripped)
			if err != nil {
				log.Debugln(err)
				continue
			}

			transport.Proxy = http.ProxyURL(parsedUrl)
		}

		resp, err = transport.RoundTrip(req)
		if err != nil {
			continue
		}

		if resp != nil {
			return resp, nil
		}
	}

	if err != nil {
		return nil, err
	}

	return nil, errors.New("unknown error")
}
