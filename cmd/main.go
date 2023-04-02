package main

import (
	"context"
	"log"
	"time"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
)

func main() {
	cfg, err := config.NewFromYaml("../config.yaml", false)
	if err != nil {
		log.Fatalf("unable to load config from .yaml file: %v\n", err)
		return
	}

	server, err := machinery.NewServer(cfg)
	if server == nil || err != nil {
		log.Fatalf("unable to create server: %v\n", err)
		return
	}

	err = server.RegisterTasks(map[string]interface{}{
		"sum":  Sum,
		"mult": Mult,
		"sub":  Sub,
	})
	if err != nil {
		log.Fatalf("unable to register tasks: %v\n", err)
	}

	worker := server.NewWorker("tag", 1)
	go func() {
		err = worker.Launch()
		if err != nil {
			log.Fatalf("failed to launch worker: %v\n", err)
			return
		}
	}()

	eta := time.Now().Add(2 * time.Second)
	signature1 := &tasks.Signature{
		Name: "sum",
		Args: []tasks.Arg{
			{
				Type:  "[]int64",
				Value: []int64{1, 2, 3, 4},
			},
		},
		ETA:          &eta,
		RetryCount:   3,
		RetryTimeout: 100,
	}

	signature2 := &tasks.Signature{
		Name: "sub",
		Args: []tasks.Arg{
			{
				Type:  "[]int64",
				Value: []int64{1, 2, 3, 4},
			},
		},
		RetryCount:   3,
		RetryTimeout: 100,
	}

	eta = time.Now().Add(4 * time.Second)
	signature3 := &tasks.Signature{
		Name: "mult",
		Args: []tasks.Arg{
			{
				Type:  "[]int64",
				Value: []int64{1, 2, 3, 4},
			},
		},
		ETA:          &eta,
		RetryCount:   3,
		RetryTimeout: 100,
	}

	group, err := tasks.NewGroup(signature1, signature2, signature3)
	if err != nil {
		log.Fatalf("failed to create group: %v\n", err)
		return
	}

	asyncResults, err := server.SendGroupWithContext(context.Background(), group, 1)
	if err != nil {
		log.Printf("failed to send tasks %v\n", err)
		return
	}

	for _, asyncResult := range asyncResults {
		res, err := asyncResult.Get(1)
		if err != nil {
			log.Printf("failed to get task's result: %v\n", err)
			continue
		}
		log.Printf(
			"%v  %v\n",
			asyncResult.Signature.Args[0].Value,
			tasks.HumanReadableResults(res),
		)
	}
	log.Println("tasks finished")
}

func Sum(args []int64) (int64, error) {
	var sum int64
	for _, arg := range args {
		sum += arg
	}
	return sum, nil
}

func Mult(args []int64) (int64, error) {
	var res int64
	for _, arg := range args {
		res *= arg
	}
	return res, nil
}

func Sub(args []int64) (int64, error) {
	res, err := Sum(args)
	return -res, err
}
