package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	useragent        = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36"
	accept           = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
	accept_encoding  = "gzip, deflate"
	accept_language  = "en-US,en;q=0.8"
	cache_control    = "max-age=0"
	connection       = "keep-alive"
	dnt              = "1"
	upgrade_insecure = "1"
)

func main() {
	client := &http.Client{}
	for _, url := range os.Args[1:] {
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", useragent)
		req.Header.Set("Accept", accept)
		req.Header.Set("Accept-Encoding", accept_encoding)
		req.Header.Set("Accept-Language", accept_language)
		req.Header.Set("Cache-Control", cache_control)
		req.Header.Set("Connection", connection)
		req.Header.Set("DNT", dnt)
		req.Header.Set("Upgrade-Insecure-Requests", upgrade_insecure)

		resp, err := client.Do(req)

		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to fetch: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		for k, v := range resp.Header {
			fmt.Fprintf(os.Stderr, "%s: %s\n", k, v[0])
		}
		responseReader := decompress(resp)
		defer responseReader.Close()

		content, err := ioutil.ReadAll(responseReader)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to read: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s", content)
	}
}

func decompress(resp *http.Response) io.ReadCloser {
	var responseReader io.ReadCloser
	var err error
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		responseReader, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to decompress response: %v\n", err)
			os.Exit(1)
		}

	default:
		responseReader = resp.Body
	}
	return responseReader
}
