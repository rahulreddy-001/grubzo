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
			logger := getLogger()
			defer logger.Sync()

			logger.Info(fmt.Sprintf("grubzo %s (revision %s)", Version, Revision))

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
			defer db.Close()
			logger.Info("database connection was established")

			// FileStorage
			logger.Info("checking file storage...")
			fs, err := getFileStorage(c)
			if err != nil {
				logger.Fatal("failed to setup file storage", zap.Error(err))
			}
			logger.Info("file storage is ok")

			// Repository
			logger.Info("setting up repository...")
			repository, _, err := repository.NewRepository(gormDB, logger, true)
			if err != nil {
				logger.Fatal("failed to initialize repository", zap.Error(err))
			}
			logger.Info("repository was set up")

			// Server
			server, err := newServer(logger, gormDB, repository, fs, c)
			if err != nil {
				logger.Fatal("failed to create server", zap.Error(err))
			}

			// if init {
			// 	logger.Info("data initializing...")

			// 	if err := repo.CreateUserRoles(role.SystemRoleModels()...); err != nil {
			// 		logger.Fatal("failed to init system user roles", zap.Error(err))
			// 	}
			// 	if err := server.SS.RBAC.Reload(); err != nil {
			// 		logger.Fatal("failed to reload rbac", zap.Error(err))
			// 	}

			// 	fid, err := file.GenerateIconFile(server.SS.FileManager, "traq")
			// 	if err != nil {
			// 		logger.Fatal("failed to generate icon file", zap.Error(err))
			// 	}
			// 	u, err := repo.CreateUser(repository.CreateUserArgs{
			// 		Name:       "grubzo",
			// 		Password:   "grubzo",
			// 		Role:       role.Admin,
			// 		IconFileID: fid,
			// 	})
			// 	if err == nil {
			// 		logger.Info("grubzo user was created", zap.Stringer("uid", u.GetID()))
			// 	} else {
			// 		logger.Fatal("failed to init admin user", zap.Error(err))
			// 	}

			// 	logger.Info("data initialization finished")
			// }

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
