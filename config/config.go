package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type configuration struct {
	PrivateKey string
	PublicKey  string
}

var Info = configuration{}

func init() {
	print("ESTO SE ESTA EJECUTANDOOOOOOOOOOOOOO")
	file_abs_path, _ := filepath.Abs("config/conf.json")
	//file, err := os.Open("../config/conf.json")
	file, err := os.Open(file_abs_path)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Info)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", Info)
}
