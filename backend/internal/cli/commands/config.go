package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/lockw1n/time-logger/internal/cli/config"
)

var errUnknownConfigKey = errors.New("unknown config key (valid: base_url, default_company, default_activity, email, ca_cert, insecure)")

var configKeys = []string{"base_url", "default_company", "default_activity", "email", "ca_cert", "insecure"}

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage tl configuration",
	}

	cmd.AddCommand(newConfigSetCmd(), newConfigGetCmd(), newConfigListCmd())

	return cmd
}

func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			key, value := args[0], args[1]
			switch key {
			case "base_url":
				cfg.BaseURL = value
			case "default_company":
				id, err := strconv.ParseUint(value, 10, 64)
				if err != nil {
					return fmt.Errorf("default_company must be a positive integer, got %q", value)
				}
				cfg.DefaultCompanyID = id
			case "default_activity":
				cfg.DefaultActivity = value
			case "email":
				cfg.Email = value
			case "ca_cert":
				cfg.CACertPath = value
			case "insecure":
				b, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("insecure must be true or false, got %q", value)
				}
				cfg.Insecure = b
			default:
				return fmt.Errorf("%q: %w", key, errUnknownConfigKey)
			}

			return config.Save(cfg)
		},
	}
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Print a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			value, err := configValue(cfg, args[0])
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), value)
			return nil
		},
	}
}

func newConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Print all configuration values",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			for _, key := range configKeys {
				value, err := configValue(cfg, key)
				if err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", key, value)
			}

			return nil
		},
	}
}

func configValue(cfg config.Config, key string) (string, error) {
	switch key {
	case "base_url":
		return cfg.BaseURL, nil
	case "default_company":
		if cfg.DefaultCompanyID == 0 {
			return "", nil
		}
		return strconv.FormatUint(cfg.DefaultCompanyID, 10), nil
	case "default_activity":
		return cfg.DefaultActivity, nil
	case "email":
		return cfg.Email, nil
	case "ca_cert":
		return cfg.CACertPath, nil
	case "insecure":
		return strconv.FormatBool(cfg.Insecure), nil
	default:
		return "", fmt.Errorf("%q: %w", key, errUnknownConfigKey)
	}
}
