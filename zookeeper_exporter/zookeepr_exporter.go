package main

import (
	"flag"
	"net/http"
	"log"
	"github.com/prometheus/client_golang/prometheus"
		"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	from = flag.String("from", "localhost:80", "The host of Zookeeper Server ")
	port = flag.String("port", "8080", "The port of \"/metrics\"  output endpoint")
	path = flag.String("path", "/metrics", "Path of output endpoint")
)

func init() {
	log.Printf("initalizing")
}

func main() {
	flag.Parse()
	prometheus.MustRegister()
	http.Handle(*path, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
				<head><title>Hadoop Exporter</title></head>
           	<body>
           		<h1>Zookeeper Exporter</h1>
				<p><a href='` + *path + `'>Metrics</a></p>
           	</body>
			</html>`))
	})
	listenAddress := ":" + *port
	log.Printf("server listing at %v", ":8080")
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
