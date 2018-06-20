// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package stubs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func StartHTTPTestServer() (srv *HTTPTestServer, addr string) {
	srv = &HTTPTestServer{}
	addr = srv.Start()

	return srv, addr
}

type httpRequestHandler func(reqData map[string]interface{}, rw http.ResponseWriter, r *http.Request)

type HTTPTestServer struct {
	reqHandler httpRequestHandler
	srv        *httptest.Server
}

func (s *HTTPTestServer) Start() (srvAddr string) {
	if s.reqHandler == nil {
		s.reqHandler = defaultHTTPHandler
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		r.Body.Close()

		if err != nil {
			panic(err)
		}

		var data map[string]interface{}
		if len(body) != 0 {
			data = make(map[string]interface{})
			err = json.Unmarshal(body, &data)

			if err != nil {
				panic(err)
			}
		}

		s.reqHandler(data, w, r)
	})
	srv := httptest.NewServer(h)

	s.srv = srv

	return srv.URL
}

func (s *HTTPTestServer) Close() {
	s.srv.Close()
}

func (s *HTTPTestServer) SetRequestHandler(h httpRequestHandler) {
	s.reqHandler = h
}

func defaultHTTPHandler(reqData map[string]interface{}, rw http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(reqData)

	if err != nil {
		panic(err)
	}

	rw.WriteHeader(200)
	rw.Write(res)
}
