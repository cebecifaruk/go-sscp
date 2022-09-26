package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	sscp "github.com/cebecifaruk/go-sscp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide host address as an argument")
		os.Exit(1)
	}

	c, err := sscp.NewPLCConnection(os.Args[1], 1, true)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Connected")

	res, err := c.Login("admin", "rw", "", 10240)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", res)

	vars := []*sscp.Variable{
		&sscp.Variable{
			Uid:    9980,
			Offset: 0,
			Length: 4,
			Value:  []byte{},
		},
		&sscp.Variable{
			Uid:    9982,
			Offset: 0,
			Length: 4,
			Value:  []byte{},
		},
	}

	fmt.Println(c.ReadVariablesDirectly(vars, nil))

	for _, v := range vars {
		var value float32
		binary.Read(bytes.NewReader(v.Value), binary.BigEndian, &value)
		fmt.Printf("%+v %+v\n", v, value)
	}

	// c.Logout()
}
