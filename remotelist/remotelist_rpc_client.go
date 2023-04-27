package main

import (
	"fmt"
	"net/rpc"
)

type appendArgs struct {
	ListIndex int
	Value     int
}

func main() {
	client, err := rpc.Dial("tcp", ":5000")
	if err != nil {
		fmt.Println("Something went wrong while dialing: %w.\n", err)
		return
	}

	for {
		fmt.Print("Please write:\n")
		fmt.Print("A: to add a number in a list;\n")
		fmt.Print("G: to get the numbers of a list;\n")
		fmt.Print("R: to remove the last number in a list;\n")
		fmt.Print("S: to get the size of a list;\n")

		var operation string
		fmt.Scan(&operation)

		if operation == "A" {
			var number int
			var listIndex int
			fmt.Print("Now please write the identification of the list that you want to add: ")
			fmt.Scan(&listIndex)
			fmt.Print("Now please write a number to add to the remote list: ")
			fmt.Scan(&number)

			var reply bool
			err = client.Call("PersistentRemoteList.Append", appendArgs{listIndex, number}, &reply)
			if err != nil {
				fmt.Printf("Something went wrong while adding the number: %w\n", err)
				continue
			}
			fmt.Printf("Successfully added number %d.\n", number)
		}
		if operation == "G" {
			var listIndex int
			fmt.Print("Now please write the identification of the list that you want get the numbers: ")
			fmt.Scan(&listIndex)
			var reply []int
			err = client.Call("PersistentRemoteList.Get", listIndex, &reply)
			if err != nil {
				fmt.Printf("Something went wrong while getting the list: %w\n", err)
				continue
			}
			fmt.Printf("The numbers on list %d are %v.\n", listIndex, reply)
		}

		if operation == "R" {
			var listIndex int
			fmt.Print("Now please write the identification of the list that you want to remove the last number: ")
			fmt.Scan(&listIndex)
			var reply int
			err = client.Call("PersistentRemoteList.Remove", listIndex, &reply)
			if err != nil {
				fmt.Printf("Something went wrong while removing the last number: %w\n", err)
				continue
			}
			fmt.Printf("Successfully removed last number: %d.\n", reply)
		}

		if operation == "S" {
			var listIndex int
			fmt.Print("Now please write the identification of the list that you want get the size: ")
			fmt.Scan(&listIndex)
			var reply int
			err = client.Call("PersistentRemoteList.Size", listIndex, &reply)
			if err != nil {
				fmt.Printf("Something went wrong while getting the size of the list: %w\n", err)
				continue
			}
			fmt.Printf("The size of list %d is %d.\n", listIndex, reply)
		}

	}
}
