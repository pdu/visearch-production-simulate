package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFile(t *testing.T) {
	filename := "/Users/peng/Downloads/ViSearch-search_service.log.2017040401"
	responses, err := parseFile(filename)
	assert.Nil(t, err)
	for i := 1; i < len(responses); i++ {
		assert.Equal(t, true, (responses[i-1].t.Before(responses[i].t) || responses[i-1].t == responses[i].t))
	}
}

func TestParseLine(t *testing.T) {
	line0 := "10.4.1.134	/mnt/logs/search-service/search.log	ViSearch:search_service - 2017-04-04T00:59:59.048Z - INFO - weardex.requests - client_ip=211.174.54.41 api_method=search http_method=GET account_app_name=Interpark@fashion_live param=[score=true&limit=50&im_name=4790689431&page=1&fq=category:001262002&score_max=0.99] user_agent=visearch-sdk-java/1.2.2-SNAPSHOT Linux/2.6.32-573.18.1.el6.x86_64 Java HotSpot(TM) 64-Bit Server VM/25.73-b02/1.8.0_73/ko_null zone_id=\"visearch-production-ap-northeast-1\" account=\"Interpark\" - 616331554674744689"
	_, yes := parseLine(line0)
	assert.Equal(t, true, yes)
	line2 := "10.4.1.134	/mnt/logs/search-service/search.log	ViSearch:search_service - 2017-04-04T00:59:59.048Z - INFO - weardex.request - client_ip=211.174.54.41 api_method=search http_method=GET account_app_name=Interpark@fashion_live param=[score=true&limit=50&im_name=4790689431&page=1&fq=category:001262002&score_max=0.99] user_agent=visearch-sdk-java/1.2.2-SNAPSHOT Linux/2.6.32-573.18.1.el6.x86_64 Java HotSpot(TM) 64-Bit Server VM/25.73-b02/1.8.0_73/ko_null zone_id=\"visearch-production-ap-northeast-1\" account=\"Interpark\" - 616331554674744689"
	_, yes = parseLine(line2)
	assert.NotEqual(t, true, yes)
	line3 := "10.4.1.134	/mnt/logs/search-service/search.log	ViSearch:search_service - 2017-04-04T00:59:59.048Z - INFO - weardex.requests - client_ip=211.174.54.41 api_method=search http_method=GET account_app_name=Interpark@fashion_live param=[score=true&limit=50&im_name=4790689431&page=1&fq=category:001262002&score_max=0.99] user_agent=visearch-sdk-java/1.2.2-SNAPSHOT Linux/2.6.32-573.18.1.el6.x86_64 Java HotSpot(TM) 64-Bit Server VM/25.73-b02/1.8.0_73/ko_null zone_id=\"visearch-production-ap-northeast-\" account=\"Interpark\" - 616331554674744689"
	_, yes = parseLine(line3)
	assert.NotEqual(t, true, yes)
}
