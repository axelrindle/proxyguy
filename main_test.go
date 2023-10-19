package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/axelrindle/proxyguy/config"
	"github.com/sirupsen/logrus"
)

const pacScript = `
// From: https://en.wikipedia.org/w/index.php?title=Proxy_auto-config&oldid=1175505362#Example

function FindProxyForURL(url, host) {
    // our local URLs from the domains below example.com don't need a proxy:
    if (shExpMatch(host, '*.example.com')) {
        return 'DIRECT';
    }

    // URLs within this network are accessed through
    // port 8080 on fastproxy.example.com:
    if (isInNet(host, '10.0.0.0', '255.255.248.0')) {
        return 'PROXY fastproxy.example.com:8080';
    }

    // All other requests go through port 8080 of proxy.example.com.
    // should that fail to respond, go directly to the WWW:
    return 'PROXY proxy.example.com:8080; DIRECT';
}
`

var pacServer *httptest.Server

func startPacServer() {
	pacServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(pacScript))
	}))
}

func stopServer() {
	pacServer.Close()
}

func TestMain(m *testing.M) {
	startPacServer()
	code := m.Run()
	stopServer()
	os.Exit(code)
}

var logger *logrus.Logger = logrus.New()

func TestProcessing(t *testing.T) {
	cfg := &config.Structure{
		PacUrl: pacServer.URL,
		Proxy: config.StructureProxy{
			Determine: "https://github.com",
		},
	}

	url, _ := Process(logger, cfg)

	if url != "proxy.example.com:8080" {
		t.Fatal("Invalid proxy endpoint returned!")
	}
}
