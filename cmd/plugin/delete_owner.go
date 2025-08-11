package main

import (
	"github.com/spf13/cobra"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// Defines logic for the delete-owner subcommand
func deleteOwnerCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "delete-owner",
		RunE: deleteOwner,
	}
}

func deleteOwner(_ *cobra.Command, _ []string) error {
	log = logf.Log.WithName("delete-owner")

	log.Info("Deleting owner")
	log.V(1).Info("debug")
	return nil
}
