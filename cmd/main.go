package main

import (
    "fmt"
    "os"
)

var version = "0.1.0"

func main() {
    if len(os.Args) < 2 {
        printUsage()
        return
    }

    switch os.Args[1] {
    case "help", "-h", "--help":
        printUsage()
    case "version", "-v", "--version":
        fmt.Println(version)
    default:
        fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
        printUsage()
        os.Exit(2)
    }
}

func printUsage() {
    fmt.Println("someday - a minimal CLI tool")
    fmt.Println()
    fmt.Println("Usage:")
    fmt.Println("  someday <command>")
    fmt.Println()
    fmt.Println("Commands:")
    fmt.Println("  help       Show help")
    fmt.Println("  version    Print version")
}

