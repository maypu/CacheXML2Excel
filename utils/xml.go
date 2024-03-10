package utils

import (
	"github.com/beevik/etree"
	"log"
	"regexp"
	"strings"
)

type CacheTableClass struct {
	Name	string
	Description string
	ClassType	string
	SqlTableName	string
	Property    []CacheTableProperty
}

type CacheTableProperty struct {
	Name         string
	Description  string
	Type		string
	ChineseType	string
	SqlFieldName string
	PointerTable	string	//类型为指针时，指向的表名
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
	cacheTableClass.Name = cacheClass.SelectAttrValue("name","")
	//计算表名
	cacheTableClassName := strings.Replace(cacheTableClass.Name, ".", "_", -1)
	cacheTableClassName = ReplaceRight(cacheTableClassName,"_",".",1)
	cacheTableClass.SqlTableName = cacheTableClassName
	//描述
	if classDescription := cacheClass.SelectElement("Description"); classDescription != nil {
		cacheTableClass.Description = classDescription.Text()
	}
	//计算表描述
	extractTableDesc := ExtractTableDesc(cacheTableClass.Description)
	if extractTableDesc != "" {
		cacheTableClass.Description = extractTableDesc
	}

	for _, property := range cacheClass.SelectElements("Property") {

		if cardinality := property.SelectElement("Cardinality"); cardinality != nil {
			if cardinality.Text() == "children" {		//父表中指向子表的字段可以不算
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
			cacheTableProperty.Description = ReplaceRight(propertyDescription.Text(), "\n", "", 1)
		} else {
			cacheTableProperty.Description = cacheTableProperty.SqlFieldName
		}
		//字段类型中文名
		cacheTableProperty.ChineseType = PropertyTypeToChinese(cacheTableProperty.Type)
		//字段指针表
		if cacheTableProperty.ChineseType == "指针" {
			pointerTable := strings.Replace(cacheTableProperty.Type, ".", "_", -1)
			pointerTable = ReplaceRight(pointerTable,"_",".",1)
			cacheTableProperty.PointerTable = pointerTable
		}
		cacheTableClass.Property = append(cacheTableClass.Property, cacheTableProperty)
	}
	return cacheTableClass
}

// ExtractTableDesc 通过正则提取表描述
func ExtractTableDesc(description string) string {
	var desc string
	reg := regexp.MustCompile(`描述: (?s:(.*?))\n`)
	descArr := reg.FindStringSubmatch(description)
	if len(descArr)>1 {
		desc = descArr[1]
	}
	desc = strings.Replace(desc,"\n", "", -1)
	return desc
}

func PropertyTypeToChinese(propertyType string) string {
	propertyTypeChinese := map[string] string {
		"%String": "字符串",
		"%Library.String": "字符串",
		"%Integer": "数字",
		"%Library.Integer": "数字",
		"%Boolean": "布尔",
		"%Library.Boolean": "布尔",
		"%Date": "日期",
		"%Library.Date": "日期",
		"%Time": "时间",
		"%Library.Time": "时间",
		"%List": "列表",
		"%Library.List": "列表",
		"%Float": "浮点数",
		"%Library.Float": "浮点数",
	}
	if value, isMapContainsKey := propertyTypeChinese[propertyType]; isMapContainsKey {
		//key exist
		return value
	} else {
		//key does not exist
		if strings.Contains(propertyType, "DHC") {
			return "指针"
		}
		return propertyType
	}
}