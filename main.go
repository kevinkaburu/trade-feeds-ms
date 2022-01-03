package main

import (
	"log"
	"os"
	TradesMs "trades/src"

	"github.com/joho/godotenv"
)

func main() {
	_ = os.Setenv("KE", "Africa/Nairobi")
	if err := godotenv.Load(); err != nil {
		log.Printf("unable to read dotenv %v", err)
	}
	TradesMs.Run()
}
