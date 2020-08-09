package main

import (
	"log"
	"sync"
)

const (
	pathStopWord = `./stopwords.txt`
	pathDictCSV  = `./wordwise-dict.csv`

	pathHtml   = "./data/extract/OEBPS/"
	pathOutput = "./data/output/"

	hintLevel = 3

	maxWorker = 10
)

var wg sync.WaitGroup

func cleanTmp() {
	// TODO
}

// loop in list file and manage worker
func Start(lstSrc []string, stopWords *map[string]bool, dict *map[string]Row) {
	poolCh := make(chan bool, maxWorker)

	// workers pop from queue and process file
	for _, src := range lstSrc {
		// log.Println("pool", len(poolCh))
		poolCh <- true
		wg.Add(1)

		worker := Worker{FilePath: src}
		go worker.run(&poolCh, stopWords, dict)
	}

	wg.Wait()
}

func main() {
	log.Println("[+] Load stopwords...")
	stopWords := loadStopWords()

	log.Println("[+] Load wordwise dict...")
	dict := loadDict()

	log.Println("[+] Load source...")
	lstSrc := loadSource()

	log.Println("[+] Starting...")
	Start(lstSrc, stopWords, dict)

	log.Println("[+] Cleaning temp files...")
	cleanTmp()

	log.Println("-> done!!!")
}
