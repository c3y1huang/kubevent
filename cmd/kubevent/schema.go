package main

import (
	"github.com/spf13/cobra"
)

func NewSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Schema, show cluster API resources info",
		Run: func(cmd *cobra.Command, args []string) {
			for _, t := range scheme.AllKnownTypes() {
				println(t.String())
			}
		},
	}

	return cmd
}
