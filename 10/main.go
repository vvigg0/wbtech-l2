package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	mysort "github.com/vvigg0/wbtech-l2/10/pkg/sorttool"
)

func main() {
	var fileName string
	config := mysort.Config{}

	flag.StringVar(&fileName, "f", "", "name of file to open")

	flag.IntVar(&config.KeyColumn, "k", 0, "sort by column N (1-based)")
	flag.BoolVar(&config.Numeric, "n", false, "numeric sort")
	flag.BoolVar(&config.Reverse, "r", false, "reverse order")
	flag.BoolVar(&config.Unique, "u", false, "unique lines only")
	flag.BoolVar(&config.Month, "M", false, "sort by month name")
	flag.BoolVar(&config.TrimSpace, "b", false, "ignore trailing blanks")
	flag.BoolVar(&config.Check, "c", false, "check if sorted")
	flag.BoolVar(&config.Human, "h", false, "human numeric sort")
	flag.StringVar(&config.Separator, "sep", "\t", "set separator")
	flag.Parse()

	cmp, err := mysort.BuildComparator(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "BuildComparator error: %v", err)
		os.Exit(1)
	}
	var lines []string
	if fileName == "" {
		lines, err = readFile(os.Stdin)
	} else {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Open file error: %v", err)
			os.Exit(1)
		}
		defer file.Close()
		lines, err = readFile(file)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "readFile error:", err)
		os.Exit(1)
	}
	if config.Check {
		ok, idx := mysort.CheckSort(lines, cmp, config.Unique)
		if ok {
			fmt.Println("Данные отсортированы")
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Данные не отсортированы, проверка завершена на %v строке\n", idx+1)
		os.Exit(1)
	}
	mysort.Sort(lines, cmp)
	if config.Unique {
		lines = mysort.UniqueSorted(lines, cmp)
	}
	for _, v := range lines {
		fmt.Println(v)
	}
}
func readFile(src *os.File) ([]string, error) {
	var lines []string

	reader := bufio.NewReader(src)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			if len(line) > 0 {
				line = strings.TrimSuffix(line, "\n")
				line = strings.TrimSuffix(line, "\r")
				lines = append(lines, line)
			}
			break
		}
		if err != nil {
			return nil, err
		}
		line = strings.TrimSuffix(line, "\n")
		line = strings.TrimSuffix(line, "\r")
		lines = append(lines, line)
	}
	return lines, nil
}
