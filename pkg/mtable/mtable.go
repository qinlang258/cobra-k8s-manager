package mtable

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

var ImageTatleColumns []string = []string{
	"NAMESPACE",
	"资源类型",
	"资源名",
	"容器名",
	"镜像地址",
}

var ResourceTatleColumns []string = []string{
	"NAMESPACE",
	"资源类型",
	"资源名",
	"容器名",
	"CPU限制",
	"CPU所需",
	"内存限制",
	"内存所需",
}

var TopTatleColumns []string = []string{
	"NAMESPACE",
	"资源类型",
	"资源名",
	"POD_NAME",
	"已使用的CPU",
	"已使用的内存",
}

var AnalysisNodeTatleColumns []string = []string{
	"节点名",
	"NAMESPACE",
	"POD_NAME",
	"容器名",
	"当前已使用的CPU",
	"CPU使用占服务器的百分比",
	"当前已使用的内存",
	"内存使用占服务器的百分比",
}

var NodeTatleColumns []string = []string{
	"节点名",
	"节点IP",
	"OS镜像",
	"Kubelet版本",
	"CONTAINER_RUNTIME_VERSION",
	"已使用的CPU",
	"CPU总大小",
	"CPU使用占服务器的百分比",
	"已使用的内存",
	"内存总大小",
	"内存使用占服务器的百分比",
}

var AnalysisCpuMemory []string = []string{
	"节点名",
	"NAMESPACE",
	"POD_NAME",
	"容器名",
	"CPU限制",
	"CPU所需",
	"最近7天已使用的CPU",
	"内存限制",
	"内存所需",
	"JAVA-XMX",
	"JAVA-XMS",
	"最近7天已使用的内存",
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
				row["资源类型"],
				row["资源名"],
				row["容器名"],
				row["镜像地址"],
			})
		}
	} else if tableName == "resource" {
		TatleColumns = ResourceTatleColumns

		table.SetHeader(TatleColumns) // Table header

		// Add rows to the table
		for _, row := range data {
			table.Append([]string{
				row["NAMESPACE"],
				row["资源类型"],
				row["资源名"],
				row["容器名"],
				row["CPU限制"],
				row["CPU所需"],
				row["内存限制"],
				row["内存所需"],
			})
		}
	} else if tableName == "top" {
		TatleColumns = TopTatleColumns

		table.SetHeader(TatleColumns) // Table header

		// Add rows to the table
		for _, row := range data {
			table.Append([]string{
				row["NAMESPACE"],
				row["资源类型"],
				row["资源名"],
				row["POD_NAME"],
				row["当前已使用的CPU"],
				row["当前已使用的内存"],
			})
		}
	} else if tableName == "node" {
		TatleColumns = NodeTatleColumns

		table.SetHeader(TatleColumns) // Table header

		// Add rows to the table
		for _, row := range data {
			table.Append([]string{
				row["节点名"],
				row["节点IP"],
				row["OS镜像"],
				row["Kubelet版本"],
				row["CONTAINER_RUNTIME_VERSION"],
				row["当前已使用的CPU"],
				row["CPU总大小"],
				row["CPU使用占服务器的百分比"],
				row["当前已使用的内存"],
				row["内存总大小"],
				row["内存使用占服务器的百分比"],
			})
		}
	} else if tableName == "analysis" {
		TatleColumns = AnalysisNodeTatleColumns

		table.SetHeader(TatleColumns) // Table header

		// Add rows to the table
		for _, row := range data {
			table.Append([]string{
				row["节点名"],
				row["NAMESPACE"],
				row["POD_NAME"],
				row["容器名"],
				row["当前已使用的CPU"],
				row["CPU使用占服务器的百分比"],
				row["当前已使用的内存"],
				row["内存使用占服务器的百分比"],
			})
		}
	} else if tableName == "analysis-cpu-memory" {
		TatleColumns = AnalysisCpuMemory

		table.SetHeader(TatleColumns) // Table header

		// Add rows to the table
		for _, row := range data {
			table.Append([]string{
				row["节点名"],
				row["NAMESPACE"],
				row["POD_NAME"],
				row["容器名"],
				row["CPU限制"],
				row["CPU所需"],
				row["最近7天已使用的CPU"],
				row["内存限制"],
				row["内存所需"],
				row["JAVA-XMX"],
				row["JAVA-XMS"],
				row["最近7天已使用的内存"],
			})
		}
	}

	// Render the table
	table.Render()

}
