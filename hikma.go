package main

import (
	"compress/gzip"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/term"
)

// --- EMBEDDING CONFIGURATION ---

//go:embed hikma.db.gz
var embeddedDBFS embed.FS

const DBName = "hikma.db"

// --- TYPES ---

type Content struct {
	Text   string
	Author string
	Sub    string
}

// --- DATABASE SETUP & EXTRACTION ---

func getDB() *sql.DB {
	// 1. Determine the storage path (Standard Linux: ~/.local/share/hikma/hikma.db)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("❌ Could not find home directory.")
	}

	dataDir := filepath.Join(homeDir, ".local", "share", "hikma")
	dbPath := filepath.Join(dataDir, DBName)

	// 2. Check if the database exists on disk
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Println("📦 First run detected. Installing database...")
		installDB(dataDir, dbPath)
		fmt.Println("✅ Installation complete.")
		fmt.Println()
	}

	// 3. Open the database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("❌ Failed to open database:", err)
	}
	return db
}

func installDB(dirPath, filePath string) {
	// A. Create the directory
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Fatal("❌ Failed to create directory:", err)
	}

	// B. Open the embedded compressed file
	srcFile, err := embeddedDBFS.Open(DBName + ".gz")
	if err != nil {
		log.Fatal("❌ Failed to read embedded data:", err)
	}
	defer srcFile.Close()

	// C. Create the gzip reader
	gzReader, err := gzip.NewReader(srcFile)
	if err != nil {
		log.Fatal("❌ Failed to create gzip reader:", err)
	}
	defer gzReader.Close()

	// D. Create the destination file
	destFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal("❌ Failed to create database file:", err)
	}
	defer destFile.Close()

	// E. Decompress and Copy
	if _, err := io.Copy(destFile, gzReader); err != nil {
		log.Fatal("❌ Failed to write database:", err)
	}
}

// --- DATA FETCHING ---

func getRandomRow(db *sql.DB, mode string, eraFilter string) *Content {
	var table, countQuery, dataQuery string
	var args []interface{}

	if mode == "poems" {
		table = "poetry"
		if eraFilter != "" {
			countQuery = "SELECT COUNT(*) FROM poetry WHERE poet_era LIKE ?"
			dataQuery = "SELECT poem_text, poet_name, poet_era FROM poetry WHERE poet_era LIKE ? LIMIT 1 OFFSET ?"
			args = append(args, "%"+eraFilter+"%")
		} else {
			countQuery = "SELECT MAX(id) FROM poetry"
			dataQuery = "SELECT poem_text, poet_name, poet_era FROM poetry WHERE id = ?"
		}
	} else {
		table = "quotes"
		condition := "category != 'hadith'"
		if mode == "hadith" {
			condition = "category = 'hadith'"
		}
		countQuery = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, condition)
		dataQuery = fmt.Sprintf("SELECT text, author, category FROM %s WHERE %s LIMIT 1 OFFSET ?", table, condition)
	}

	var max int
	err := db.QueryRow(countQuery, args...).Scan(&max)
	if err != nil || max == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())

	// Retry loop for gaps in IDs
	for i := 0; i < 3; i++ {
		offset := rand.Intn(max)
		if mode == "poems" && eraFilter == "" {
			offset = offset + 1
		} // ID lookup starts at 1

		finalArgs := append(args, offset)
		row := db.QueryRow(dataQuery, finalArgs...)

		var c1, c2, c3 string
		if err := row.Scan(&c1, &c2, &c3); err == nil {
			if mode == "poems" {
				// Poetry Formatting
				lines := strings.Split(c1, "\n")
				cleanLines := []string{}
				for _, l := range lines {
					if len(strings.TrimSpace(l)) > 2 {
						cleanLines = append(cleanLines, strings.TrimSpace(l))
					}
				}

				text := c1
				if len(cleanLines) >= 2 {
					idx := 0
					if len(cleanLines) > 2 {
						limit := len(cleanLines)
						if limit%2 != 0 {
							limit--
						}
						idx = rand.Intn(limit/2) * 2
					}
					text = cleanLines[idx+1] + "   ۞   " + cleanLines[idx]
				} else if len(cleanLines) == 1 {
					text = cleanLines[0]
				}
				return &Content{Text: text, Author: c2, Sub: c3}
			} else {
				// Quote/Hadith Formatting
				sub := "Wisdom"
				if c3 == "hadith" {
					sub = "حديث نبوي"
					c2 = translateAuthor(c2)
				}
				return &Content{Text: c1, Author: c2, Sub: sub}
			}
		}
	}
	return nil
}

