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

	entry := 127
	e, err := st.GetEntry(entry)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("got entry %d: %q\n", entry, e)

	// NOTE: something wrong with the delte
	e, err = st.DeleteData()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("got entry (and deleted): %q\n", e)

	entry = 322
	o, err := st.GetEntryOffset(entry)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("got entry %d's offset: %d\n", entry, o)

	ent, err := st.GetAllEntries()
	if err != nil {
		log.Panic(err)
	}
	for i, e := range ent {
		fmt.Printf("%d: entry at offset %d\n", i, e)
	}

	//writeData(st, 500)

	err = st.Close()
	if err != nil {
		log.Panic(err)
	}
}

func writeData(st *data.Store, n int) {
	for i := 0; i < n; i++ {
		d := fmt.Sprintf("{'id':%d,'name':'record number %d','active':true}\n", i, i)
		err := st.WriteData([]byte(d))
		if err != nil {
			log.Printf("error: %s\n", err)
			break
		}
	}
}

func readData(st *data.Store, n int) {
	for i := 0; i < n; i++ {
		b, err := st.ReadData()
		if err != nil {
			log.Printf("error: %s\n", err)
			break
		}
		log.Printf("read: %q\n", b)
	}
}
