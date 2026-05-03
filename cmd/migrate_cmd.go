package cmd

import (
	"grubzo/internal/migration"
	"log"

	"github.com/spf13/cobra"
)

func migrateCommand() *cobra.Command {
	var dropDB bool
	cmd := cobra.Command{
		Use:   "migrate",
		Short: "Execute database schema migration only",
		RunE: func(_ *cobra.Command, _ []string) error {
			log.Println("[migrate] start")
			engine, err := getDatabase(c)
			if err != nil {
				log.Printf("[migrate] getDatabase error: %v", err)
				return err
			}
			log.Println("[migrate] connected to DB engine")

			sqlDB, err := engine.DB()
			if err != nil {
				log.Printf("[migrate] engine.DB error: %v", err)
				return err
			}
			defer func() {
				_ = sqlDB.Close()
				log.Println("[migrate] closed sql DB")
			}()

			if dropDB {
				log.Println("[migrate] reset flag true: dropping all tables")
				if err := migration.DropAll(engine); err != nil {
					log.Printf("[migrate] DropAll err: %v", err)
					return err
				}
				log.Println("[migrate] DropAll success")
			}

			init, err := migration.Migrate(engine)
			if err != nil {
				log.Printf("[migrate] migration failed: %v", err)
				return err
			}
			log.Printf("[migrate] migration finished, init=%v", init)
			return nil
		},
	}
	flags := cmd.Flags()
	flags.BoolVar(&dropDB, "reset", false, "whether to truncate database (drop all tables)")
	return &cmd
}
