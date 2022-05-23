package build

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/gopacket"
)

const (
	Port    = 80
	Version = "0.1"
)

const (
	CmdPort = "-p"
)

type H struct {
	port    int
	version string
}

var hp *H

func NewInstance() *H {
	if hp == nil {
		hp = &H{
			port:    Port,
			version: Version,
		}
	}
	return hp
}

func (m *H) ResolveStream(net, transport gopacket.Flow, buf io.Reader) {
	bio := bufio.NewReader(buf)
	for {
		// log.Println(net, transport)
		if strconv.Itoa(m.port) == transport.Dst().String() {
			decodeRequest(bio)
		} else {
			decodeResponse(bio)
		}
	}
}

func decodeRequest(bio *bufio.Reader) {
	req, err := http.ReadRequest(bio)
	if err == io.EOF {
		log.Println("decode request err", err)
		return
	} else if err != nil {
		log.Println("decode request err", err)
		return
	} else {
		var msg = "request ["
		msg += req.Method
		msg += "] ["
		msg += req.Host + req.URL.String()
		msg += "] ["
		msg += fmt.Sprintf("%+v", req.Header)
		msg += "] ["

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println("read request body err", err)
			return
		}
		// req.ParseForm()
		// msg += req.Form.Encode()
		msg += string(body)
		msg += "]"

		log.Println(msg)

		req.Body.Close()
	}
}

func decodeResponse(bio *bufio.Reader) {
	rsp, err := http.ReadResponse(bio, nil)
	if err == io.EOF {
		// log.Println("decode response err", err)
		return
	} else if err != nil {
		// log.Println("decode response err", err)
		return
	} else {
		var msg = "response ["
		msg += fmt.Sprintf("%+v", rsp.Header)
		msg += "] ["
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			// log.Println("read response body err", err)
			return
		}
		msg += string(body)
		msg += "]"

		// log.Println(msg)

		rsp.Body.Close()
	}
}

func (m *H) BPFFilter() string {
	return "tcp and port " + strconv.Itoa(m.port)
}

func (m *H) Version() string {
	return Version
}

func (m *H) SetFlag(flg []string) {

	c := len(flg)

	if c == 0 {
		return
	}
	if c>>1 == 0 {
		fmt.Println("ERR : Http Number of parameters")
		os.Exit(1)
	}
	for i := 0; i < c; i = i + 2 {
		key := flg[i]
		val := flg[i+1]

		switch key {
		case CmdPort:
			port, err := strconv.Atoi(val)
			m.port = port
			if err != nil {
				panic("ERR : port")
			}
			if port < 0 || port > 65535 {
				panic("ERR : port(0-65535)")
			}
			break
		default:
			panic("ERR : mysql's params")
		}
	}
}
