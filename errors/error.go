package errors

import (
	"fmt"
	"os"
)

func PrintErrors(es ...error) {
	for _, e := range es {
		fmt.Fprintf(os.Stderr, "%s\n", e.Error())
	}
}