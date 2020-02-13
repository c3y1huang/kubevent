package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Version string
	Build   string
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Printf("Version: %s Build: %s\n", Version, Build)
			return err
		},
	}
}
