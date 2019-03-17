package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/seveirbian/comparator/compare"
    "github.com/sirupsen/logrus"
)

var logger = logrus.WithField("cmd", "compare")

var compareUsage = `Usage:  comparator compare DIR1 DIR2`

func init() {
    rootCmd.AddCommand(compareCmd)
    compareCmd.SetUsageTemplate(compareUsage)
}   

var compareCmd = &cobra.Command{
    Use:   "compare",
    Short: "Compares files in two folders to see if it contains files with same content",
    Long:  `Compares files in two folders to see if it contains files with same content`,
    Args:  cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        comparator, err := compare.Init(args[0], args[1])
        if err != nil {
            logger.Fatal("Fail to init...")
        }

        err = comparator.Compare()
        if err != nil {
            logger.Fatal("Fail to compare...")
        }

        fmt.Printf("Common size: %d bytes\n", comparator.CommonFilesSize)
    },
}