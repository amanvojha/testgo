package main

import (
	"net/http"
	"fmt"
	"log"
	"github.com/basho/riak-go-client"
	"os"
)

var debug = true

//connect to tcp ports for cluster
var s1 = os.Getenv("RIAK1_N1")
var s2 = os.Getenv("RIAK1_N2")
var s3 = os.Getenv("RIAK1_N3")
var s4 = os.Getenv("RIAK1_N4")
var s5 = os.Getenv("RIAK1_N5")

var t1 = os.Getenv("RIAK2_N1")
var t2 = os.Getenv("RIAK2_N2")
var t3 = os.Getenv("RIAK2_N3")
var t4 = os.Getenv("RIAK2_N4")
var t5 = os.Getenv("RIAK2_N5")

var cluster1 *riak.Cluster
var cluster2 *riak.Cluster

func initCluster1(){
	nodeOpts1 := &riak.NodeOptions{
		RemoteAddress: s1,
	}

	nodeOpts2 := &riak.NodeOptions{
		RemoteAddress: s2,
	}

	nodeOpts3 := &riak.NodeOptions{
		RemoteAddress: s3,
	}

	nodeOpts4 := &riak.NodeOptions{
		RemoteAddress: s4,
	}

	nodeOpts5 := &riak.NodeOptions{
		RemoteAddress: s5,
	}

	var node1, node2, node3, node4, node5 *riak.Node
	var err error

	if node1, err = riak.NewNode(nodeOpts1); err != nil {
		fmt.Println(err.Error())
	}

	if node2, err = riak.NewNode(nodeOpts2); err != nil {
		fmt.Println(err.Error())
	}

	if node3, err = riak.NewNode(nodeOpts3); err != nil {
		fmt.Println(err.Error())
	}

	if node4, err = riak.NewNode(nodeOpts4); err != nil {
		fmt.Println(err.Error())
	}

	if node5, err = riak.NewNode(nodeOpts5); err != nil {
		fmt.Println(err.Error())
	}

	nodes := []*riak.Node{node1, node2, node3, node4, node5}
	opts := &riak.ClusterOptions{
		Nodes: nodes,
	}

	log.Println( nodes )

	cluster1, err = riak.NewCluster(opts)
	if err != nil {
		fmt.Println(err.Error())
	}

	if err := cluster1.Start(); err != nil {
		fmt.Println(err.Error())
	}
}

func initCluster2(){
	nodeOpts1 := &riak.NodeOptions{
		RemoteAddress: t1,
	}

	nodeOpts2 := &riak.NodeOptions{
		RemoteAddress: t2,
	}

	nodeOpts3 := &riak.NodeOptions{
		RemoteAddress: t3,
	}

	nodeOpts4 := &riak.NodeOptions{
		RemoteAddress: t4,
	}

	nodeOpts5 := &riak.NodeOptions{
		RemoteAddress: t5,
	}

	var node1, node2, node3, node4, node5 *riak.Node
	var err error

	if node1, err = riak.NewNode(nodeOpts1); err != nil {
		fmt.Println(err.Error())
	}

	if node2, err = riak.NewNode(nodeOpts2); err != nil {
		fmt.Println(err.Error())
	}

	if node3, err = riak.NewNode(nodeOpts3); err != nil {
		fmt.Println(err.Error())
	}

	if node4, err = riak.NewNode(nodeOpts4); err != nil {
		fmt.Println(err.Error())
	}

	if node5, err = riak.NewNode(nodeOpts5); err != nil {
		fmt.Println(err.Error())
	}

	nodes := []*riak.Node{node1, node2, node3, node4, node5}
	opts := &riak.ClusterOptions{
		Nodes: nodes,
	}

	log.Println( nodes )

	cluster2, err = riak.NewCluster(opts)
	if err != nil {
		fmt.Println(err.Error())
	}

	if err := cluster2.Start(); err != nil {
		fmt.Println(err.Error())
	}
}

func init() {
	initCluster1()
	initCluster2()
}

func handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	fmt.Fprintf(w, "Hi there! Welcome to goBurger")
}

func main() {
	http.HandleFunc("/hi", handler)
	http.HandleFunc("/getRestaurants", getRestaurants)
	http.HandleFunc("/getMenu", getMenu)
	http.HandleFunc("/cart", cart)
	http.HandleFunc("/order", order)
	http.HandleFunc("/orders", getOrders)
	http.HandleFunc("/user", user)

	http.ListenAndServe(":8080", nil)

	defer func() {
		if err := cluster1.Stop(); err != nil {
			log.Println(err.Error())
		}
	}()

	defer func() {
		if err := cluster2.Stop(); err != nil {
			log.Println(err.Error())
		}
	}()
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func insertObjects(bucket string, key string, body []byte, cluster *riak.Cluster) error {
	obj := &riak.Object{
		ContentType:     "application/json",
		Key:             key,
		Value:           body,
	}

	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucket(bucket).
		WithContent(obj).
		Build()
	if err != nil {
		return err
	}

	if err = cluster.Execute(cmd); err != nil {
		return err
	}

	return nil
}

func queryObjects(bucket string, key string, cluster *riak.Cluster ) (*riak.FetchValueResponse, error) {
	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucket(bucket).
		WithKey(key).
		Build()
	if err != nil {
		return nil, err
	}

	if err = cluster.Execute(cmd); err != nil {
		return nil, err
	}

	fvc := cmd.(*riak.FetchValueCommand)
	rsp := fvc.Response

	if debug {
		if rsp.Values != nil {
			log.Println(string(rsp.Values[0].Value))
		}
	}

	return rsp, nil
}

func updateObjects(bucket string, key string, newval []byte, cluster *riak.Cluster) (*riak.FetchValueResponse, error) {
	cmd, err := riak.NewFetchValueCommandBuilder().
		WithBucket(bucket).
		WithKey(key).
		Build()
	if err != nil {
		return nil, err
	}

	if err = cluster.Execute(cmd); err != nil {
		return nil, err
	}

	fvc := cmd.(*riak.FetchValueCommand)
	rsp := fvc.Response
	if rsp.Values !=  nil {
		obj := rsp.Values[0]

		if debug {
			if obj != nil {
				log.Println(string(obj.Value))
			}
		}

		obj.Value = newval

		cmd, err = riak.NewStoreValueCommandBuilder().
			WithBucket(bucket).
			WithKey(key).
			WithContent(obj).
			Build()

		if err != nil {
			return nil, err
		}

		if err = cluster.Execute(cmd); err != nil {
			return nil, err
		}

		return rsp, nil
	}
	return nil, nil
}

func deleteObjects(bucket string, key string, cluster *riak.Cluster) error{
	cmd, err := riak.NewDeleteValueCommandBuilder().
		WithBucket(bucket).
		WithKey(key).
		Build()
	if err != nil {
		return err
	}

	return cluster.Execute(cmd)
}


//shard the database based on userid
func getCluster(userid string) *riak.Cluster {
	if len(userid)%2 == 0 {
		return cluster1
	} else {
		return cluster2
	}
}
