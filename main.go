package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
  if len(os.Args) != 3 {
    printUsage()
    os.Exit(1)
  }
  if l, err := NewLane(os.Args[1], os.Args[2]); err != nil {
    log.Fatalln(err)
  } else {
    fmt.Println(l.String())
  }
}

func printUsage() {
  fmt.Println("Provide exactly two arguments: start and destination.")
  fmt.Println("Enclose in double quotes if either contains spaces.")
  fmt.Println("\nExample:\n\t$ toll \"Cleveland, OH\" \"Houston, TX\"")
  fmt.Println()
}
