package mail

import (
	"bytes"
	"html/template"
	"path"
	"reflect"
	"time"

	"labix.org/v2/mgo/bson"
)

var TemplateDir string

func RenderMail(templatePath string, body map[string]interface{}) (data string, err error) {
	funcM := template.FuncMap{
		"formatTime": func(t interface{}) string {
			switch i := t.(type) {
			case int64:
				return time.Unix(i, 0).Format("2006-01-02 15:04:05")
			case bson.MongoTimestamp:
				return time.Unix(int64(i)/1000000000, 0).Format("2006-01-02 15:04:05")
			case time.Time:
				return i.Format("2006-01-02 15:04:05")
			default:
			}
			return "unknown type: " + reflect.TypeOf(t).Name()
		},
	}

	tempName := path.Base(templatePath)
	tpl, err := template.New(tempName).Funcs(funcM).ParseFiles(path.Join(TemplateDir, templatePath))

	if err != nil {
		return
	}

	buf := bytes.Buffer{}
	err = tpl.Execute(&buf, body)
	data = buf.String()
	return
}
