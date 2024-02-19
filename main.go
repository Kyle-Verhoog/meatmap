package main

import (
	"fmt"
	"github.com/namsral/flag"
	mux "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "static/index.html")
}

func serveCreate(w http.ResponseWriter, r *http.Request, hotel *Hotel) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	room, err := hotel.createRoom()
	if err != nil {
		log.Println("create failed!")
		http.Redirect(w, r, "/todo", 301)
		return
	}

	roomName := room.name
	http.Redirect(w, r, fmt.Sprintf("/r/%v", roomName), 301)
}

func roomHandler(w http.ResponseWriter, r *http.Request, hotel *Hotel) {
	http.ServeFile(w, r, "static/room.html")
}

func main() {
	var host, port string
	flag.StringVar(&host, "host", "", "http host")
	flag.StringVar(&port, "port", "8080", "http port")
	flag.Parse()

	tracer.Start(tracer.WithService("meatmap"))
	defer tracer.Stop()

	r := mux.NewRouter()
	hotel := newHotel()

	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		serveCreate(w, r, hotel)
	}).Methods("POST")

	// normal room handler
	r.HandleFunc("/r/{room}", func(w http.ResponseWriter, r *http.Request) {
		roomHandler(w, r, hotel)
	}).Methods("GET")

	// "secure" rooms
	r.HandleFunc("/x", func(w http.ResponseWriter, r *http.Request) {
		roomHandler(w, r, hotel)
	}).Methods("GET")
	r.HandleFunc("/ws/{room}", hotel.serveHotel)
	http.Handle("/", r)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("listening on %s:%s\n", host, port)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		os.Exit(0)
	}()

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
