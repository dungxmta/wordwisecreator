package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	bookFile     string
	bookPath     string
	bookFilename string

	specialChar []string
	lstReplace  []string
	// re = regexp.MustCompile(`/[^ \w]+/`)
)

func init() {
	specialChar = []string{
		`,`, `<`, `>`, `;`, `&`, `*`, `~`, `/`,
		`"`, `[`, `]`, `#`, `?`, "`", `–`, `.`,
		"`", `"`, `"`, `!`, `“`, `”`, `:`, `.`,
	}
	for _, v := range specialChar {
		lstReplace = append(lstReplace, v, " ")
	}
}

type Row struct {
	Word     string
	FullDef  string
	ShortDef string
	Example  string
	HintLv   int
}

// Load Dict from CSV
func loadDict() *map[string]Row {

	file, err := os.Open(pathDictCSV)
	if err != nil {
		log.Fatalln("Error when open ", pathStopWord, "->", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	reader := csv.NewReader(file)

	dict := make(map[string]Row)
	row := Row{}

	// Read each record from csv
	// skip header
	record, err := reader.Read()
	if err == io.EOF {
		log.Fatalln("Empty csv file")
	}

	count := 0

	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Error when scan word ", count, "->", err)
		}

		if len(record) < 6 {
			log.Println("invalid word: ", record)
			continue
		}

		hintLv, err := strconv.Atoi(record[5])
		if err != nil {
			log.Println("cant get hint_level: ", record, "->", err)
			continue
		}

		row = Row{
			Word:     record[1],
			FullDef:  record[2],
			ShortDef: record[3],
			Example:  record[4],
			HintLv:   hintLv,
		}

		dict[row.Word] = row
		count++
	}

	log.Println("-> csv words:", count)
	return &dict
}

// Load Stop Words from txt
func loadStopWords() *map[string]bool {
	dict := make(map[string]bool)

	file, err := os.Open(pathStopWord)
	if err != nil {
		log.Fatalln("Error when open ", pathStopWord, "->", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	// scanner.Split(bufio.ScanWords)

	count := 0
	for scanner.Scan() {
		word := scanner.Text()
		if strings.HasPrefix(word, "#") {
			continue
		}

		if word != "" {
			dict[word] = true
			count++
			// log.Println(word)
		}
	}

	if scanner.Err() != nil {
		log.Fatalln("Error when scan word ", "->", err)
	}

	log.Println("-> stop words:", count)
	return &dict
}

func loadSource() []string {
	// TODO
	// change .epub to .zip and extract it

	// get list .html/xhtml file
	var lstSrc []string

	files, err := ioutil.ReadDir(pathHtml)
	if err != nil {
		log.Fatalln("Error when read dir ", pathHtml, "->", err)
	}

	for idx, file := range files {
		// log.Println(file.Name())
		if file.IsDir() || (!strings.HasSuffix(file.Name(), "html") &&
			!strings.HasSuffix(file.Name(), "xhtml")) {
			continue
		}
		lstSrc = append(lstSrc, path.Join(pathHtml, file.Name()))

		// TODO: test... remove me
		if idx > 0 {
			// break
		}
	}

	if len(lstSrc) == 0 {
		log.Fatalln("No source file...")
	}

	log.Println("-> files: ", len(lstSrc))
	return lstSrc
}
