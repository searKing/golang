// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashring_test

import (
	"fmt"
	"log"
	"slices"

	"github.com/searKing/golang/go/exp/container/hashring"
)

func ExampleNew() {
	c := hashring.New[string]()
	c.AddNodes("NodeA")
	c.AddNodes("NodeB")
	c.AddNodes("NodeC")
	var nodes []string
	for node := range c.All() {
		nodes = append(nodes, node)
	}
	slices.Sort(nodes)
	fmt.Printf("all nodes: %v\n", nodes)

	users := []string{"Alice", "Bob  ", "Eve  ", "Carol", "Dave "}
	fmt.Printf("locate nodes...\n")
	for _, u := range users {
		server, has := c.Get(u)
		if !has {
			log.Fatal()
		}
		fmt.Printf("	%s => %s\n", u, server)
	}

	// Output:
	// all nodes: [NodeA NodeB NodeC]
	// locate nodes...
	//	Alice => NodeB
	//	Bob   => NodeA
	//	Eve   => NodeA
	//	Carol => NodeC
	//	Dave  => NodeA
}

func ExampleHashRing_AddNodes() {
	c := hashring.New[string]()
	c.AddNodes("NodeA")
	c.AddNodes("NodeB")
	c.AddNodes("NodeC")
	users := []string{"Alice", "Bob  ", "Eve  ", "Carol", "Dave "}
	fmt.Println("initial state [A, B, C]")
	for _, u := range users {
		server, has := c.Get(u)
		if !has {
			log.Fatal()
		}
		fmt.Printf("%s => %s\n", u, server)
	}
	c.AddNodes("NodeD")
	c.AddNodes("NodeE")
	fmt.Println("\nwith NodeD, NodeE [A, B, C, D, E]")
	for _, u := range users {
		server, has := c.Get(u)
		if !has {
			log.Fatal()
		}
		fmt.Printf("%s => %s\n", u, server)
	}
	// Output:
	// initial state [A, B, C]
	// Alice => NodeB
	// Bob   => NodeA
	// Eve   => NodeA
	// Carol => NodeC
	// Dave  => NodeA
	//
	// with NodeD, NodeE [A, B, C, D, E]
	// Alice => NodeB
	// Bob   => NodeA
	// Eve   => NodeA
	// Carol => NodeE
	// Dave  => NodeA
}

func ExampleHashRing_RemoveNodes() {
	c := hashring.New[string]()
	c.AddNodes("NodeA")
	c.AddNodes("NodeB")
	c.AddNodes("NodeC")
	//users := []string{"Alice", "Bob", "Eve", "Carol", "Dave", "Isaac", "Justin", "Mallory", "Oscar", "Pat", "Victor", "Trent", "Walter"}
	users := []string{"Alice", "Bob  ", "Eve  ", "Carol", "Dave "}
	fmt.Println("initial state [A, B, C]")
	for _, u := range users {
		server, has := c.Get(u)
		if !has {
			log.Fatal()
		}
		fmt.Printf("%s => %s\n", u, server)
	}
	c.RemoveNodes("NodeA")
	fmt.Println("\nNodeA removed [B, C]")
	for _, u := range users {
		server, has := c.Get(u)
		if !has {
			log.Fatal()
		}
		fmt.Printf("%s => %s\n", u, server)
	}
	// Output:
	// initial state [A, B, C]
	// Alice => NodeB
	// Bob   => NodeA
	// Eve   => NodeA
	// Carol => NodeC
	// Dave  => NodeA
	//
	// NodeA removed [B, C]
	// Alice => NodeB
	// Bob   => NodeC
	// Eve   => NodeB
	// Carol => NodeC
	// Dave  => NodeB
}
