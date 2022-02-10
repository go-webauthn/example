package handler

import (
	"embed"
	"html/template"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-webauthn/example/internal/middleware"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

var templates map[string]*template.Template

func NewEmbedFS(embedFS embed.FS, prefix string, indexFiles []string, templateFiles []string) (handler middleware.RequestHandler) {

	if len(indexFiles) == 0 {
		indexFiles = []string{"index.html"}
	}

	efs, _ := fs.Sub(embedFS, prefix)

	content := http.FS(efs)

	templates = make(map[string]*template.Template, len(templateFiles))

	log := zap.L()

	for _, name := range templateFiles {
		name = assetPath(name)

		log.Debug("loading template", zap.String("name", name))

		file, err := content.Open(name)
		if err != nil {
			log.Error("failed to load template", zap.String("name", name), zap.Error(err))

			continue
		}

		stat, err := file.Stat()
		if err != nil {
			log.Error("failed to stat template", zap.String("name", name), zap.Error(err))

			continue
		}

		data := make([]byte, stat.Size())

		_, err = file.Read(data)
		if err != nil {
			log.Error("failed to read template", zap.String("name", name), zap.Error(err))

			continue
		}

		t, err := template.New(name).Parse(string(data))
		if err != nil {
			log.Error("failed to parse template", zap.String("name", name), zap.Error(err))

			continue
		}

		log.Debug("loaded template successfully", zap.String("name", name))

		templates[name] = t
	}

	return func(ctx *middleware.RequestCtx) {
		requestPath := assetPath(string(ctx.RequestURI()))

		if handleTemplate(ctx, requestPath) {
			return
		}

		found := false

		file, err := content.Open(requestPath)
		if err != nil {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)

			return
		}

		stat, err := file.Stat()
		if err != nil {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)

			return
		}

		if stat.IsDir() {
			for _, indexFile := range indexFiles {
				indexPath := path.Join(requestPath, indexFile)

				if handleTemplate(ctx, indexPath) {
					return
				}

				file, err = content.Open(indexPath)
				if err != nil {
					continue
				}

				stat, err = file.Stat()
				if err != nil || stat.IsDir() {
					continue
				}

				found = true
			}
		} else {
			found = true
		}

		if found {
			setMIME(ctx, stat.Name())

			ctx.SetBodyStream(file, int(stat.Size()))

			return
		} else {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)

			return
		}
	}
}

func setMIME(ctx *middleware.RequestCtx, name string) {
	ext := strings.ToLower(filepath.Ext(name))
	contentType := mime.TypeByExtension(ext)
	if contentType != "" {
		ctx.SetContentType(contentType)
	}
}

func handleTemplate(ctx *middleware.RequestCtx, path string) (handled bool) {
	ctx.Log.Debug("checking if requested path is a templated asset")

	if t, ok := templates[path]; ok {
		ctx.Log.Debug("requested asset is templated")

		err := t.Execute(ctx.Response.BodyWriter(), struct{ ExternalURL string }{ctx.Config.ExternalURL.String()})
		if err == nil {
			ctx.Log.Debug("rendered template")

			setMIME(ctx, path)

			return true
		}

		ctx.Log.Error("error occurred rendering template", zap.Error(err))
	}

	ctx.Log.Debug("requested asset is not templated")

	return false
}

func assetPath(name string) (path string) {
	if !strings.HasPrefix(name, "/") {
		return "/" + name
	}

	return name
}
