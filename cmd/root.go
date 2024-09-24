package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func obtainToken() string {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// TODO:
	// viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	token := viper.Get("TOKEN").(string)
	return token
}

func getHighlight() string {
	const url = "https://readwise.io/api/v2/review/"
	token := obtainToken()

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(content))
	return string(content)
}

var rootCmd = &cobra.Command{
	Use:   "readwise-daily",
	Short: "readwise-daily print daily highlights to the command line",
	Run: func(cmd *cobra.Command, args []string) {
		highlight := getHighlight()
		fmt.Println(highlight)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
