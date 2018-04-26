package stubs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func StartHTTPTestServer() (srv *HTTPTestServer, addr string, closeSrv func()) {
	srv = &HTTPTestServer{}
	addr = srv.Start()

	return srv, addr, srv.Close
}

type httpRequestHandler func(reqData map[string]interface{}, rw http.ResponseWriter)

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

		data := make(map[string]interface{})
		err = json.Unmarshal(body, &data)

		if err != nil {
			panic(err)
		}

		s.reqHandler(data, w)
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

func defaultHTTPHandler(reqData map[string]interface{}, rw http.ResponseWriter) {
	res, err := json.Marshal(reqData)

	if err != nil {
		panic(err)
	}

	rw.WriteHeader(200)
	rw.Write(res)
}
