package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
	"encoding/json"

  vegeta "github.com/tsenart/vegeta/v12/lib"
) 

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
			latency := runTest().String()
			resp := make(map[string]string)
			resp["latency"] = latency
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
	})

	log.Println("Listening on localhost:80")

	log.Fatal(http.ListenAndServe(":80", nil))
}


func runTest() time.Duration {
  rate := vegeta.Rate{Freq: 100, Per: time.Second}
  duration := 4 * time.Second
  targeter := vegeta.NewStaticTargeter(vegeta.Target{
    Method: "GET",
    URL:    "http://google.it/",
  })
  attacker := vegeta.NewAttacker()

  var metrics vegeta.Metrics
  for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
    metrics.Add(res)
  }
  metrics.Close()

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	return metrics.Latencies.P99
}
