package main

import "sync"

func Serial(url string, fetcher Fetcher, fetched map[string]bool) {

}

type fetchState struct {
	mu      sync.Mutex
	fetched map[string]bool
}

type Fetcher interface {
	Fetch(url string) (urls []string, err error)
}
