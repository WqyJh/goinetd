//MIT License
//Copyright (c) 2020 Qiying Wang
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	help    bool
	version bool
	conf    string
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
