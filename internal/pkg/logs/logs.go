package logs

import (
	"fmt"
	"os"
)

func Error(err error) {
	fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
}

func ErrorMsg(str string) {
	fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", str)
}
