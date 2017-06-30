package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
)

func searchCall(req *Request) error {
	url := fmt.Sprintf("%s/%s?%s", *host, req.method, req.param)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Wrap(err, "NewRequest fail")
	}

	access, secret, err := getCred(req.username)
	if err != nil {
		return errors.Wrap(err, "getCred fail")
	}
	httpReq.SetBasicAuth(access, secret)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return errors.Wrap(err, "Do http request fail")
	}
	resp.Body.Close()

	return nil
}

func playback(reqs Requests) {
	id := 0

	c := time.Tick(time.Millisecond * time.Duration(1000 / *qps))

	var wg sync.WaitGroup

	for range c {

		if id >= len(reqs) {
			break
		}
		wg.Add(1)

		go func(req *Request) {
			defer wg.Done()

			err := searchCall(req)
			if err != nil {
				log.Printf("searchCall failed,err:%v\n", err)
			} else {
				log.Printf("searchCall ok,req:%v\n", req)
			}

		}(reqs[id])

		id++
	}

	wg.Wait()
}
