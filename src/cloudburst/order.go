package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"github.com/satori/go.uuid"
	"fmt"
)

func order(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "POST" {
		createOrder(w,r)
	} else if r.Method == "PUT" {
		updateOrder(w,r)
	} else if r.Method == "GET" {
		getOrder(w,r)
	}
}

func createOrder(w http.ResponseWriter, r *http.Request){
	enableCors(&w)

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var order Order
	err = json.Unmarshal(b, &order)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if order.Id == "" {
		uuid, _ := uuid.NewV4()
		order.Id = uuid.String()
	}

	if order.UserId == "" {
		http.Error(w, "User ID is not sent", 500)
		return
	}

	if order.RestaurantId == 0 {
		http.Error(w, "Restaurant ID is not sent", 500)
		return
	}

	order.OrderStatus = "Order Placed"

	if debug { fmt.Printf("%+v\n", order) }

	output, err := json.Marshal(order)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if insertObjects("orders", order.Id, output, getCluster(order.UserId)) == nil {
		//update orderlist
		updateOrderList(order.UserId, order.Id)

		//delete cart
		err := deleteObjects("cart", order.UserId, getCluster(order.UserId))
		if err != nil {
			log.Println("[RIAK DEBUG] " + err.Error())
		}

		w.WriteHeader(http.StatusOK)
		w.Write(output)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func updateOrder(w http.ResponseWriter, r *http.Request){
	enableCors(&w)

	var userid, orderid string
	orderid = r.URL.Query().Get("orderid")
	orderid = r.URL.Query().Get("userid")

	if orderid != "" {
		resp, err := queryObjects("orders", orderid, getCluster(userid))
		if err != nil {
			log.Println("[RIAK DEBUG] " + err.Error())
		}

		var order Order
		err = json.Unmarshal(resp.Values[0].Value, &order)
		order.OrderStatus = "Order Processed"

		output, err := json.Marshal(order)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		newrsp, err := updateObjects("orders", orderid, []byte(output), getCluster(userid))
		if err != nil {
			log.Println("[RIAK DEBUG] " + err.Error())
		}else {
			if len(newrsp.Values) > 0 {
				w.Write(newrsp.Values[0].Value)
			}
		}
	}
}

func updateOrderList(userid string, orderid string){
	if orderid != "" {
		resp, err := queryObjects("orderlist", userid, getCluster(userid))
		if err != nil {
			log.Println("[RIAK DEBUG] " + err.Error())
		}

		var orders []string

		if resp.Values != nil {
			err = json.Unmarshal(resp.Values[0].Value, &orders)
			if err != nil {
				log.Println("updateOrderList: json unmarshalling error")
			}
		}

		orders = append(orders, orderid)

		output, err := json.Marshal(orders)
		if err != nil {
			log.Println("updateorderlist: json marshal error"+ err.Error())
		}

		if resp.Values != nil {
			_, err = updateObjects("orderlist", userid, []byte(output), getCluster(userid))
		} else {
			err = insertObjects("orderlist", userid, []byte(output), getCluster(userid))
		}

		if err != nil {
			log.Println("[RIAK DEBUG] " + err.Error())
		}
	}
}

func getOrder(w http.ResponseWriter, r *http.Request){
	enableCors(&w)

	var userid, orderid string
	orderid = r.URL.Query().Get("orderid")
	orderid = r.URL.Query().Get("userid")

	if orderid != "" {
		resp, err := queryObjects("orders", orderid, getCluster(userid))
		if err != nil {
			log.Println("[RIAK DEBUG] " + err.Error())
		} else {
			if len(resp.Values) > 0 {
				w.Write(resp.Values[0].Value)
			}
		}
	}
}

func getOrders(w http.ResponseWriter, r *http.Request){
enableCors(&w)
	if r.Method == "GET" {
		var userid string
		userid = r.URL.Query().Get("userid")

		resp, err := queryObjects("orderlist", userid, getCluster(userid))
		if err != nil {
			log.Println("[RIAK DEBUG] " + err.Error())
			return
		}

		var orderids []string
		err = json.Unmarshal(resp.Values[0].Value, &orderids)
		if err != nil {
			log.Println("getOrders: json unmarshalling error")
			return
		}

		var orders []Order
		var order Order
		for _, orderid := range orderids {
			if orderid != "" {
				resp, err := queryObjects("orders", orderid, getCluster(userid))
				if err != nil {
					log.Println("[RIAK DEBUG] " + err.Error())
				}

				err = json.Unmarshal(resp.Values[0].Value, &order)
				if err != nil {
					log.Println("getOrders unmarshalling error: " + err.Error())
					return
				}
				orders = append(orders, order)
			}
		}

		output, err := json.Marshal(orders)
		if err != nil {
			log.Println("getOrders marshalling error: " + err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(output)
	}
}