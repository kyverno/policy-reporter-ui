package cmd

import (
	"flag"
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kyverno/policy-reporter-ui/pkg/config"
	"github.com/kyverno/policy-reporter-ui/pkg/logging"
)

var configFile string

func newRunCMD() *cobra.Command {
	c := &config.Config{}

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run Policy Reporter UI",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := config.Load(c, configFile)
			if err != nil {
				return err
			}
			logger, err := logging.New(c.Logging)
			if err != nil {
				return fmt.Errorf("failed to setup logger: %w", err)
			}
			if err := config.SetupMemLimit(c); err != nil {
				return fmt.Errorf("failed to setup memlimit: %w", err)
			}

			resolver := config.NewResolver(c)

			serv, err := resolver.Server(cmd.Context())
			if err != nil {
				return err
			}

			logger.Info("Server starts", zap.Int("port", c.Server.Port))

			g := &errgroup.Group{}
			g.Go(serv.Start)

			if c.CRDs.CustomBoard {
				g.Go(func() error {
					informer, err := resolver.CustomBoardInformer()
					if err != nil {
						return err
					}

					stop := make(chan struct{})
					informer.Run(stop)
					logger.Info("custom board informer starts")
					<-stop

					return nil
				})
			}

			return g.Wait()
		},
	}

	// For local usage
	clientcmd.BindOverrideFlags(&c.KubeConfig, cmd.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))

	cmd.Flags().StringVarP(&configFile, "config", "c", "", "target configuration file")
	cmd.Flags().BoolVar(&c.Server.OverwriteHost, "overwrite-host", false, "Overwrite Proxy Host and set Forward Header")
	cmd.Flags().IntVar(&c.Server.Port, "port", 8080, "PolicyReporter UI port")
	cmd.Flags().BoolVar(&c.Server.CORS, "dev", false, "Enable CORS Header for development")
	cmd.Flags().BoolVar(&c.UI.Disabled, "no-ui", false, "Disable the embedded frontend")
	cmd.Flags().BoolVar(&c.Local, "local", false, "use kube config to connect to cluster")
	cmd.Flags().BoolVar(&c.AutoMemoryLimit.Enabled, "auto-memory-enabled", true, "Enable automatic GOMEMLIMIT configuration based on container or system memory.")
	cmd.Flags().Float64Var(&c.AutoMemoryLimit.Ratio, "auto-memory-ratio", 0.9, "The ratio of reserved GOMEMLIMIT memory to the detected maximum container or system memory. Must be greater than 0 and less than or equal to 1.")
	flag.Parse()

	return cmd
}
