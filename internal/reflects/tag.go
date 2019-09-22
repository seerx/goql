package reflects

import (
	"reflect"
	"strings"

	"github.com/seerx/goql/internal/parser/types"

	"github.com/seerx/goql/internal/reflects/validators"
)

const (
	tagDesc   = "desc"
	tagPrefix = "prefix"
	tagLimit  = "limit"
	tagRegexp = "regexp"
	//tagExplodeParams = "explode"

)

//const (
//	fieldPrefix  = "Prefix"
//	fieldExplode = "ExplodeParams"
//)

// GqlTag 解析 tag
type GqlTag struct {
	FieldName   string // 字段名称,使用 json 定义，如果没有则使用 fieldName
	Prefix      string // 所有函数的前缀
	Description string // 说明
	Limit       string // 限制
	Regexp      string //正则表达式限制
}

func (tag *GqlTag) GenerateValidators(typ reflect.Type) []validators.Validator {
	validatorAry := []validators.Validator{}

	if tag.Limit != "" {
		if types.IsInt(typ.Kind()) {
			v := validators.CreateIntegerLimit(tag.FieldName, tag.Limit)
			if v != nil {
				validatorAry = append(validatorAry, v)
			}
		}

		if types.IsFloat(typ.Kind()) {
			v := validators.CreateFloatLimit(tag.FieldName, tag.Limit)
			if v != nil {
				validatorAry = append(validatorAry, v)
			}
		}

		if types.IsString(typ.Kind()) {
			v := validators.CreateStringLimit(tag.FieldName, tag.Limit)
			if v != nil {
				validatorAry = append(validatorAry, v)
			}
		}
	}
	if tag.Regexp != "" {
		if types.IsString(typ.Kind()) {
			v := validators.CreateRegexpValidator(tag.FieldName, tag.Limit)
			if v != nil {
				validatorAry = append(validatorAry, v)
			}
		}
	}
	return validatorAry
}

// ParseTag 解析 Tag
func ParseTag(field *reflect.StructField) *GqlTag {
	gTag := &GqlTag{}
	tag := field.Tag

	gTag.FieldName = parseFieldName(&tag, field)
	gqlTag := tag.Get("gql")
	if gqlTag != "" {
		mp := map[string]string{}
		ary := strings.Split(gqlTag, ",")
		for _, item := range ary {
			sub := strings.Split(item, "=")
			if len(sub) == 1 {
				mp[sub[0]] = sub[0]
			}
			if len(sub) == 2 {
				mp[sub[0]] = sub[1]
			}
		}

		gTag.Prefix = mp[tagPrefix] //  parseGqlPrefix(mp)
		gTag.Description = mp[tagDesc]
		gTag.Regexp = mp[tagRegexp]
		gTag.Limit = mp[tagLimit]
	}

	return gTag
}

// parseGqlPrefix 解析函数前缀
//func parseGqlPrefix(pairs map[string]string) string {
//	return pairs[tagPrefix]
//}

// parseFieldName 解析字段名称
func parseFieldName(tag *reflect.StructTag, field *reflect.StructField) string {
	name := tag.Get("json")
	if name == "" {
		name = field.Name
	}
	if name == "-" {
		// json 中的忽略
		return ""
	}
	ary := strings.Split(name, ",")
	if len(ary) == 1 {
		return name
	}

	for _, item := range ary {
		if item != "omitempty" {
			return item
		}
	}

	return ""
}