func translateAuthor(name string) string {
	translations := map[string]string{
		"Bukhari": "الإمام البخاري", "Muslim": "الإمام مسلم",
		"Tirmidhi": "الترمذي", "Abu Dawood": "أبو داود",
		"Abudawood": "أبو داود", "Ibn Majah": "ابن ماجه",
		"Ibnmajah": "ابن ماجه", "Nasai": "النسائي",
		"Malik": "الإمام مالك", "Ahmed": "الإمام أحمد",
	}
	if val, ok := translations[strings.TrimSpace(name)]; ok {
		return val
	}
	return name
}

// --- DISPLAY ---

func printFancy(data *Content) {
	if data == nil {
		return
	}

	cyan, gold, green, grey, reset := "\033[1;96m", "\033[1;93m", "\033[1;92m", "\033[90m", "\033[0m"

	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	if width > 80 {
		width = 80
	}
	if width < 40 {
		width = 40
	}

	bar := strings.Repeat("▀", width-8)
	bottom := strings.Repeat("▄", width-8)

	fmt.Println()
	fmt.Printf("    %s%s%s\n", gold, bar, reset)
	fmt.Printf("    %s  %s%s\n", cyan, data.Text, reset)

	subColor := grey
	if strings.Contains(data.Sub, "Hadith") || strings.Contains(data.Sub, "حديث") {
		subColor = green
	}

	fmt.Printf("\n    %s  %s | %s%s\n", subColor, data.Author, data.Sub, reset)
	fmt.Printf("    %s%s%s\n", gold, bottom, reset)
	fmt.Println()
}

// --- MAIN ---

func main() {
	var showPoems, showQuotes, showHadith bool
	var eraFilter string

	// Custom Help
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "\nUsage: %s [options]\n\n", os.Args[0])
		fmt.Fprintln(w, "A terminal companion for Arabic Poetry, Quotes, and Hadith.\n")
		fmt.Fprintln(w, "Options:")
		fmt.Fprintln(w, "  -d, --hadith      Show Prophetic Hadith")
		fmt.Fprintln(w, "  -p, --poems       Show Poetry")
		fmt.Fprintln(w, "  -q, --quotes      Show Wisdom Quotes")
		fmt.Fprintln(w, "  -e, --era <name>  Filter Poems by Era (e.g., 'Abbasid')")
		fmt.Fprintln(w, "  -h, --help        Show this help message")
		fmt.Fprintln(w, "")
	}

	flag.BoolVar(&showPoems, "p", false, "")
	flag.BoolVar(&showPoems, "poems", false, "")
	flag.BoolVar(&showQuotes, "q", false, "")
	flag.BoolVar(&showQuotes, "quotes", false, "")
	flag.BoolVar(&showHadith, "d", false, "")
	flag.BoolVar(&showHadith, "hadith", false, "")
	flag.BoolVar(&showHadith, "hadeeth", false, "")
	flag.StringVar(&eraFilter, "e", "", "")
	flag.StringVar(&eraFilter, "era", "", "")

	flag.Parse()

	db := getDB()
	defer db.Close()

	mode := "all"
	if showPoems || eraFilter != "" {
		mode = "poems"
	} else if showQuotes {
		mode = "quotes"
	} else if showHadith {
		mode = "hadith"
	} else {
		modes := []string{"poems", "quotes", "hadith"}
		rand.Seed(time.Now().UnixNano())
		mode = modes[rand.Intn(len(modes))]
	}

	content := getRandomRow(db, mode, eraFilter)
	if content == nil && mode == "all" {
		content = getRandomRow(db, "quotes", "")
		if content == nil {
			content = getRandomRow(db, "poems", "")
		}
	}

	if content != nil {
		printFancy(content)
	} else {
		fmt.Println("⚠️  No content found for this category.")
	}
}
