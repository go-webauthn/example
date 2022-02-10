package server

import (
	"net"

	"github.com/valyala/fasthttp"

	"github.com/go-webauthn/example/internal/configuration"
	"github.com/go-webauthn/example/internal/handler"
	"github.com/go-webauthn/example/internal/middleware"
)

func Run(config *configuration.Config) (err error) {
	var providers *middleware.Providers

	if providers, err = middleware.NewProviders(config); err != nil {
		return err
	}

	efs := handler.NewEmbeddedFS(handler.EmbeddedFSConfig{
		Prefix:     "public_html",
		IndexFiles: []string{"index.html"},
		TemplatedFiles: map[string]handler.TemplatedEmbeddedFSFileConfig{
			"index.html": {
				Data: struct{ ExternalURL string }{config.ExternalURL.String()},
			},
		},
	}, assets)

	if err = efs.Load(); err != nil {
		return err
	}

	r := middleware.NewRouter(config, providers)

	r.GET("/", middleware.CORS(efs.Handler()))
	r.GET("/{filepath:*}", middleware.CORS(efs.Handler()))

	r.OPTIONS("/debug", middleware.CORS(handler.Nil))
	r.GET("/debug", middleware.CORS(handler.DebugGET))

	r.OPTIONS("/api/info", middleware.CORS(handler.Nil))
	r.GET("/api/info", middleware.CORS(handler.InfoGET))

	r.OPTIONS("/api/login", middleware.CORS(handler.Nil))
	r.POST("/api/login", middleware.CORS(handler.LoginPOST))

	r.OPTIONS("/api/logout", middleware.CORS(handler.Nil))
	r.GET("/api/logout", middleware.Authenticated(middleware.CORS(handler.LogoutGET)))

	r.OPTIONS("/api/u2f/register", middleware.CORS(handler.Nil))
	r.GET("/api/u2f/register", middleware.Authenticated(middleware.CORS(handler.U2FRegisterGET)))
	r.POST("/api/u2f/register", middleware.Authenticated(middleware.CORS(handler.U2FRegisterPOST)))

	r.OPTIONS("/api/webauthn/attestation", middleware.CORS(handler.Nil))
	r.GET("/api/webauthn/attestation", middleware.Authenticated(middleware.CORS(handler.AttestationGET)))
	r.POST("/api/webauthn/attestation", middleware.Authenticated(middleware.CORS(handler.AttestationPOST)))

	r.OPTIONS("/api/webauthn/assertion", middleware.CORS(handler.Nil))
	r.GET("/api/webauthn/assertion", middleware.CORS(handler.AssertionGET))
	r.POST("/api/webauthn/assertion", middleware.CORS(handler.AssertionPOST))

	server := &fasthttp.Server{
		Handler:               r.Handler,
		NoDefaultServerHeader: true,
	}

	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		return err
	}

	return server.Serve(listener)
}
