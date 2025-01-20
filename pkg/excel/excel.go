package excel

import (
	"context"
	"fmt"
	"time"

	"k8s-manager/pkg/config"
	"k8s-manager/pkg/mtable"

	"github.com/xuri/excelize/v2"
	"k8s.io/klog/v2"
)

func ExportXlsx(ctx context.Context, tableName string, dataList []map[string]string, kubeconfig string) {
	// 如果 dataList 为空，直接返回
	if len(dataList) == 0 {
		fmt.Println("No data to export")
		return
	}

	name, err := config.GetClusterNameFromPrometheusUrl(kubeconfig)
	if err != nil {
		klog.Error(ctx, err.Error())
	}

	// 创建 Excel 文件
	file := excelize.NewFile()
	sheetName := "Sheet1"

	// 根据 tableName 选择相应的表头
	var headers []string
	if tableName == "image" {
		headers = mtable.ImageTatleColumns
	} else if tableName == "resource" {
		headers = mtable.ResourceTatleColumns
	} else if tableName == "top" {
		headers = mtable.TopTatleColumns
	} else if tableName == "node" {
		headers = mtable.NodeTatleColumns
	} else if tableName == "analysis" {
		headers = mtable.AnalysisNodeTatleColumns
	} else if tableName == "analysis-cpu-memory" {
		headers = mtable.AnalysisCpuMemory
	}

	// 设置表头（从 A1 开始，横向填充表头）
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string('A'+i))
		file.SetCellValue(sheetName, cell, header)
	}

	// 填充数据
	for i, app := range dataList {
		row := i + 2 // 从第2行开始填充数据
		for j, header := range headers {
			cell := fmt.Sprintf("%s%d", string('A'+j), row)
			file.SetCellValue(sheetName, cell, app[header])
		}
	}

	// 加载东八区（Asia/Shanghai）时区
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("Error loading time zone:", err)
		return
	}

	// 获取当前东八区时间
	currentTime := time.Now().In(location)
	excelName := fmt.Sprintf("%s集群-%s-%s.xlsx", name, tableName, currentTime.Format("2006-01-02-150405"))
	// 保存文件
	if err := file.SaveAs(excelName); err != nil {
		panic(err)
	}
}
