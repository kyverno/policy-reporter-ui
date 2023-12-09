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

			serv := resolver.Server(cmd.Context())

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
	flag.Parse()

	return cmd
}
