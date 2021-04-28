package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {

	Exec("bash")
	out := Exec("ping", "-c5", "8.8.8.8")

	fmt.Println(out)

}

func Exec(name string, arg ...string) string {
	b, err := exec.Command(name, arg...).Output()
	if err != nil {
		log.Fatalf("executing command, error: %v\n", err)
	}
	return string(b)
}
