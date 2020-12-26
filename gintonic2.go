package main

import (
    // "github.com/rivo/tview"
    "fmt"
    "flag"
    "path"
    "os"
    "log"
    "bufio"
    "unicode"
)

var launch_db = make(map[string]string)

var shuffle bool
var help bool
var launch_db_filepath string

var inputs []string

func readConf() {
    if len(launch_db_filepath) == 0 {
        home := os.Getenv("HOME")
        launch_db_filepath = path.Join(home, ".gintonic2/launch_db")
    }
    f, err := os.Open(launch_db_filepath)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    scanner :=bufio.NewScanner(f)
    var key string
    var value string
    for scanner.Scan() {
        line := []rune(scanner.Text())
        if len(line) > 0 {
            if unicode.IsSpace(line[0]) {
                value = string(line)
                if len(key)>0 {
                    launch_db[key] = value
                }
            } else {
                key = string(line)
            }
        }
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

func readArgs() {
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage:\n\n")
        fmt.Fprintf(os.Stderr, "%s [flags] paths     Launch for specified paths\n\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "%s [flags]           Launch for paths from stdin\n\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "Flags:\n")
        flag.PrintDefaults()
    }
    flag.BoolVar(&shuffle, "s", false, "shuffle targets")
    flag.StringVar(&launch_db_filepath, "c", "", "path to a custom launch_db file")
    flag.BoolVar(&shuffle, "d", false, "print debug messages")
    flag.BoolVar(&help, "h", false, "show this message and exit")
    flag.Parse()
    if help {
        flag.Usage()
        os.Exit(0)
    }
    if flag.NArg() > 0 {
        inputs = flag.Args()
    } else {
        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            line := scanner.Text()
            if len(line)>0 {
                inputs = append(inputs, line)
            }
        }
        if err := scanner.Err(); err != nil {
                log.Println(err)
        }
    }
    fmt.Println(inputs)
}

func main() {
    readArgs()
    readConf()
}
