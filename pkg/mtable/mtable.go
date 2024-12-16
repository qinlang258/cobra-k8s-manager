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

var AnalysisNodeTatleColumns []string = []string{
	"NODE_NAME",
	"NAMESPACE",
	"POD_NAME",
	"CPU_USED",
	"CPU_PERCENT",
	"MEMORY_USED",
	"MEMORY_PERCENT",
}

var NodeTatleColumns []string = []string{
	"NODE_NAME",
	"NODE_ADDRESS",
	"OS_IMAGE",
	"KUBELET_VERSION",
	"CONTAINER_RUNTIME_VERSION",
	"CPU_USED",
	"CPU_TOTAL",
	"CPU_PERCENT",
	"MEMORY_USED",
	"MEMORY_TOTAL",
	"MEMORY_PERCENT",
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
	} else if tableName == "node" {
		TatleColumns = NodeTatleColumns

		table.SetHeader(TatleColumns) // Table header

		// Add rows to the table
		for _, row := range data {
			table.Append([]string{
				row["NODE_NAME"],
				row["NODE_ADDRESS"],
				row["OS_IMAGE"],
				row["KUBELET_VERSION"],
				row["CONTAINER_RUNTIME_VERSION"],
				row["CPU_USED"],
				row["CPU_TOTAL"],
				row["CPU_PERCENT"],
				row["MEMORY_USED"],
				row["MEMORY_TOTAL"],
				row["MEMORY_PERCENT"],
			})
		}
	} else if tableName == "analysis" {
		TatleColumns = AnalysisNodeTatleColumns

		table.SetHeader(TatleColumns) // Table header

		// Add rows to the table
		for _, row := range data {
			table.Append([]string{
				row["NODE_NAME"],
				row["NAMESPACE"],
				row["POD_NAME"],
				row["CPU_USED"],
				row["CPU_PERCENT"],
				row["MEMORY_USED"],
				row["MEMORY_PERCENT"],
			})
		}
	}

	// Render the table
	table.Render()

}
