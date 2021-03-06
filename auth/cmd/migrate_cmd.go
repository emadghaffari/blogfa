package cmd

import (
	"blogfa/auth/config"
	"blogfa/auth/client/mysql"
	"blogfa/auth/model"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var migrateCMD = cobra.Command{
	Use:     "migrate",
	Long:    "migrate database strucutures. This will migrate tables",
	Aliases: []string{"m"},
	Run:     Runner.migrate,
}

// migrate database with fake data
func (c *command) migrate(cmd *cobra.Command, args []string) {
	// Current working directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	// read from file
	config.Load(dir + "/config.yaml")

	if err := mysql.Storage.Connect(config.Global); err != nil {
		fmt.Println(err.Error())
		return
	}

	if !config.Global.MYSQL.Automigrate {
		fmt.Println("CHANGE MYSQL AUTO MIGRATE IN CONFIGS")
		return
	}

	err = mysql.Storage.GetDatabase().AutoMigrate(
		model.User{},
		model.Role{},
		model.Permission{},
		model.Provider{},
	)
	if err != nil {
		fmt.Println(err)
	}
}
