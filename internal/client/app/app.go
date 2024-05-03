package app

import (
	"fmt"
	"log"
	"os"

	"github.com/mamadeusia/file-transfer-proof/internal/client/service"
	"github.com/spf13/cobra"
)

var (
	serverAddr    string
	directoryPath string
	rootCmd       = &cobra.Command{
		Use:   "client",
		Short: "tool for upload and downloading files",
		Run: func(cmd *cobra.Command, args []string) {
			clientService := service.New(serverAddr, directoryPath)
			if err := clientService.SendFile(); err != nil {
				log.Fatal(err)
			}
		},
	}

	collectionHash string
	indexNum       int32
	downloadSubCmd = &cobra.Command{
		Use: "download",
		Run: func(cmd *cobra.Command, args []string) {
			clientService := service.New(serverAddr, directoryPath)
			if err := clientService.DownloadFile(indexNum, collectionHash); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&serverAddr, "addr", "a", "", "server address")
	rootCmd.Flags().StringVarP(&directoryPath, "file", "f", "", "file path")
	if err := rootCmd.MarkFlagRequired("file"); err != nil {
		log.Fatal(err)
	}
	if err := rootCmd.MarkFlagRequired("addr"); err != nil {
		log.Fatal(err)
	}

	downloadSubCmd.Flags().StringVarP(&serverAddr, "addr", "a", "", "server address")
	downloadSubCmd.Flags().StringVarP(&collectionHash, "collhash", "c", "", "collection hash")
	downloadSubCmd.Flags().Int32VarP(&indexNum, "index", "i", 0, "index file")
	if err := downloadSubCmd.MarkFlagRequired("collhash"); err != nil {
		log.Fatal(err)
	}
	if err := downloadSubCmd.MarkFlagRequired("index"); err != nil {
		log.Fatal(err)
	}
	if err := downloadSubCmd.MarkFlagRequired("addr"); err != nil {
		log.Fatal(err)
	}
	rootCmd.AddCommand(downloadSubCmd)

}
