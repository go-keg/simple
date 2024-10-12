package migrate

import (
	"context"
	"fmt"
	"log"

	"github.com/go-keg/simple/conf"
	"github.com/go-keg/simple/data"
	"github.com/go-keg/simple/data/ent"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "migrate",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("conf")
		cfg, err := conf.Load(path)
		if err != nil {
			panic(err)
		}
		client, err := data.NewEntClient(cfg)
		if err != nil {
			panic(err)
		}
		defer func(client *ent.Client) {
			err := client.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(client)
		// Run the auto migration tool.
		if err := client.Debug().Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
	},
}
