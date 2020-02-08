package main

import (
	"flag"
	"fmt"
	ink "github.com/dennwc/inkview"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Strubbl/wallabago"
	"github.com/bmaupin/go-epub"
)

const version = "0.1"
const defaultConfigJSON = "/mnt/ext1/system/config/wallabook.json"

var debug = flag.Bool("d", false, "get debug output (implies verbose mode)")
var debugDebug = flag.Bool("dd", false, "get even more debug output like data (implies debug mode)")
var v = flag.Bool("v", false, "print version")
var verbose = flag.Bool("verbose", false, "verbose mode")
var configJSON = flag.String("config", defaultConfigJSON, "file name of config JSON file")

func handleFlags() {
	flag.Parse()
	if *debug && len(flag.Args()) > 0 {
		log.Printf("handleFlags: non-flag args=%v", strings.Join(flag.Args(), " "))
	}
	// version first, because it directly exits here
	if *v {
		fmt.Printf("version %v\n", version)
		os.Exit(0)
	}
	// test verbose before debug because debug implies verbose
	if *verbose && !*debug && !*debugDebug {
		log.Printf("verbose mode")
	}
	if *debug && !*debugDebug {
		log.Printf("handleFlags: debug mode")
		// debug implies verbose
		*verbose = true
	}
	if *debugDebug {
		log.Printf("handleFlags: debug² mode")
		// debugDebug implies debug
		*debug = true
		// and debug implies verbose
		*verbose = true
	}
}

func checkError(err error, op string) {
	if err != nil {
		log.Printf("Op error %s: %s", op, err)
		os.Exit(1)
	}
}

func runWallabook(w io.Writer) {
	ink.InitCerts()
	os.Chdir(filepath.Dir(os.Args[0]))
	log.SetOutput(w)
	handleFlags()
	// check for config
	if *verbose {
		log.Println("reading config", *configJSON)
	}
	err := wallabago.ReadConfig(*configJSON)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	numArticles, err := wallabago.GetNumberOfTotalArticles()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	log.Printf("There are %d articles", numArticles)

	entries, err := wallabago.GetAllEntries()
	checkError(err, "GetAllEntries")

	e := epub.NewEpub("Wallabooks")
	e.SetAuthor("Wallabook")

	for _, entry := range entries {
		if len(entry.Content) > 500 {
			html := fmt.Sprintf("<h1>%s</h1>\n%s", entry.Title, entry.Content)
			e.AddSection(html, entry.Title, "", "")
		} else {
			log.Printf("Skipping %s, too short (%d)\n", entry.Title, len(entry.Content))
		}

	}

	err = e.Write("/mnt/ext1/Wallabook/result.epub")
	checkError(err, "WriteEpub")
}
