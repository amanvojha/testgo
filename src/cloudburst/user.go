package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func user(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "POST" {
		createUser(w,r)
	} else if r.Method == "PUT" {
		updateUser(w,r)
	} else if r.Method == "GET" {
		getUser(w,r)
	} else if r.Method == "DELETE" {
		deleteUser(w,r)
	}
}

func createUser(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var user User
	err = json.Unmarshal(b, &user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Println(user)

	if user.Id == "" {
		http.Error(w, "User ID is not sent", 500)
		return
	}

	if user.Password == "" {
		http.Error(w, "Password is not sent", 500)
		return
	}

	output, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if insertObjects("users", user.Id, output, getCluster(user.Id)) == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var user User
	err = json.Unmarshal(b, &user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if user.Id != "" {
		_, err := queryObjects("users", user.Id, getCluster(user.Id))
		if err != nil {
			log.Println(user.Id + "not present in RIAK")
			log.Println(err.Error())
			return
		}

		output, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		newrsp, err := updateObjects("users", user.Id, []byte(output), getCluster(user.Id))
		if err != nil {
			log.Println("[RIAK DEBUG] " + err.Error())
		} else {
			if len(newrsp.Values) > 0 {
				w.Write(newrsp.Values[0].Value)
			}
		}
	}
}

func getUser(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	var userid string
	userid = r.URL.Query().Get("id")

	if userid != "" {
		resp, err := queryObjects("users", userid, getCluster(userid))
		if err != nil {
			log.Println("[RIAK DEBUG] " + err.Error())
		} else {
			if resp.Values != nil {
				w.Write(resp.Values[0].Value)
			}
		}
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	var userid string
	userid = r.Header.Get("id")

	err := deleteObjects("users", userid, getCluster(userid))
	err = deleteObjects("orderlist", userid, getCluster(userid))
	if err != nil {
		log.Println("[RIAK DEBUG] " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}