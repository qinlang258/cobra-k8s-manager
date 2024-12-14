package mtable

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

var ImageTatleColumns []string = []string{
	"NAMESPACE",
	"TYPE",
	"RESOURCE_NAME",
	"CONTAINER_NAME",
	"IMAGE",
}

var ResourceTatleColumns []string = []string{
	"NAMESPACE",
	"TYPE",
	"RESOURCE_NAME",
	"CONTAINER_NAME",
	"CPU_LIMIT",
	"CPU_REQUESTS",
	"MEMORY_LIMIT",
	"MEMORY_REQUESTS",
}

var TopTatleColumns []string = []string{
	"NAMESPACE",
	"TYPE",
	"RESOURCE_NAME",
	"POD_NAME",
	"CPU_USED",
	"MEMORY_USED",
}

func TablePrint(tableName string, data []map[string]string) {
	// Use tablewriter to create a nice formatted table
	table := tablewriter.NewWriter(os.Stdout)
	var TatleColumns []string
	if tableName == "image" {
		TatleColumns = ImageTatleColumns

		table.SetHeader(TatleColumns) // Table header

		// Add rows to the table
		for _, row := range data {
			table.Append([]string{
				row["NAMESPACE"],
				row["TYPE"],
				row["RESOURCE_NAME"],
				row["CONTAINER_NAME"],
				row["IMAGE"],
			})
		}
	} else if tableName == "resource" {
		TatleColumns = ResourceTatleColumns

		table.SetHeader(TatleColumns) // Table header

		// Add rows to the table
		for _, row := range data {
			table.Append([]string{
				row["NAMESPACE"],
				row["TYPE"],
				row["RESOURCE_NAME"],
				row["CONTAINER_NAME"],
				row["CPU_LIMIT"],
				row["CPU_REQUESTS"],
				row["MEMORY_LIMIT"],
				row["MEMORY_REQUESTS"],
			})
		}
	} else if tableName == "top" {
		TatleColumns = TopTatleColumns

		table.SetHeader(TatleColumns) // Table header

		// Add rows to the table
		for _, row := range data {
			table.Append([]string{
				row["NAMESPACE"],
				row["TYPE"],
				row["RESOURCE_NAME"],
				row["POD_NAME"],
				row["CPU_USED"],
				row["MEMORY_USED"],
			})
		}
	}

	// Render the table
	table.Render()
}
