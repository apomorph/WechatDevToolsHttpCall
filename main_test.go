package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"testing"
)

func TestHttp(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:10001/login?storeId=10000")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}

func TestConcurrency(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 100000; i++ {
		s := strconv.Itoa(i)
		wg.Add(1)
		go func(storeId string) {
			defer wg.Add(-1)
			resp, err := http.Get("http://127.0.0.1:10001/login?storeId=" + storeId)
			if err != nil {
				// fmt.Println(err)
				return
			}
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// fmt.Println(err)
				return
			}
			fmt.Println(string(data))
		}(s)
	}

	wg.Wait()
}
