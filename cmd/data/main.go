package main

import (
	"fmt"
	"github.com/scottcagno/net-tools/pkg/data"
	"log"
)

func main() {

	st, err := data.OpenStore("cmd/data/test/data.txt")
	if err != nil {
		log.Panic(err)
	}

	for i := 0; i < 250; i++ {
		// skip the odd records
		if i%2 != 0 {
			err := st.NextData()
			if err != nil {
				fmt.Printf("error: %s\n", err)
				break
			}
			continue
		}
		b, err := st.ReadData()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			break
		}
		fmt.Printf("%s\n", b)
	}

	err = st.Close()
	if err != nil {
		log.Panic(err)
	}
}
