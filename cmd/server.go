package cmd

import (
	"context"
	"grubzo/internal/config"
	"grubzo/internal/repository"
	"grubzo/internal/router"
	"grubzo/internal/services"
	"grubzo/internal/utils/storage"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

func newServer(logger *zap.Logger, db *gorm.DB, repository *repository.Repository, fs storage.FileStorage, config *config.Config) (*Server, error) {
	services, err := services.Setup(logger, db, repository, fs, config)
	if err != nil {
		logger.Fatal("failed to initialize services", zap.Error(err))
	}
	gin := router.Setup(logger, db, repository, services, config)
	server := &Server{
		L:      logger,
		Router: gin,
		// SS:     services,
	}
	return server, nil
}

func (s *Server) Start(address string) error {
	s.HTTP = &http.Server{
		Addr:    address,
		Handler: s.Router,
	}
	return s.HTTP.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		err := s.HTTP.Shutdown(ctx)
		s.L.Info("Router shutdown")
		return err
	})
	return eg.Wait()
}
