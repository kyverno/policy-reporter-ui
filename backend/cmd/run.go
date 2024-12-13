package cmd

import (
	"flag"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kyverno/policy-reporter-ui/pkg/config"
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

			resolver := config.NewResolver(c)
			logger := resolver.Logger()

			serv, err := resolver.Server(cmd.Context())
			if err != nil {
				return err
			}

			logger.Info("Server starts", zap.Int("port", c.Server.Port))

			return serv.Start()
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
	flag.Parse()

	return cmd
}
