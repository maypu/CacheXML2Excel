package utils

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
)

func SetExcelHeader(f *excelize.File, sheetIndex int) *excelize.File {
	//找名称
	sheetName := f.GetSheetName(sheetIndex)
	// 设置表头
	f.SetCellValue(sheetName, "A1", "产品组")
	f.SetCellValue(sheetName, "B1", "表名")
	f.SetCellValue(sheetName, "C1", "中文名")
	f.SetCellValue(sheetName, "D1", "表说明")
	f.SetCellValue(sheetName, "E1", "命名空间")
	f.SetCellValue(sheetName, "F1", "数据类型")
	f.SetCellValue(sheetName, "G1", "iMedical版本号")
	f.SetCellValue(sheetName, "H1", "年预估增长量")
	f.SetCellValue(sheetName, "I1", "管理岗位")
	f.SetCellValue(sheetName, "J1", "字段")
	f.SetCellValue(sheetName, "K1", "字段名称")
	f.SetCellValue(sheetName, "L1", "字段说明")
	f.SetCellValue(sheetName, "M1", "字段类型")
	f.SetCellValue(sheetName, "N1", "关联表")
	f.SetCellValue(sheetName, "O1", "值域范围")
	f.SetCellValue(sheetName, "P1", "重要等级")
	f.SetCellValue(sheetName, "Q1", "iMedical版本号")
	f.SetCellValue(sheetName, "R1", "项目组(非标版)")
	//设置单元格的值
	return f
}

func GetRowIndex()  {
	f, err := excelize.OpenFile("D:\\Document\\公卫\\HOS表结构导入模板222.xlsx")
	if err != nil {
		log.Fatal(err)
		ExitProgram()
	}
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		ExitProgram()
	}
	fmt.Println(len(rows))
}