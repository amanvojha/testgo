package main

import (
	"net/http"
	"log"
)

func getRestaurants(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "GET" {
		pincode := r.URL.Query().Get("pincode")

		if pincode != "" {
			resp, err := queryObjects("restaurants", pincode, cluster1)
			if err != nil {
				log.Println("[RIAK DEBUG] " + err.Error())
			}else {
				if len(resp.Values) > 0 {
					w.Write(resp.Values[0].Value)
				}
			}
		}
	}
}

func getMenu(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "GET" {
		resp, err := queryObjects("restaurants", "menu", cluster1)
		if err != nil {
			log.Println("[RIAK DEBUG] " + err.Error())
		}else {
			if len(resp.Values) > 0 {
				w.Write(resp.Values[0].Value)
			}
		}
	}
}