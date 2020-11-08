package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/daheige/go-proj/app/rpc/service"
	"github.com/daheige/go-proj/pb"
	"github.com/golang/protobuf/proto"
)

func main() {
	httpMux := http.NewServeMux() //创建一个http ServeMux实例
	httpMux.HandleFunc("/debug/pprof/", pprof.Index)
	httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	httpMux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		//接收json格式，然后转换为pb格式，调用rpc的逻辑
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		defer r.Body.Close()

		req := &pb.HelloReq{}
		FormatMessage(r.Header.Get("Content-Type"), body, req)

		log.Println("recv request: ", req)

		log.Println("pb message:")
		log.Println(proto.Marshal(req))
		log.Println(proto.MarshalTextString(req))

		//调用service
		s := &service.GreeterService{}
		res, err := s.SayHello(r.Context(), req)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		log.Println("res: ", res)

		b, err := json.Marshal(res)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		var b2 = []byte{10, 7, 100, 97, 104, 101, 105, 103, 101}

		log.Println("b2 string: ", string(b2))

		w.Write(b)
	})

	log.Println("server has run")
	if err := http.ListenAndServe(":8080", httpMux); err != nil {
		log.Println(err)
	}
}

// FormatMessage 解析消息题到req中
func FormatMessage(cType string, body []byte, req proto.Message) (proto.Message, error) {
	unmarshalFunc := proto.Unmarshal
	// 兼容JSON请求
	if strings.Contains(cType, "application/json") {
		unmarshalFunc = func(b []byte, m proto.Message) error {
			return json.Unmarshal(b, m)
		}
	}

	if err := unmarshalFunc(body, req); err != nil {
		log.Println("error:", err)

		return nil, err
	}

	return req, nil
}
