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
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	localIP    string
	localPort  uint16
	remoteIP   string
	remotePort uint16
}

func (c *Config) Verify() bool {
	return net.ParseIP(c.localIP) != nil && net.ParseIP(c.remoteIP) != nil
}

func ParseConfig(file string) ([]*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	var result []*Config
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		splited := strings.Fields(text)
		if len(splited) != 4 {
			return nil, errors.New(fmt.Sprintf("config error: %s\n", text))
		}
		sport, err := strconv.Atoi(splited[1])
		if err != nil {
			return nil, err
		}
		dport, err := strconv.Atoi(splited[3])
		if err != nil {
			return nil, err
		}
		cfg := &Config{splited[0], uint16(sport), splited[2], uint16(dport)}
		if cfg.Verify() {
			result = append(result, cfg)
		} else {
			return nil, errors.New(fmt.Sprintf("config error: %s\n", text))
		}
	}
	return result, nil
}
