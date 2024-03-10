package main

import (
	"CacheXML2Excel/utils"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
	"time"
)

func main() {
	fmt.Println("请输入你想要扫描的xml文件列表的绝对路径，按回车结束（如：D:\\开发环境\\传染病\\class）")
	fmt.Print("-> ")
	var inputText string
	fmt.Scanln(&inputText)
	fmt.Println("接收到的内容是：",inputText)
	fmt.Println("回车继续，如需修改请重新输入：")
	fmt.Print("-> ")
	var inputText2 string
	fmt.Scanln(&inputText2)
	if inputText2 != "" {
		inputText = inputText2
	}
	if inputText != "" {
		fileNameList := utils.GetDirFiles(inputText)
		if len(fileNameList)<1 {
			fmt.Println("没有找到表结构的xml文件！")
			utils.ExitProgram()
		}
		//创建Excel
		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		}()
		// 创建一个工作表
		sheetIndex, err := f.NewSheet("Sheet1")
		if err != nil {
			log.Fatal(err)
		}
		sheetName := f.GetSheetName(sheetIndex)
		//设置表头
		f = utils.SetExcelHeader(f, sheetIndex)

		rowIndex := 2
		for _, fileName := range fileNameList {
			cacheTableClass := utils.ReadXml(fileName)
			if len(cacheTableClass.Property)<2 {
				fmt.Println("找到服务类：" + cacheTableClass.Name + " 已跳过")
				continue
			}
			fmt.Println("找到表：" + cacheTableClass.SqlTableName)
			//插入表名等信息
			columnIndex := 'B'	//rune
			//rowIndex = rowIndex + 1
			f.SetCellValue(sheetName, string(columnIndex) + strconv.Itoa(rowIndex), cacheTableClass.SqlTableName)	//B2
			columnIndex = columnIndex + 1
			f.SetCellValue(sheetName, string(columnIndex) + strconv.Itoa(rowIndex), cacheTableClass.Description)	//C2
			columnIndex = columnIndex + 2
			f.SetCellValue(sheetName, string(columnIndex) + strconv.Itoa(rowIndex), "DHC-APP")	//E2
			columnIndex = columnIndex + 2
			f.SetCellValue(sheetName, string(columnIndex) + strconv.Itoa(rowIndex), "iMedical 8.0")	//G2

			for _, property := range cacheTableClass.Property {
				fmt.Println("-> 找到字段：" + property.Name + " " + property.ChineseType + " " + property.Description)
				//插入字段名等信息
				columnIndex = 'J'
				f.SetCellValue(sheetName, string(columnIndex) + strconv.Itoa(rowIndex), property.SqlFieldName)	//J2
				columnIndex = columnIndex + 1
				f.SetCellValue(sheetName, string(columnIndex) + strconv.Itoa(rowIndex), property.Description)	//K2
				columnIndex = columnIndex + 2
				f.SetCellValue(sheetName, string(columnIndex) + strconv.Itoa(rowIndex), property.ChineseType)	//M2
				columnIndex = columnIndex + 1
				f.SetCellValue(sheetName, string(columnIndex) + strconv.Itoa(rowIndex), property.PointerTable)	//N2
				rowIndex = rowIndex + 1
			}
		}
		// 设置工作簿的默认工作表
		f.SetActiveSheet(sheetIndex)
		nowDateTime := time.Now().Format("20060102") + time.Now().Format("150405")
		excelFileName := "Cache表结构整理_CacheXML2Excel_"+ nowDateTime +".xlsx"
		f.SetColWidth(sheetName, "A","R",20)
		// 根据指定路径保存文件
		if err := f.SaveAs(excelFileName); err != nil {
			fmt.Println("文件保存失败！")
			log.Fatal(err)
		} else {
			fmt.Println("文件生成成功！")
			fmt.Println("文件保存在当前程序根目录下，名称为：" + excelFileName)
		}

	} else {
		fmt.Println("未输入文件夹路径！")
	}
	utils.ExitProgram()
}