// Package gosnowth contains an IRONdb client library written in Go.
package gosnowth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const numericTestData = "[[1380000000,50],[1380000300,60]]"

const numericTestAllData = `[
    [
        1379998800,
        {
            "count": 60,
            "value": 10,
            "stddev": 0,
            "derivative": 0,
            "derivative_stddev": 0,
            "counter": 0,
            "counter_stddev": 0,
            "derivative2": 0,
            "derivative2_stddev": 0,
            "counter2": 0,
            "counter2_stddev": 0
        }
    ],
    [
        1380002400,
        {
            "count": 60,
            "value": 10,
            "stddev": 0,
            "derivative": 0,
            "derivative_stddev": 0,
            "counter": 0,
            "counter_stddev": 0,
            "derivative2": 0,
            "derivative2_stddev": 0,
            "counter2": 0,
            "counter2_stddev": 0
        }
    ],
    [
        1380006000,
        {
            "count": 60,
            "value": 10,
            "stddev": 1,
            "derivative": 1,
            "derivative_stddev": 1,
            "counter": 1,
            "counter_stddev": 1,
            "derivative2": 1,
            "derivative2_stddev": 1,
            "counter2": 1,
            "counter2_stddev": 1
        }
    ]
]`

const numericTestWriteData = `[
	{
		"count": 1,
		"value": 10,
		"derivative": 1,
		"counter": 1,
		"stddev": 1,
		"derivative_stddev": 1,
		"counter_stddev": 1,
		"metric": "test",
		"id": "fc85e0ab-f568-45e6-86ee-d7443be8277d",
		"offset": 0,
		"parts": {
			"period": 1,
			"data": [
				{
					"count": 1,
					"value": 1,
					"derivative": 1,
					"counter": 1,
					"stddev": 1,
					"derivative_stddev": 1,
					"counter_stddev": 1
				}
			]
		}
	}
]`

func TestNumericValue(t *testing.T) {
	nv := NumericValueResponse{}
	if err := json.Unmarshal([]byte(numericTestData), &nv); err != nil {
		t.Error("error decoding JSON: ", err)
	}

	if nv.Data[0].Time != time.Unix(1380000000, 0) {
		t.Error("invalid time parsing")
	}

	if nv.Data[1].Time != time.Unix(1380000300, 0) {
		t.Error("invalid time parsing")
	}

	if nv.Data[0].Value != 50 {
		t.Error("invalid value parsing")
	}

	if nv.Data[1].Value != 60 {
		t.Error("invalid value parsing")
	}
}

func TestNumericAllValue(t *testing.T) {
	nv := NumericAllValueResponse{}
	if err := json.Unmarshal([]byte(numericTestAllData), &nv); err != nil {
		t.Error("error decoding JSON: ", err)
	}
}

func TestNumericReadWrite(t *testing.T) {
	ms := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		if r.RequestURI == "/state" {
			_, _ = w.Write([]byte(stateTestData))
			return
		}

		if r.RequestURI == "/stats.json" {
			_, _ = w.Write([]byte(statsTestData))
			return
		}

		u := "/read/1529509020/1529509200/1/" +
			"fc85e0ab-f568-45e6-86ee-d7443be8277d/count/test"
		if strings.HasPrefix(r.RequestURI, u) {
			_, _ = w.Write([]byte(numericTestData))
			return
		}

		u = "/read/1529509020/1529509200/1/" +
			"fc85e0ab-f568-45e6-86ee-d7443be8277d/all/test"
		if strings.HasPrefix(r.RequestURI, u) {
			_, _ = w.Write([]byte(numericTestAllData))
			return
		}

		u = "/write/nnt"
		if strings.HasPrefix(r.RequestURI, u) {
			w.WriteHeader(200)
			return
		}
	}))

	defer ms.Close()
	sc, err := NewSnowthClient(false, ms.URL)
	if err != nil {
		t.Fatal("Unable to create snowth client", err)
	}

	u, err := url.Parse(ms.URL)
	if err != nil {
		t.Fatal("Invalid test URL")
	}

	node := &SnowthNode{url: u}
	res, err := sc.ReadNumericValues(time.Unix(1529509020, 0),
		time.Unix(1529509200, 0), 1, "count",
		"fc85e0ab-f568-45e6-86ee-d7443be8277d", "test", node)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 2 {
		t.Fatalf("Expected results: 2, got: %v", len(res))
	}

	if res[0].Value != 50 {
		t.Errorf("Expected value: 50, got: %v", res[0].Value)
	}

	resA, err := sc.ReadNumericAllValues(time.Unix(1529509020, 0),
		time.Unix(1529509200, 0), 1, "fc85e0ab-f568-45e6-86ee-d7443be8277d",
		"test", node)
	if err != nil {
		t.Fatal(err)
	}

	if len(resA) != 3 {
		t.Fatalf("Expected results: 3, got: %v", len(resA))
	}

	if resA[0].Value != 10 {
		t.Errorf("Expected value: 10, got: %v", resA[0].Value)
	}

	nv := []NumericWrite{}
	err = json.NewDecoder(bytes.NewBufferString(numericTestWriteData)).Decode(&nv)
	if err != nil {
		t.Fatal(err)
	}

	err = sc.WriteNumeric(nv, node)
	if err != nil {
		t.Fatal(err)
	}
}
