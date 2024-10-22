package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
)

func main() {
	verbose := flag.Bool("v", false, "Enable verbose debugging output")
	dateStr := flag.String("date", time.Now().Format("2006-01-02"), "Date to extract (YYYY-MM-DD)")
	width := flag.Int("width", 0, "Force specific terminal width (0 for auto-detect)")
	flag.Parse()

	debug := log.New(os.Stderr, "DEBUG: ", log.Ltime|log.Lmicroseconds)

	// Get terminal width
	terminalWidth := *width
	if terminalWidth == 0 {
		if ws, err := getWindowSize(); err == nil {
			// Subtract a small buffer to prevent wrapping issues
			terminalWidth = int(ws.Col) - 2
			if *verbose {
				debug.Printf("Detected terminal width: %d, using width: %d", ws.Col, terminalWidth)
			}
		} else {
			terminalWidth = 80
			if *verbose {
				debug.Printf("Failed to detect terminal width, using default: %d", terminalWidth)
			}
		}
	}

	// Parse the date and get the day of week
	date, err := time.Parse("2006-01-02", *dateStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid date format: %v\n", err)
		os.Exit(1)
	}

	// Format the heading we're looking for
	targetHeading := fmt.Sprintf("### %s %s", date.Format("Monday"), date.Format("2006-01-02"))
	if *verbose {
		debug.Printf("Looking for heading: %q", targetHeading)
	}

	file, err := os.Open("lesson-planning.md")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var capturing bool
	var foundContent bool
	var content []string
	var lastLineEmpty bool

	// Read line by line
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// If we find our target heading, start capturing
		if trimmedLine == targetHeading {
			if *verbose {
				debug.Printf("Found target heading")
			}
			capturing = true
			foundContent = true
			content = append(content, line)
			lastLineEmpty = false
			continue
		}

		// If we're capturing and hit another heading, stop
		if capturing && strings.HasPrefix(trimmedLine, "### ") {
			if *verbose {
				debug.Printf("Found next heading, stopping")
			}
			break
		}

		// Collect lines while we're capturing
		if capturing {
			if trimmedLine == "" {
				if !lastLineEmpty {
					content = append(content, "")
					lastLineEmpty = true
				}
			} else {
				content = append(content, line)
				lastLineEmpty = false
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if !foundContent {
		fmt.Printf("No lessons found for %s\n", *dateStr)
		os.Exit(0)
	}

	// Create glamour renderer with explicit width setting
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(terminalWidth),
		glamour.WithPreservedNewLines(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating renderer: %v\n", err)
		os.Exit(1)
	}

	// Join and render the content
	finalContent := strings.Join(content, "\n")

	if *verbose {
		debug.Printf("Rendering content with width: %d", terminalWidth)
	}

	rendered, err := renderer.Render(finalContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering content: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(rendered)
}

type windowSize struct {
	Row, Col uint16
}
