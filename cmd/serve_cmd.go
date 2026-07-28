package cmd

import (
	"context"
	"fmt"
	"grubzo/internal/repository"
	"grubzo/internal/services"
	"grubzo/internal/utils/gormzap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Server struct {
	L      *zap.Logger
	SS     *services.Services
	Router *gin.Engine
	HTTP   *http.Server
	Repo   repository.Repository
}

func serveCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "serve",
		Short: "Serve grubzo API",
		Run: func(_ *cobra.Command, _ []string) {
			// Logger
			logger, closeLogger := getLogger()
			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_ = closeLogger(ctx)
			}()
			logger = logger.With(zap.String("instance", c.Instance))

			logger.Info(fmt.Sprintf("grubzo %s (revision %s)", Version, Revision))

			// Tracer
			shutdown, err := initTracer("GrubzoServer", c, logger.Named("tracer"))
			if err != nil {
				logger.Fatal("failed to initilize tracer", zap.Error(err))
			}
			logger.Info("tracer initilized successfully")
			defer func() {
				if err := shutdown(context.Background()); err != nil {
					logger.Warn("failed to shutdown tracer", zap.Error(err))
				}
			}()

			//Redis
			logger.Info("connecting redis...")
			redis, err := getRedisClient(c)
			if err != nil {
				logger.Fatal("failed to connect redis", zap.Error(err))
			}
			logger.Info("redis connection established")

			// Database
			logger.Info("connecting database...")
			gormDB, err := getDatabase(c)
			if err != nil {
				logger.Fatal("failed to connect database", zap.Error(err))
			}
			gormDB.Logger = gormzap.New(logger.Named("gorm"))
			db, err := gormDB.DB()
			if err != nil {
				logger.Fatal("failed to get *sql.DB", zap.Error(err))
			}
			defer func() {
				if err := db.Close(); err != nil {
					logger.Warn("failed to close database", zap.Error(err))
				}
			}()
			logger.Info("database connection established")

			// FileStorage
			logger.Info("checking file storage...")
			fs, err := getFileStorage(c)
			if err != nil {
				logger.Fatal("failed to setup file storage", zap.Error(err))
			}
			logger.Info("file storage is ok")

			// Repository
			logger.Info("setting up repository...")
			repository, _, err := repository.NewRepository(gormDB, redis, logger, true)
			if err != nil {
				logger.Fatal("failed to initialize repository", zap.Error(err))
			}
			logger.Info("repository was set up")

			// Server
			server, err := newServer(logger, gormDB, redis, repository, fs, c)
			if err != nil {
				logger.Fatal("failed to create server", zap.Error(err))
			}

			go func() {
				if err := server.Start(fmt.Sprintf(":%d", c.App.Port)); err != nil {
					logger.Info("shutting down the server")
				}
			}()

			logger.Info("grubzo started")
			waitSIGINT()
			logger.Info("grubzo shutting down...")

			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.ShutdownTimeout)*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				logger.Warn("abnormal shutdown", zap.Error(err))
			}
			logger.Info("grubzo shutdown")
		},
	}
	return &cmd
}

func waitSIGINT() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	signal.Stop(quit)
	close(quit)
	for range quit {
		continue
	}
}
