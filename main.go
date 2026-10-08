package main

import (
	"fmt"
	"gator/internal/config"
)

func main() {
	// Read configuration
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set user
	err = cfg.SetUser("tom")
	if err != nil {
		fmt.Println(err)
		return
	}

	//Re-read the configuration
	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(cfg)
}
