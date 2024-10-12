//go:generate go run -mod=readonly github.com/google/wire/cmd/wire

package main

import (
	"os"

	"github.com/go-keg/keg/contrib/log"
	"github.com/go-keg/keg/contrib/tracing"
	"github.com/go-keg/simple/cmd/migrate"
	"github.com/go-keg/simple/conf"
	_ "github.com/go-keg/simple/data/ent/runtime"
	"github.com/go-keg/simple/job"
	"github.com/go-keg/simple/schedule"
	"github.com/go-kratos/kratos/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/spf13/cobra"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

var (
	// Name is the name of the compiled software.
	Name = "admin"
	// Version is the version of the compiled software.
	Version = "latest"

	id, _ = os.Hostname()

	rootCmd = &cobra.Command{
		Use:     Name,
		Version: Version,
	}
)

func init() {
	rootCmd.AddCommand(migrate.Cmd)
	rootCmd.PersistentFlags().String("conf", "./configs/admin.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("conf")
		cfg, err := conf.Load(path)
		if err != nil {
			panic(err)
		}
		logger := log.NewLoggerFromConfig(cfg.Log, Name,
			log.ServiceName(Name),
			log.ServiceVersion(Version),
			log.ServiceInstanceID(id),
			log.DeploymentEnvironment(os.Getenv("APP_ENV")),
		)
		tracing.SetTracerProvider(cfg.Trace.Endpoint,
			semconv.ServiceName(Name),
			semconv.ServiceVersion(Version),
			semconv.ServiceInstanceID(id),
			semconv.DeploymentEnvironment(os.Getenv("APP_ENV")),
		)
		app, cleanup, err := initApp(logger, cfg)
		if err != nil {
			panic(err)
		}
		defer cleanup()

		// start and wait for stop signal
		if err = app.Run(); err != nil {
			panic(err)
		}
	}
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func newApp(logger klog.Logger, hs *http.Server, job *job.Job, schedule *schedule.Schedule) *kratos.App {
	return kratos.New(
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Logger(logger),
		kratos.Server(hs, job, schedule),
	)
}
