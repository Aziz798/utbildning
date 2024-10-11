package files

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

// Define a struct to represent the data
type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Mail      string `json:"mail"`
}

func Analyze() {
	// Open the Excel file
	filePath := "./file.xlsx"
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer f.Close()

	// Create a slice to store the extracted data
	var people []Person

	// Get the first sheet name
	sheetName := f.GetSheetName(0)

	// Get all rows in the sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("Failed to get rows: %v", err)
	}

	// Iterate over the rows, skipping the header (assuming the first row is the header)
	for i, row := range rows {
		if i == 0 {
			continue // Skip the header row
		}

		// Check if the row has enough columns (C is 2nd index, D is 3rd, and J is 9th index)
		if len(row) > 9 {
			// Extract the values from columns C, D, and J
			person := Person{
				FirstName: row[2], // Column C (2nd index)
				LastName:  row[3], // Column D (3rd index)
				Mail:      row[9], // Column J (9th index)
			}
			// Add the person to the slice
			people = append(people, person)
		}
	}

	// Create a JSON file
	outputFile, err := os.Create("output.json")
	if err != nil {
		log.Fatalf("Failed to create output.json: %v", err)
	}
	defer outputFile.Close()

	// Encode the data into JSON and write it to the file
	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ") // Pretty print the JSON
	err = encoder.Encode(people)
	if err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	fmt.Println("Data successfully written to output.json")
}
