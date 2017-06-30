package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	logger       = "weardex.requests"
	zoneID       = "visearch-production-ap-northeast-1"
	search       = "api_method=search"
	uploadsearch = "api_method=uploadsearch"
)

type Request struct {
	t        time.Time
	method   string
	username string
	param    string
}

type Requests []*Request

func (slice Requests) Len() int {
	return len(slice)
}

func (slice Requests) Less(i, j int) bool {
	return slice[i].t.Before(slice[j].t)
}

func (slice Requests) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func parseFiles(filenames []string) (Requests, error) {
	var responses Requests

	for _, filename := range filenames {
		response, err := parseFile(filename)
		if err != nil {
			return responses, errors.Wrap(err, "parseFile fail")
		}
		responses = append(responses, response...)
	}

	sort.Sort(responses)

	return responses, nil
}

func parseFile(filename string) (Requests, error) {
	var response Requests
	file, err := os.Open(filename)
	if err != nil {
		return response, errors.Wrapf(err, "Open file failed,filename:%v", filename)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var line string
	for err == nil {
		line, err = reader.ReadString('\n')
		req, yes := parseLine(line)
		if yes {
			response = append(response, req)
		}
	}

	if err != io.EOF {
		return response, errors.Wrapf(err, "bufio reader error,filename:%v", filename)
	}

	sort.Sort(response)

	return response, nil
}

func parseLine(line string) (*Request, bool) {
	if !strings.Contains(line, logger) || !strings.Contains(line, zoneID) || (!strings.Contains(line, search) && !strings.Contains(line, uploadsearch)) {
		return nil, false
	}

	parts := strings.Split(line, " ")
	if len(parts) <= 12 {
		log.Printf("log format error,line:%v\n", line)
		return nil, false
	}
	t, err := time.Parse("2006-01-02T15:04:05.000Z", parts[2])
	if err != nil {
		log.Printf("time parse error,origin:%v err:%v\n", parts[2], err)
		return nil, false
	}
	method := strings.Split(parts[9], "=")[1]
	username := strings.Split(parts[11], "=")[1]
	param := parts[12][7 : len(parts[12])-1]
	// fmt.Println(t, method, username, param)

	if username != "Interpark@fashion_live" || method != "search" {
		return nil, false
	}

	return &Request{
		t:        t,
		method:   method,
		username: username,
		param:    param,
	}, true
}
