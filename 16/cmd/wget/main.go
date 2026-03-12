package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vvigg0/wbtech-l2/16/internal/wget"
)

func main() {
	cfg := wget.Config{}
	flag.IntVar(&cfg.Workers, "w", 1, "Количество одновременных загрузок(min=1)")
	flag.IntVar(&cfg.Depth, "d", -1, "Глубина скачивания ссылок. d=0 - одна страница,default = бесконечная глубина(d=-1).")
	flag.StringVar(&cfg.RootPath, "p", "", "Путь для скачивания")
	flag.Parse()
	if cfg.Workers < 1 || cfg.Depth < -1 {
		flag.Usage()
		os.Exit(1)
	}
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "не указан url")
		os.Exit(1)
	}
	cfg.URL = args[0]
	if cfg.RootPath == "" {
		path, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка получения корневой директории:", err)
			os.Exit(1)
		}
		cfg.RootPath = path
	}
	if err := wget.Run(&cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
