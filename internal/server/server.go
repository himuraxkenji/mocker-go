package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type MockServer struct {
	Port       string         `json:"port,omitempty"`
	Endpoints  []EndpointMock `json:"endpoints,omitempty"`
	ServerName string         `json:"server_name,omitempty"`
}

type EndpointMock struct {
	Endpoint string `json:"endpoint,omitempty"`
	Method   string `json:"method,omitempty"`
	Response struct {
		Status int `json:"status,omitempty"`
		Body   any `json:"body,omitempty"`
	} `json:"response,omitempty"`
}

func InitServer() {
	servers := readResponseFiles()
	var wg sync.WaitGroup

	for _, s := range servers {
		wg.Add(1)
		go startServer(s, &wg)
	}

	wg.Wait()
}

func startServer(s MockServer, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := http.NewServeMux()

	for _, r := range s.Endpoints {
		route := fmt.Sprintf("%s %s/", strings.ToUpper(r.Method), r.Endpoint)
		fmt.Println(route)
		mux.HandleFunc(route, func(w http.ResponseWriter, rq *http.Request) {
			if rq.Method != r.Method {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(r.Response.Status)
			json.NewEncoder(w).Encode(r.Response.Body)
		})
	}
	log.Printf("Server %s running on port %s...\n", s.ServerName, s.Port)
	log.Fatal(http.ListenAndServe(s.Port, mux))
}

func readResponseFiles() []MockServer {
	d := "/home/himura/code/go/tools/mocker-go/internal/responses/"

	entries, err := os.ReadDir(d)
	var result []MockServer

	if err != nil {
		fmt.Println("Error when read directory -", err)
		return result
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		fn := e.Name()

		if strings.HasSuffix(fn, ".json") {
			fp := filepath.Join(d, fn)
			c, err := os.ReadFile(fp)

			if err != nil {
				fmt.Println("Error on read file json", err)
				continue
			}
			var r MockServer
			err = json.Unmarshal(c, &r)
			if err != nil {
				fmt.Printf("a error has occur when decode - %v\n", err)
			}
			result = append(result, r)
		}
	}

	fmt.Println("read all files in directory - ")

	return result

}
