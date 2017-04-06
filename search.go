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
	_, err = client.Do(httpReq)
	if err != nil {
		return errors.Wrap(err, "Do http request fail")
	}

	return nil
}

func playback(reqs Requests) {
	id := 0
	now := time.Now()
	c := time.Tick(time.Millisecond)
	var wg sync.WaitGroup

	for current := range c {

		past := current.Sub(now)

		for id < len(reqs) && reqs[id].t.Sub(reqs[0].t) <= past {
			wg.Add(1)

			go func(req *Request) {
				wg.Done()

				err := searchCall(req)
				if err != nil {
					log.Printf("searchCall failed,err:%v\n", err)
				}

			}(reqs[id])

			id++
		}

		if id >= len(reqs) {
			break
		}
	}

	wg.Wait()
}
