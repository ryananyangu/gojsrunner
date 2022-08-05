package utils

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"strings"
	"time"
)

func ReadFile(file string) (string, error) {
	openedFile, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer openedFile.Close()

	byteValue, _ := ioutil.ReadAll(openedFile)

	return string(byteValue[:]), nil

}

func TimeTrack(start time.Time, name string) (string, time.Duration) {
	elapsed := time.Since(start)
	return name, elapsed
}

func Request(request string, headers map[string][]string, urlPath string, method string) (string, error) {

	reqURL, _ := url.Parse(urlPath)

	reqBody := ioutil.NopCloser(strings.NewReader(request))

	req := &http.Request{
		Method: method,
		URL:    reqURL,
		Header: headers,
		Body:   reqBody,
	}

	// res, err := http.DefaultClient.Do(req)

	res, err := ExternalRequestTimer(req)
	if err != nil {
		// log http error
		Log.Errorf("SEND REQUEST | URL : %s | METHOD : %s | BODY : %s | ERROR : %v", urlPath, method, request, err)
		return "", err
	}

	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	Log.Infof("SEND REQUEST | URL : %s | METHOD : %s | BODY : %s | STATUS : %s | HTTP_CODE : %d", urlPath, method, request, res.Status, res.StatusCode)

	if res.StatusCode > 299 || res.StatusCode <= 199 {
		Log.Errorf("SEND REQUEST | URL : %s | METHOD : %s | BODY : %s | STATUS : %s | HTTP_CODE : %d", urlPath, method, request, res.Status, res.StatusCode)
		return res.Status, fmt.Errorf("%d | %s", res.StatusCode, data)
	}

	resbody := string(data)
	Log.Infof("SEND REQUEST | URL : %s | METHOD : %s | BODY : %s | STATUS : %s | HTTP_CODE : %d", urlPath, method, resbody, res.Status, res.StatusCode)

	return resbody, nil
}

func ExternalRequestTimer(req *http.Request) (*http.Response, error) {

	var start, connect, dns, tlsHandshake time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			Log.Infof("DNS Done: %v\n", time.Since(dns))
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			Log.Infof("TLS Handshake: %v\n", time.Since(tlsHandshake))
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			Log.Infof("Connect time: %v\n", time.Since(connect))
		},

		GotFirstResponseByte: func() {
			Log.Infof("Time from start to first byte: %v\n", time.Since(start))
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return res, err
	}
	return res, nil
}

func init() {
	if os.Getenv("AMQP_SERVER_URL") == "" {
		Log.Fatal("Amq url not loaded")
	}
}
