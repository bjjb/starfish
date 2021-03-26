package cmd

import "github.com/spf13/cobra"

func New(settings ...Setting) *cobra.Command {
	c := &cobra.Command{}
	for _, set := range settings {
		c = set(c)
	}
	return c
}
