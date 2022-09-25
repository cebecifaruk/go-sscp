# SSCP Driver for Go

![Test and Build Workflow](https://github.com/cebecifaruk/go-sscp/actions/workflows/test.yaml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/cebecifaruk/go-sscp.svg)](https://pkg.go.dev/github.com/cebecifaruk/go-sscp)

SSCP protocol implmentation for Go.

This implementation heavily depends on the specification presented in the doc folder of the project root.
SSCP protocol is simply a protocol for communicating with PLCs. Now, it is possible to read/write variables
with go language.

## Implementation

This implementation covers the features listed below. Please open an issue and check the repository regularly
if you need a feature not implemented.


| Feature                 | Implemented  | Tested  |
|-------------------------|--------------|---------|
| GetBasicInfo            | ❌           | ❌       |
| Login                   | ✅           | ✅       |
| Logout                  | ✅           | ✅       |
| LargeBinarySend         | ❌           | ❌       |
| LargeBinaryRecv         | ❌           | ❌       |
| GetPLCStatistics        | ✅           | ✅       |
| GetTaskStatistics       | ✅           | ✅       |
| GetChannelStatistics    | ✅           | ✅       |
| ReadVariablesDirectly   | ✅           | ✅       |
| WriteVariablesDirectly  | ✅           | ✅       |
| ReadVariablesFileMode   | ❌           | ❌       |
| WriteVariablesFileMode  | ❌           | ❌       |
| TimeSetup               | ✅           | ❌       |
| TimeSetupExtended       | ✅           | ❌       |



## Example Usage


```go
	conn, err := sscp.NewPLCConnecetion(os.Args[1], 1, true)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	res, err := c.Login("admin", "rw", "", 10240)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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

	err = c.ReadVariablesDirectly(vars)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, v := range vars {
		var value float32
		binary.Read(bytes.NewReader(v.Value), binary.BigEndian, &value)
		fmt.Printf("%+v %+v\n", v, value)
	}
```