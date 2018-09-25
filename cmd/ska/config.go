package main

import "os"

var templatesDir string

func init() {
	rootCmd.PersistentFlags().StringVarP(&templatesDir, "templates", "t", os.Getenv("HOME")+"/.ska", "Templates dir")
}
