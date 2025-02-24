package api

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/theplant/luhn"
	"github.com/vysogota0399/gophermart_portal/internal/api/logging"
	"github.com/vysogota0399/gophermart_portal/internal/api/models"
	"github.com/vysogota0399/gophermart_portal/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type HTTPServer struct {
	Router *Router
	srv    *http.Server
}

type Controller interface {
	CreateRoutes(*Router) []*Route
}

type Route struct {
	Handler gin.HandlerFunc
	Path    string
	Method  string
	Meta    map[string]bool
}

var AuthorizationRequire = "authorizetion_require"

func (r *Route) IsAuthorizationRequire() bool {
	_, ok := r.Meta[AuthorizationRequire]
	return ok
}

type Router struct {
	routes map[string]map[string]*Route
	router *gin.Engine
}

type AuthorizationService interface {
	Call(ctx context.Context, req *http.Request) (*models.Session, error)
}

var luhnableNumber validator.Func = func(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	number, err := strconv.Atoi(v)
	if err != nil {
		return false
	}

	return luhn.Valid(number)
}

func NewRouter(ctrs []Controller, auth AuthorizationService, lg *logging.ZapLogger) *Router {
	r := &Router{router: gin.Default(), routes: make(map[string]map[string]*Route)}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("luhnablenumber", luhnableNumber)
	}

	r.router.Use(
		Authorize(r, auth, lg),
		AddJSONContentType(),
	)

	for _, c := range ctrs {
		for _, route := range c.CreateRoutes(r) {
			r.AddRoute(route)
		}
	}

	return r
}

func (router *Router) AddRoute(r *Route) {
	paths, ok := router.routes[r.Method]
	if !ok {
		paths = make(map[string]*Route)
		router.routes[r.Method] = paths
	}
	paths[r.Path] = r

	switch r.Method {
	case http.MethodPost:
		router.router.POST(r.Path, r.Handler)
	case http.MethodGet:
		router.router.GET(r.Path, r.Handler)
	}
}

func (router *Router) Find(r *http.Request) *Route {
	paths, ok := router.routes[r.Method]
	if !ok {
		return nil
	}

	path := paths[r.URL.Path]
	return path
}

func NewHTTPServer(lc fx.Lifecycle, cfg *config.Config, r *Router) *HTTPServer {
	s := &HTTPServer{
		srv:    &http.Server{Addr: cfg.HTTPAddress, Handler: r.router, ReadHeaderTimeout: time.Minute},
		Router: r,
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return s.Start(ctx)
			},
			OnStop: func(ctx context.Context) error {
				return s.Shutdown(ctx)
			},
		},
	)

	return s
}

type TestHTTPServer struct {
	Srv    *httptest.Server
	router *Router
}

func NewTestHTTPServer(address string, r *Router) *TestHTTPServer {
	return &TestHTTPServer{
		Srv:    httptest.NewServer(r.router),
		router: r,
	}
}

func (s *HTTPServer) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return err
	}

	go s.srv.Serve(ln)

	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

var CurrentUserIDKey string = "current_user"

func Authorize(r *Router, authorizer AuthorizationService, lg *logging.ZapLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		route := r.Find(c.Request)
		if route == nil {
			lg.ErrorCtx(c, "path not found", zap.String("path", c.Request.URL.Path))
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if !route.IsAuthorizationRequire() {
			c.Next()
			return
		}

		sess, err := authorizer.Call(c, c.Request)
		if err != nil {
			lg.ErrorCtx(c, "authorization failed", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		lg.DebugCtx(c, "authorized user", zap.Any("session", sess))
		c.Set(CurrentUserIDKey, sess.Sub)
		c.Next()
	}
}

func AddJSONContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Type") == "application/json" {
			c.Writer.Header().Set("Content-Type", "application/json")
		}

		c.Next()
	}
}
