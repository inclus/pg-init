package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/inclus/pg-init/pkg/database"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "pg-init",
	Short: "wait for database to be ready",
	Long:  `wait for database to be ready`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("postgres", database.BuildConnectionString())
		if err != nil {
			log.Fatalf(fmt.Errorf("failed to create connection %w", err).Error())
		}

		err = database.WaitForWorkingConnection(db)
		if err != nil {
			log.Fatalf(fmt.Errorf("failed to check connection %w", err).Error())
		}

		err = database.SetupPostgisExtension(db)
		if err != nil {
			log.Fatalf(fmt.Errorf("failed to setup postgis extension %w", err).Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pg-init.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", "5432")
	viper.SetDefault("db.database", "postgres")
	viper.SetDefault("db.user", "postgres")
	viper.SetDefault("db.password", "")
	viper.SetDefault("db.extra", "sslmode=disable")
	viper.SetDefault("retry.attempts", 10)

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
}
