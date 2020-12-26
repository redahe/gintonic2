package main

import (
    // "github.com/rivo/tview"
    "bufio"
    "flag"
    "fmt"
    "log"
    "os"
    "os/exec"
    "path"
    "path/filepath"
    "strings"
    "unicode"
)

var launch_db = make(map[string]string)

var shuffle bool
var help bool
var debug bool
var no_ui bool
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
                value = strings.TrimLeft(string(line), "\t ")
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
    if debug {
        log.Print("Launch db: ", launch_db)
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
    flag.BoolVar(&shuffle, "s", false, "Shuffle targets")
    flag.StringVar(&launch_db_filepath, "c", "", "Path to a custom launch_db file")
    flag.BoolVar(&debug, "d", false, "Print debug messages")
    flag.BoolVar(&help, "h", false, "Show this message and exit")
    flag.BoolVar(&no_ui, "n", false, "No UI, just launch all targets one by one")
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
    if debug {
        log.Print("Inputs: ", inputs)
    }
}

func launchTarget(target string, command string) {
    cmd := exec.Command(command, target)
    cmd.Stdout = os.Stdout // cmd.Stdout -> stdout
    cmd.Stderr = os.Stderr // cmd.Stderr -> stderr
    cmd.Stdin = os.Stdin // cmd.Stdin <- stdin
    if debug {
        log.Println("Running:", command, target)
    }
    err := cmd.Run()
    if err != nil {
        log.Printf("Command finished with error: %v", err)
    }
}

func launchAll() {
    for _, path := range inputs {
        for pattern, command := range launch_db {
            match, err := filepath.Match(pattern, path)
            if  err != nil {
                log.Fatal(err)
            }
            if match {
                if debug {
                    log.Print("Match: ", path, " :", pattern)
                }
                launchTarget(path, command)
            }
        }
    }
}

func menuLoop() {
}

func main() {
    readArgs()
    readConf()
    if no_ui {
        launchAll()
    } else {
        menuLoop()
    }
}
