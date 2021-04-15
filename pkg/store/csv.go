package store

type csvStore struct {
}

/*
func (c *csvStore) Save(path string, v interface{}) error {
	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		return fmt.Errorf("failed to open file: %v\n", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	return nil
}
*/
func (c *csvStore) Load(path string, v interface{}) error {

	return nil
}
