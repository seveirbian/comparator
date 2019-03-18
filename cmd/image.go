package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/seveirbian/comparator/image"
)

var imageUsage = `Usage:  comparator image IMAGENAME:TAG IMAGENAME:TAG`

func init() {
    rootCmd.AddCommand(imageCmd)
    imageCmd.SetUsageTemplate(imageUsage)
}   

var imageCmd = &cobra.Command{
    Use:   "image",
    Short: "Compares files in two images' fs to see if it contains files with same content",
    Long:  `Compares files in two images' fs to see if it contains files with same content`,
    Args:  cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        comparator, err := image.Init(args[0], args[1])
        if err != nil {
            logger.Fatal("Fail to init...")
        }

        err = comparator.Compare()
        if err != nil {
            logger.Fatal("Fail to image...")
        }

        fmt.Printf("Common size: %d bytes\n", comparator.Comparator.CommonFilesSize)
    },
}