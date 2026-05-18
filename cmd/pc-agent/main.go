package main

// TODO: implement pc-agent command
import (
	"fmt"
	"os"
	"pc-app/internal/agent"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		fmt.Printf("[ENV] failed to load .env: %v\n", err)
	}

	agent.Run()
}
