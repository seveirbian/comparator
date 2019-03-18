package cmd

import (
    "os"
    "github.com/spf13/cobra"
    "github.com/sirupsen/logrus"
)

var logger = logrus.WithField("comparator", "compare")

var rootCmd = &cobra.Command{
    Use:   "comparator",
    Short: "Comparator is a tool that compares file in two folders",
    Long: `A tool that compares files in two folders to see 
if it contains files with same content`,
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    logrus.WithFields(logrus.Fields{
                "err": err, 
                }).Fatal("Fail execute rootCmd")
    os.Exit(1)
  }
}