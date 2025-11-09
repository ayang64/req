package client

import (
	"cmp"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type Payload struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Cur   int    `json:"cur"`
	Limit int    `json:"limit"`
}

func (p Payload) More() bool {
	return p.Cur <= p.Limit
}

func TestEvery(t *testing.T) {
	p := Payload{
		Name:  "foo",
		Value: "bar",
	}

	cur, limit := 1, rand.Intn(20)+10
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// t.Logf("got url: %s", r.URL)
		// t.Logf("got headers %#v", r.Header)
		p.Cur, p.Limit = cur, limit
		if err := json.NewEncoder(w).Encode(&p); err != nil {
			t.Logf("wtf?")
		}
		cur++
	}))
	t.Cleanup(func() { svr.Close() })

	h := http.Header{
		"SOME-HEADER":    []string{"some-value"},
		"ANOTHER-HEADER": []string{"another-value"},
	}
	v := url.Values{
		"q":    []string{cmp.Or("", "go tips and tricks")},
		"sort": []string{cmp.Or("", "time")},
	}

	e, err := Every[Payload](t.Context(), http.MethodGet, svr.URL, h, v, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("got: %#v", e)
}

func TestAll(t *testing.T) {
	p := Payload{
		Name:  "foo",
		Value: "bar",
	}

	cur, limit := 1, rand.Intn(20)+10
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p.Cur, p.Limit = cur, limit
		if err := json.NewEncoder(w).Encode(&p); err != nil {
			t.Logf("wtf?")
		}
		cur++
	}))
	t.Cleanup(func() { svr.Close() })

	h := http.Header{
		"SOME-HEADER":    []string{"some-value"},
		"ANOTHER-HEADER": []string{"another-value"},
	}
	v := url.Values{
		"q":    []string{cmp.Or("", "go tips and tricks")},
		"sort": []string{cmp.Or("", "time")},
	}

	for v, err := range All[Payload](t.Context(), http.MethodGet, svr.URL, h, v, nil, nil) {
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("got: %#v", v)
	}
}

func TestGet(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := Payload{
			Name:  "foo",
			Value: "bar",
		}

		t.Logf("got url: %s", r.URL)
		t.Logf("got headers %#v", r.Header)

		if err := json.NewEncoder(w).Encode(&p); err != nil {
			t.Logf("wtf?")
		}
	}))
	t.Cleanup(func() { svr.Close() })

	h := http.Header{
		"SOME-HEADER":    []string{"some-value"},
		"ANOTHER-HEADER": []string{"another-value"},
	}
	v := url.Values{
		"q":    []string{cmp.Or("", "go tips and tricks")},
		"sort": []string{cmp.Or("", "time")},
	}

	p, err := Send[Payload](t.Context(), http.MethodGet, svr.URL, h, v, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v", p)
}
