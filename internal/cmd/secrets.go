package cmd

import (
	"fmt"

	"github.com/aykay76/ici/internal/secrets"
	"github.com/spf13/cobra"
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manage local secrets for workflows",
	Long: `Manage secrets that will be injected into workflows at runtime.

Secrets are stored locally in ~/.ici/secrets.json and are only available
to workflows run on this machine. They are not synchronized with GitHub.

Examples:
  ici secrets set MY_SECRET my-secret-value
  ici secrets list
  ici secrets remove MY_SECRET`,
}

var secretsSetCmd = &cobra.Command{
	Use:   "set <name> <value>",
	Short: "Store a secret",
	Long:  `Store a secret that will be available to workflows via environment variables.`,
	Args:  cobra.ExactArgs(2),
	RunE:  secretsSet,
}

var secretsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets",
	Long:  `List all stored secret names (values are not displayed for security).`,
	RunE:  secretsList,
}

var secretsRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a secret",
	Long:  `Remove a stored secret.`,
	Args:  cobra.ExactArgs(1),
	RunE:  secretsRemove,
}

func init() {
	rootCmd.AddCommand(secretsCmd)
	secretsCmd.AddCommand(secretsSetCmd)
	secretsCmd.AddCommand(secretsListCmd)
	secretsCmd.AddCommand(secretsRemoveCmd)
}

func secretsSet(cmd *cobra.Command, args []string) error {
	name := args[0]
	value := args[1]

	store := secrets.DefaultStore()
	if err := store.Set(name, value); err != nil {
		return fmt.Errorf("failed to store secret: %w", err)
	}

	fmt.Printf("✓ Secret '%s' stored successfully\n", name)
	return nil
}

func secretsList(cmd *cobra.Command, args []string) error {
	store := secrets.DefaultStore()
	names, err := store.List()
	if err != nil {
		return fmt.Errorf("failed to list secrets: %w", err)
	}

	if len(names) == 0 {
		fmt.Println("No secrets stored")
		return nil
	}

	fmt.Println("Stored secrets:")
	for _, name := range names {
		fmt.Printf("  • %s\n", name)
	}

	return nil
}

func secretsRemove(cmd *cobra.Command, args []string) error {
	name := args[0]

	store := secrets.DefaultStore()
	if err := store.Delete(name); err != nil {
		return fmt.Errorf("failed to remove secret: %w", err)
	}

	fmt.Printf("✓ Secret '%s' removed successfully\n", name)
	return nil
}
