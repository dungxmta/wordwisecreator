package main

import (
	"bufio"
	"fmt"
	strip "github.com/grokify/html-strip-tags-go"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Worker struct {
	FilePath string
}

// Read line by line input file =/> raw line
// Process line:
//  strip html tag, replace special char with space
//  split line to words
//  loop words & check in wordwise dict to replace with ruby tag =/> normalized dict
//  loop normalized dict to replace normalized word in raw line =/> normalized line
//  write normalized line to output file
//
func (obj *Worker) run(poolCh *chan bool, stopWords *map[string]bool, dict *map[string]Row) {
	s := strings.Split(obj.FilePath, "/")
	fileName := s[len(s)-1]

	log.Println("[+] worker start...", fileName)

	defer func() {
		if poolCh != nil {
			<-*poolCh
		}
		wg.Done()
		log.Println("[-] worker done !!!", fileName)
	}()

	// read inp file
	fi, err := os.Open(obj.FilePath)
	if err != nil {
		log.Fatalln("Error when open ", obj.FilePath, "->", err)
	}
	defer fi.Close()

	// create out file
	outFile := path.Join(pathOutput, fileName)
	fo, err := os.Create(outFile)
	if err != nil {
		log.Fatalln("Error when create output ", outFile, "->", err)
	}
	defer fo.Close()

	scanner := bufio.NewScanner(fi)
	scanner.Split(bufio.ScanLines)
	// scanner.Split(bufio.ScanWords)

	writer := bufio.NewWriter(fo)

	count := 0
	var ignore bool
	var cleanLine = ""
	var words []string

	for scanner.Scan() {
		line := scanner.Text()
		// if strings.HasPrefix(word, "#") {
		// 	continue
		// }

		// ignore some row
		ignore = skipLine(line)
		if ignore {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				log.Fatalln(err)
			}
			continue
		}

		// clear html tag and some char
		cleanLine = CleanLine(line)

		// split clean line to words
		words = strings.Split(cleanLine, " ")
		// var wordsOut []string
		normalized := make(map[string]string)

		// build dict before replace normalized value
		for _, word := range words {
			word = strings.Trim(word, " ")

			if _, ok := (*stopWords)[word]; ok {
				// wordsOut = append(wordsOut, word)
				continue
			}

			ws, ok := (*dict)[word]
			if !ok {
				continue
			}

			if ws.HintLv > hintLevel {
				continue
			}

			normalized[word] = fmt.Sprintf("<ruby>%v<rt>%v</rt></ruby>", word, ws.ShortDef)
		}

		// replace normalized value in full line
		if len(normalized) > 0 {
			var r []string
			for k, v := range normalized {
				r = append(r, k, v)
			}

			replacer := strings.NewReplacer(r...)
			line = replacer.Replace(line)
		}

		count++
		// log.Println(cleanLine)
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatalln(err)
		}
	}

	if err = writer.Flush(); err != nil {
		log.Fatalln(err)
	}

	if scanner.Err() != nil {
		log.Fatalln("Error when scan word ", "->", err)
	}

	// ss := rand.Intn(10)
	// time.Sleep(time.Second * time.Duration(ss))
}

func skipLine(s string) bool {
	s = strings.Trim(s, " ")
	if s == "" {
		return true
	}

	for _, v := range []string{
		`<?xml`,
		`<html`,
		`<head`,
		`<title`,
		`<link`,
		`</head`,
		`<body`,
		`</body`,
		`</html`,
	} {
		if strings.HasPrefix(s, v) {
			return true
		}
	}

	return false
}

// strip html tags
// replace special char with scpace
func CleanLine(s string) string {
	out := strip.StripTags(s)

	replacer := strings.NewReplacer(lstReplace...)
	out = replacer.Replace(out)

	return out
}
