package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	help bool
	version bool
	conf string
)

func parseArgs() {
	flag.BoolVar(&help, "h", false, "show this help")
	flag.BoolVar(&version, "version", false, "show version")
	flag.StringVar(&conf, "c", "/etc/goinetd.conf", "config file path")
	flag.Parse()
}

func init() {
	parseArgs()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	if version {
		fmt.Println("0.1.0")
		os.Exit(0)
	}
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	configs, err := ParseConfig(conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	var proxies []TcpProxy

	for _, c := range configs {
		p := TcpProxy{c, nil}
		proxies = append(proxies, p)
	}

	for i, _ := range proxies {
		err := proxies[i].ListenAndServe()
		if err != nil {
			fmt.Printf("Failed to establish tcp forwarding: %v %s\n", proxies[i].Config, err)
		}
	}

	fmt.Println("awaiting signal")
	<-done

	for _, proxy := range proxies {
		proxy.Stop()
	}
	fmt.Println("exiting")
}
