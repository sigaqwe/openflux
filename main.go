package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/wlynxg/anet"
	"universal-bypass-tool/socks5"
	"universal-bypass-tool/transport"
	"universal-bypass-tool/transport/oneme"
	"universal-bypass-tool/transport/yandex"
	"universal-bypass-tool/tunnel"
	"universal-bypass-tool/utils"
)

var (
	globalDocUrl string
	maxToken     string
	maxUid       string
)

func startClient(url string, socksAddr string, transportType string, token string, uid string) error {
	config := transport.DefaultConfig()

	var trans transport.Transport

	switch transportType {
	case "yandex":
		trans = yandex.NewYandexDocsTransport(url, config)

	case "oneme":
		uidint, err := strconv.ParseInt(uid, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid MAX user id: %w", err)
		}
		trans = oneme.NewOneMeTransport(false, token, uidint, config)

	default:
		return fmt.Errorf("unknown transport type: %s", transportType)
	}

	if err := trans.Start(); err != nil {
		return fmt.Errorf("failed to start transport: %w", err)
	}

	tun := tunnel.NewTCPTunnel(trans, false)

	log.Printf("Running as CLIENT (SOCKS5 on %s)", socksAddr)

	socks5Server := socks5.NewSOCKS5Server(socksAddr, tun)

	return socks5Server.Start()
}

//export RunMainClient
func RunMainClient(url *C.char) {
	docURL := C.GoString(url)

	go func() {
		if err := startClient(
			docURL,
			"127.0.0.1:1080",
			"yandex",
			"",
			"",
		); err != nil {
			log.Printf("OpenFlux client stopped: %v", err)
		}
	}()
}

//export RunMain
func RunMain() {
	docURL := globalDocUrl

	go func() {
		if err := startClient(
			docURL,
			"127.0.0.1:1080",
			"yandex",
			maxToken,
			maxUid,
		); err != nil {
			log.Printf("OpenFlux client stopped: %v", err)
		}
	}()
}

func main() {
	fmt.Print("written by p1neappleXpress\n")

	exitNode := flag.Bool("exit-node", false, "Run as exit node (needs root)")
	client := flag.Bool("client", false, "Run as client")
	debug := flag.Bool("debug", false, "Enable verbose debug logging")
	socksAddr := flag.String("socks5", ":1080", "SOCKS5 address")
	transportType := flag.String("transport", "yandex", "Transport type (yandex, oneme)")

	flag.StringVar(
		&globalDocUrl,
		"url",
		"http://#",
		"Document URL. If you use Yandex.Docs transport",
	)

	flag.StringVar(
		&maxToken,
		"maxToken",
		"",
		"MAX user token",
	)

	flag.StringVar(
		&maxUid,
		"maxUid",
		"",
		"MAX user ID",
	)

	flag.Parse()

	if !*exitNode && !*client {
		flag.Usage()
		os.Exit(1)
	}

	if *debug {
		utils.EnableDebug()
	}

	log.Printf("=== Universal Bypass Tool ===")
	log.Printf(
		"Mode: %s",
		map[bool]string{
			true:  "EXIT NODE",
			false: "CLIENT",
		}[*exitNode],
	)

	log.Printf("Transport: %s", *transportType)

	if *client {
		if err := startClient(
			globalDocUrl,
			*socksAddr,
			*transportType,
			maxToken,
			maxUid,
		); err != nil {
			log.Fatal(err)
		}
		return
	}

	config := transport.DefaultConfig()

	var trans transport.Transport

	switch *transportType {
	case "yandex":
		trans = yandex.NewYandexDocsTransport(globalDocUrl, config)

	case "oneme":
		uidint, _ := strconv.ParseInt(maxUid, 10, 64)
		trans = oneme.NewOneMeTransport(
			*exitNode,
			maxToken,
			uidint,
			config,
		)

	default:
		log.Fatalf("Unknown transport type: %s", *transportType)
	}

	if err := trans.Start(); err != nil {
		log.Fatalf("Failed to start transport: %v", err)
	}

	_ = tunnel.NewTCPTunnel(trans, *exitNode)

	select {}
}
