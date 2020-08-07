package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

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

	configs, err := ParseConfig("goinetd.conf")
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
