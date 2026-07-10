package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"ascii-art-web/art_web"
	"ascii-art-web/handlers"
)

func main() {
	if len(os.Args) != 1 {
		fmt.Println("Wrong Number of Arguments")
		return
	}

	start := time.Now()

	if err := web.EnsureBanners(); err != nil {
		log.Fatalf("Banner integrity check failed: %v", err)
	}
	fmt.Println("✅ Banner files are present and valid.")

	elapsed := time.Since(start)
	fmt.Printf("🚀 Server Finished Loading, time taken: %s\n", elapsed)

	handlers.Routing()
}
