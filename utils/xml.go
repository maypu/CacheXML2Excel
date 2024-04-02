package utils

import (
	"github.com/beevik/etree"
	"log"
	"regexp"
	"strings"
)

type CacheTableClass struct {
	Name         string
	Description  string
	ClassType    string
	SqlRowIdName string
	SqlTableName string
	Property     []CacheTableProperty
}

type CacheTableProperty struct {
	Name         string
	Description  string
	Type         string
	ChineseType  string
	SqlFieldName string
	PointerTable string //类型为指针时，指向的表名
}

func ReadXml(path string) CacheTableClass {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		log.Fatal(err)
	}
	export := doc.SelectElement("Export")
	cacheClass := export.SelectElement("Class")

	var cacheTableClass CacheTableClass
	//类名
	cacheTableClass.Name = cacheClass.SelectAttrValue("name", "")
	//计算表名
	cacheTableClassName := strings.Replace(cacheTableClass.Name, ".", "_", -1)
	cacheTableClassName = ReplaceRight(cacheTableClassName, "_", ".", 1)
	//RowID
	if sqlRowIdName := cacheClass.SelectElement("SqlRowIdName"); sqlRowIdName != nil {
		cacheTableClass.SqlRowIdName = sqlRowIdName.Text() //老表如 DHCMed.EPD.Epidemic 会自定义ID名称
	} else {
		cacheTableClass.SqlRowIdName = "ID" //新表默认使用 ID
	}
	cacheTableClass.SqlTableName = cacheTableClassName
	//描述
	if classDescription := cacheClass.SelectElement("Description"); classDescription != nil {
		classDesc := classDescription.Text()
		classDesc = strings.TrimPrefix(classDesc, "\n")
		classDesc = strings.TrimSuffix(classDesc, "\n")
		classDesc = strings.TrimSpace(classDesc)
		cacheTableClass.Description = classDesc
	}
	//计算表描述
	extractTableDesc := ExtractTableDesc(cacheTableClass.Description)
	if extractTableDesc != "" {
		cacheTableClass.Description = extractTableDesc
	}
	for _, property := range cacheClass.SelectElements("Property") {

		if cardinality := property.SelectElement("Cardinality"); cardinality != nil {
			if cardinality.Text() == "children" { //父表中指向子表的字段可以不算
				continue
			}
		}
		cacheTableProperty := CacheTableProperty{
			//字段名
			Name: property.SelectAttrValue("name", ""),
			//字段类型
			Type: property.SelectElement("Type").Text(),
		}
		//字段表名
		if sqlFieldName := property.SelectElement("SqlFieldName"); sqlFieldName != nil {
			cacheTableProperty.SqlFieldName = sqlFieldName.Text()
		} else {
			cacheTableProperty.SqlFieldName = cacheTableProperty.Name
		}
		//字段描述（没有写备注的字段，取SqlFieldName或者Name）
		if propertyDescription := property.SelectElement("Description"); propertyDescription != nil {
			propertyDesc := propertyDescription.Text()
			propertyDesc = strings.TrimPrefix(propertyDesc, "\n")
			propertyDesc = strings.TrimSuffix(propertyDesc, "\n")
			cacheTableProperty.Description = propertyDesc
		} else {
			cacheTableProperty.Description = cacheTableProperty.SqlFieldName
		}
		//字段类型中文名
		cacheTableProperty.ChineseType = PropertyTypeToChinese(cacheTableProperty.Type)
		//字段指针表
		if cacheTableProperty.ChineseType == "指针" {
			pointerTable := strings.Replace(cacheTableProperty.Type, ".", "_", -1)
			pointerTable = ReplaceRight(pointerTable, "_", ".", 1)
			cacheTableProperty.PointerTable = pointerTable
		}
		cacheTableClass.Property = append(cacheTableClass.Property, cacheTableProperty)
	}
	//SqlRowIdName
	idProperty := CacheTableProperty{
		Name:         cacheTableClass.SqlRowIdName,
		Description:  "表id唯一主键",
		Type:         "%Integer",
		ChineseType:  PropertyTypeToChinese("%Integer"),
		SqlFieldName: cacheTableClass.SqlRowIdName,
	}
	cacheTableClass.Property = append([]CacheTableProperty{idProperty}, cacheTableClass.Property...)
	return cacheTableClass
}

// ExtractTableDesc 通过正则提取表描述
func ExtractTableDesc(description string) string {
	var desc string
	reg := regexp.MustCompile(`描述: (?s:(.*?))\n`)
	descArr := reg.FindStringSubmatch(description)
	if len(descArr) > 1 {
		desc = descArr[1]
	}
	if len(desc) > 0 {
		desc = strings.Replace(desc, "\n", "", -1)
	} else {
		desc = description
	}
	return desc
}

func PropertyTypeToChinese(propertyType string) string {
	propertyTypeChinese := map[string]string{
		"%String":          "字符串",
		"%Library.String":  "字符串",
		"%Integer":         "数字",
		"%Library.Integer": "数字",
		"%Boolean":         "布尔型",
		"%Library.Boolean": "布尔型",
		"%Date":            "日期",
		"%Library.Date":    "日期",
		"%Time":            "时间",
		"%Library.Time":    "时间",
		"%List":            "列表",
		"%Library.List":    "列表",
		"%Float":           "浮点型",
		"%Library.Float":   "浮点型",
		"%Numeric":         "数字",
		"%Library.Numeric": "数字",
	}
	if value, isMapContainsKey := propertyTypeChinese[propertyType]; isMapContainsKey {
		//key exist
		return value
	} else {
		//key does not exist
		if strings.Contains(propertyType, "DHC") || strings.Contains(propertyType, "Dr") {
			return "指针"
		}
		return propertyType
	}
}
