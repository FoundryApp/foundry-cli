// +build debug

package cmd

import "github.com/spf13/cobra"

func AddRootFlags(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVarP(&debugFile, "debug-file", "d", "", "path to file where the debug logs are written --d='path/to/file.txt'")
}
