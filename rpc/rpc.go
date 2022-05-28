package rpc

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	rpcmw "github.com/crypto-pe/backend/rpc/middleware"

	"github.com/crypto-pe/backend/config"
	"github.com/crypto-pe/backend/proto"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog"
)

type RPC struct {
	Config *config.Config
	Log    zerolog.Logger

	HTTP *http.Server

	JWTAuth *jwtauth.JWTAuth

	running   int32
	startTime time.Time
}

func NewRPC(cfg *config.Config, logger zerolog.Logger) (*RPC, error) {
	httpServer := &http.Server{
		Addr:              cfg.Service.Listen,
		ReadTimeout:       45 * time.Second,
		WriteTimeout:      45 * time.Second,
		IdleTimeout:       45 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
	s := &RPC{
		Config:  cfg,
		Log:     logger.With().Str("ps", "rpc").Logger(),
		JWTAuth: jwtauth.New("HS256", []byte(cfg.Auth.JWTSecret), nil),
		HTTP:    httpServer,
	}
	return s, nil
}

func (s *RPC) Run(ctx context.Context) error {
	if s.IsRunning() {
		return fmt.Errorf("rpc: already running")
	}

	s.Log.Info().Str("op", "run").Msgf("-> rpc: listening on %s", s.HTTP.Addr)

	atomic.StoreInt32(&s.running, 1)
	defer atomic.StoreInt32(&s.running, 0)

	// Setup HTTP server handler
	s.HTTP.Handler = s.handler()
	// Handle stop signal to ensure clean shutdown
	go func() {
		<-ctx.Done()
		s.Stop(context.Background())
	}()

	// Start the http server and serve!
	err := s.HTTP.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *RPC) GetLogger(ctx context.Context) *zerolog.Logger {
	lg := httplog.LogEntry(ctx)
	return &lg
}

func (s *RPC) Stop(timeoutCtx context.Context) {
	if !s.IsRunning() || s.IsStopping() {
		return
	}
	atomic.StoreInt32(&s.running, 2)

	s.Log.Info().Str("op", "stop").Msg("-> rpc: stopping..")
	s.HTTP.Shutdown(timeoutCtx)
	s.Log.Info().Str("op", "stop").Msg("-> rpc: stopped.")
}

func (s *RPC) IsRunning() bool {
	return atomic.LoadInt32(&s.running) == 1
}

func (s *RPC) IsStopping() bool {
	return atomic.LoadInt32(&s.running) == 2
}

func (s *RPC) handler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.NoCache)
	// r.Use(honeybadger.Handler)
	r.Use(middleware.Heartbeat("/ping"))
	// HTTP request logger
	r.Use(httplog.RequestLogger(s.Log))
	// Timeout any request after 28 seconds as Cloudflare has a 30 second limit anyways.
	r.Use(middleware.Timeout(28 * time.Second))

	// Rate limiting
	r.Use(httprate.LimitByIP(200, 1*time.Minute))
	// CORS
	r.Use(s.corsHandler())

	// Quick pages
	r.Use(middleware.PageRoute("/", http.HandlerFunc(indexHandler)))
	r.Use(middleware.PageRoute("/favicon.ico", http.HandlerFunc(stubHandler(""))))

	// // Seek and verify JWT tokens, and put on request context
	r.Use(jwtauth.Verifier(s.JWTAuth))

	// // Session middleware
	r.Use(rpcmw.Session)

	// // Access control
	r.Use(rpcmw.AccessControl)
	// Mount rpc endpoints
	rpcHandler := proto.NewAPIServer(s)
	// r.Handle("/rpc/ArcadeumAPI/*", chi.Chain(middleware.PathRewrite("/rpc/ArcadeumAPI/", "/rpc/API/")).Handler(rpcHandler))
	r.Post("/rpc/*", rpcHandler.ServeHTTP)

	return r
}

func (s *RPC) corsHandler() func(next http.Handler) http.Handler {
	// CORS options for trusted https://*.sequence.app apps, where we allow
	// authorization headers to pass.

	// TODO: certain endpoints would be sequence.app only, but other endpoints
	// should be usable by dapps like skyweaver

	trustedOrigins := []string{
		// add origins
	}
	corsTrustedOptions := cors.Options{
		AllowedOrigins:   trustedOrigins,
		AllowedMethods:   []string{"HEAD", "GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Release"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           600,
	}

	// TODO: in case a handler panics, or returns 500, we should still return
	// the cors origin header, so we don't confuse the end user with cors error
	// but really its issue with the endpoint request

	// For local dev mode, allow all traffic
	if s.Config.Mode == config.DevelopmentMode {
		corsTrustedOptions.AllowOriginFunc = func(r *http.Request, origin string) bool {
			return true
		}
		return cors.Handler(corsTrustedOptions)
	}

	// CORS options for third-party apps, where we block authorization headers.
	corsThirdPartyOptions := cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Release"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false, // IMPORTANT! must be false
		MaxAge:           600,
	}

	// Here we use the RouteHeaders middleware to split the request paths depending
	// on the Origin request header value.
	return middleware.RouteHeaders().
		RouteAny("Origin", trustedOrigins, cors.Handler(corsTrustedOptions)).
		Route("Origin", "*", cors.Handler(corsThirdPartyOptions)).
		Handler
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("."))
}

func stubHandler(respBody string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(respBody))
	})
}
