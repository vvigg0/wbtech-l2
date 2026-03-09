package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	mygrep "github.com/vvigg0/wbtech-l2/12/internal/grepanalog"
	ringbuffer "github.com/vvigg0/wbtech-l2/12/internal/ringbuffer"
)

func main() {
	config := mygrep.Config{}
	var around int

	flag.IntVar(&config.After, "A", 0, "вывести N строк после совпадения")
	flag.IntVar(&config.Before, "B", 0, "вывести N строк до совпадения")
	flag.IntVar(&around, "C", 0, "вывести N строк контекста до и после совпадения")
	flag.BoolVar(&config.CountOnly, "c", false, "вывести только количество совпадений")
	flag.BoolVar(&config.IgnoreCase, "i", false, "игнорировать регистр")
	flag.BoolVar(&config.Invert, "v", false, "вывести строки без совпадения")
	flag.BoolVar(&config.Fixed, "F", false, "искать фиксированную строку, а не регулярное выражение")
	flag.BoolVar(&config.LineNum, "n", false, "вывести номер строки")

	flag.Parse()

	config.Pattern = flag.Arg(0)
	fileName := flag.Arg(1)

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(2)
	}

	if around != 0 {
		config.After = around
		config.Before = around
	}
	if config.CountOnly {
		config.After = 0
		config.Before = 0
		config.LineNum = false
	}

	match, err := mygrep.BuildMatcher(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "BuildMatcher error: %v\n", err)
		os.Exit(1)
	}

	var r io.Reader = os.Stdin
	if fileName != "" {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Open file error: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		r = file
	}

	buff := ringbuffer.Create(config.Before)

	var lineNo int
	var count int
	var after int
	var lastPrintedNo int
	var printedAny bool

	printEntry := func(e ringbuffer.Entry) {
		if e.No <= lastPrintedNo {
			return
		}
		if config.LineNum {
			fmt.Printf("%d:%s\n", e.No, e.Text)
		} else {
			fmt.Println(e.Text)
		}
		lastPrintedNo = e.No
	}

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		lineNo++
		m := match(line)
		if config.CountOnly {
			if m {
				count++
			}
			continue
		}
		if m {
			snap := buff.Snapshot()

			start := lineNo
			for _, e := range snap {
				if e.No > lastPrintedNo {
					start = e.No
					break
				}
			}

			if printedAny && start > lastPrintedNo+1 {
				fmt.Println("--")
			}

			for _, e := range snap {
				printEntry(e)
			}
			if config.After > after {
				after = config.After
			}
			printEntry(ringbuffer.Entry{No: lineNo, Text: line})
			printedAny = true
		} else if after > 0 {
			printEntry(ringbuffer.Entry{No: lineNo, Text: line})
			after--
		}
		buff.Push(ringbuffer.Entry{No: lineNo, Text: line})
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v\n", err)
		os.Exit(1)
	}

	if config.CountOnly {
		fmt.Println(count)
	}
}
