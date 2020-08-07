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
