package main

import (
	"fmt"

	"github.com/Brandon-Butterbaugh/gator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("Brandon")
	cfg = config.Read()
	fmt.Println(cfg)
}
