package handler

import (
	"bytes"
	"embed"
	"errors"
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

func NewEmbeddedFS(config EmbeddedFSConfig, embedFS embed.FS) (efs EmbeddedFS) {
	efs = EmbeddedFS{
		config: config,
	}

	if len(efs.config.IndexFiles) == 0 {
		efs.config.IndexFiles = []string{"index.html"}
	}

	fsfs, _ := fs.Sub(embedFS, efs.config.Prefix)

	efs.filesystem = http.FS(fsfs)

	return efs
}

func (e *EmbeddedFS) Load() (err error) {
	if e.loaded {
		return errors.New("templates already loaded")
	}

	log := zap.L()

	var (
		name     string
		file     http.File
		fileInfo fs.FileInfo
	)

	for tName, templatedFile := range e.config.TemplatedFiles {
		name = e.assetPath(tName)

		log.Debug("loading template", zap.String("name", name))

		if file, err = e.filesystem.Open(name); err != nil {
			log.Error("failed to load template", zap.String("name", name), zap.Error(err))

			return err
		}

		if fileInfo, err = file.Stat(); err != nil {
			log.Error("failed to stat template", zap.String("name", name), zap.Error(err))

			return err
		}

		data := make([]byte, fileInfo.Size())

		if _, err = file.Read(data); err != nil {
			log.Error("failed to read template", zap.String("name", name), zap.Error(err))

			return err
		}

		if templatedFile.template, err = template.New(name).Parse(string(data)); err != nil {
			log.Error("failed to parse template", zap.String("name", name), zap.Error(err))

			return err
		}

		log.Debug("parsed template successfully", zap.String("name", name))

		log.Debug("executing pre-rendered template", zap.String("name", name), zap.Error(err))

		out := bytes.Buffer{}

		if err = templatedFile.template.Execute(&out, templatedFile.Data); err != nil {
			log.Error("failed to execute pre-rendered template", zap.String("name", name), zap.Error(err))

			return err
		}

		templatedFile.rendered = out.Bytes()

		e.config.TemplatedFiles[tName] = templatedFile
	}

	return nil
}

func (e EmbeddedFS) assetPath(name string) (path string) {
	if !strings.HasPrefix(name, "/") {
		return "/" + name
	}

	return name
}

func (e EmbeddedFS) setMIME(ctx *middleware.RequestCtx, path string) {
	ext := strings.ToLower(filepath.Ext(path))
	contentType := mime.TypeByExtension(ext)
	if contentType != "" {
		ctx.SetContentType(contentType)
	}
}

func (e EmbeddedFS) template(ctx *middleware.RequestCtx, path string) (handled bool) {
	ctx.Log.Debug("checking if requested path is a templated asset", zap.String("name", path))

	if t, ok := e.config.TemplatedFiles[path[1:]]; ok {
		ctx.Log.Debug("requested asset is templated", zap.String("name", path))

		e.setMIME(ctx, path)

		ctx.SetBody(t.rendered)

		return true
	}

	ctx.Log.Debug("requested asset is not templated")

	return false
}

func (e *EmbeddedFS) Handler() (handler middleware.RequestHandler) {
	return func(ctx *middleware.RequestCtx) {
		requestPath := e.assetPath(string(ctx.RequestURI()))

		if e.template(ctx, requestPath) {
			return
		}

		found := false

		file, err := e.filesystem.Open(requestPath)
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
			ctx.Log.Debug("requested asset is directory", zap.String("name", stat.Name()))

			for _, indexFile := range e.config.IndexFiles {
				indexPath := path.Join(requestPath, indexFile)

				ctx.Log.Debug("checking index file", zap.String("name", indexPath))
				if e.template(ctx, indexPath) {
					return
				}

				if file, err = e.filesystem.Open(indexPath); err != nil {
					continue
				}

				if stat, err = file.Stat(); err != nil || stat.IsDir() {
					continue
				}

				found = true
			}
		} else {
			found = true
		}

		if found {
			e.setMIME(ctx, stat.Name())

			ctx.SetBodyStream(file, int(stat.Size()))

			return
		} else {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)

			return
		}
	}
}

type EmbeddedFS struct {
	filesystem http.FileSystem
	config     EmbeddedFSConfig
	loaded     bool
}

type EmbeddedFSConfig struct {
	Prefix         string
	IndexFiles     []string
	TemplatedFiles map[string]TemplatedEmbeddedFSFileConfig
}

type TemplatedEmbeddedFSFileConfig struct {
	Data     interface{}
	template *template.Template
	rendered []byte
}
