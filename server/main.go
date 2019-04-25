/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"

	pb "proto_dice"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// PlayDice implements helloworld.GreeterServer
func (s *server) PlayDice(ctx context.Context, in *pb.DiceRequest) (*pb.DiceReply, error) {
	ReceivedDice, _ := strconv.Atoi(in.Dice)
	SentDice := rand.Intn(6-1) + 1
	Condition := ""
	if ReceivedDice > SentDice {
		Condition = "Won"
	} else if SentDice > ReceivedDice {
		Condition = "Lose"
	} else {
		Condition = "Draw"
	}
	fmt.Println("Client Dice : ", in.Dice)
	fmt.Println("Our Dice : ", SentDice)
	fmt.Println("Result for client : ", Condition)
	return &pb.DiceReply{
		ClientDice: strconv.Itoa(ReceivedDice),
		ServerDice: strconv.Itoa(SentDice),
		Message:    Condition,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPlayServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
