package item

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func Serve(w http.ResponseWriter, r *http.Request) {
	templ := template.New("")

	err := filepath.Walk("../site/layouts", func(path string, info fs.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			_, err := templ.Parse(path)
			if err != nil {
				log.Err(err).Msg("wrong parse files")
			}
		}

		return err
	})

	log.Print(templ.DefinedTemplates())

	if err != nil {
		panic(err)
	}

	templ.Execute(w, nil)

	// tmp, err := template.ParseFS()("../site/layouts")
	// if err != nil {
	// 	log.Err(err).Msg("kek")
	// 	return
	// }

	// if err := tmp.Execute(w, 0); err != nil {
	// 	log.Err(err).Msg("wrong execute template")
	// }
}
