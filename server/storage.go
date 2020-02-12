package main

import "io/ioutil"

func saveDump(data []byte) error {
	return ioutil.WriteFile(dumpfile, data, 0600)
}

func loadDump() ([]byte, error) {
	return ioutil.ReadFile(dumpfile)
}
