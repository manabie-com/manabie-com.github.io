package examples

import (
	"fmt"

	"github.com/google/go-jsonnet"
)

func MakeConfig() error {
	vm := jsonnet.MakeVM()
	jsonData, err := vm.EvaluateFile(`dashboard.jsonnet`)
	if err != nil {
		return err
	}
	fmt.Println(jsonData)
	return nil
}
