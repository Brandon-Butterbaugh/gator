package main

import (
	"github.com/Brandon-Butterbaugh/gator.git/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("Brandon")
	cfg = config.Read()
	println(cfg)
}
