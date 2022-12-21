package template

import (
	"bytes"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/app"
	"github.com/vivi-app/vivi/extensions/passport"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/tester"
	"strings"
	"text/template"
)

type funcs struct {
	Search,
	Prepare,
	Stream,
	Download,
	Test string
}

type fields struct {
	Display string
	About   string
	Layers  string
}

type metaInfo struct {
	Fields   *fields
	App      string
	Passport *passport.Passport
	Fn       *funcs
}

var meta *metaInfo

func init() {
	meta = &metaInfo{}

	meta.App = app.Name
	meta.Fields = &fields{
		Display: scraper.FieldDisplay,
		About:   scraper.FieldAbout,
		Layers:  scraper.FieldLayers,
	}
	meta.Fn = &funcs{
		Search:   scraper.FunctionSearch,
		Prepare:  scraper.FunctionPrepare,
		Stream:   scraper.FunctionStream,
		Download: scraper.FunctionDownload,
		Test:     tester.FunctionTest,
	}
}

func execTemplate(tmpl string) []byte {
	if strings.Contains(tmpl, "template:skip") {
		return []byte(tmpl)
	}

	parsed := lo.Must(template.New("").Parse(tmpl))
	var b bytes.Buffer
	lo.Must0(parsed.Execute(&b, meta))

	return b.Bytes()
}
