package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/yuridevx/app-example/cmd/example/wired"
	"os"
	"os/signal"
	"sync"
)

func run() *cobra.Command {
	var cmd = &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {
			// Run with graceful shutdown
			appObjects, err := wired.InitApp()
			if err != nil {
				panic(err)
			}
			app := appObjects.Builder.Build()
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			//ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			wg := &sync.WaitGroup{}
			app.Run(ctx, wg)
			appObjects.Health.OnAppStarted()
			wg.Wait()
			appObjects.Logger.Info("exiting gracefully")
		},
	}

	return cmd
}
