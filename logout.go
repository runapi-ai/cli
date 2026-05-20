package main

import (
	"github.com/spf13/cobra"
)

func (c *cli) logoutCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove saved API token from config",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}

			if cfg.APIKey == "" {
				c.logf("Not logged in")
				return c.writeJSON(map[string]any{"logged_out": true})
			}

			// Clear auth fields, preserve base_url and other settings
			cfg.APIKey = ""
			cfg.CreatedAt = ""
			if err := saveConfig(cfg); err != nil {
				return err
			}

			configPath, _ := configFilePath()
			c.logf("✓ Logged out. Token removed from %s", configPath)

			return c.writeJSON(map[string]any{"logged_out": true})
		},
	}
}
