package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
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

func getHighlight() readwiseHighlight {
	const url = "https://readwise.io/api/v2/review/"
	token := obtainToken()

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("Failed to establish the HTTP request")
		log.Error(err.Error())
		panic(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	res, err := client.Do(req)
	if err != nil {
		log.Error("Connection the the Readwise API failed")
		panic(err)
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	var obj JSONPayload
	err = json.Unmarshal(content, &obj)
	if err != nil {
		panic(err)
	}
	highlights := obj.Highlights
	n := rand.Intn(len(highlights))
	highlight := highlights[n]
	return highlight
}

type readwiseHighlight struct {
	Text          string    `json:"text"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Url           string    `json:"url"`
	SourceUrl     string    `json:"source_url"`
	SourceType    string    `json:"source_type"`
	Category      string    `json:"category"`
	LocationType  string    `json:"location_type"`
	Location      int       `json:"location"`
	Note          string    `json:"note"`
	HighlightedAt time.Time `json:"highlighted_at"`
	HightlightUrl string    `json:"highlight_url"`
	ImageUrl      string    `json:"image_url"`
	Id            int       `json:"id"`
	ApiSource     string    `json:"api_source"`
}

type JSONPayload struct {
	ReviewId        int                 `json:"review_id"`
	ReviewUrl       string              `json:"review_url"`
	ReviewCompleted bool                `json:"review_completed"`
	Highlights      []readwiseHighlight `json:"highlights"`
}

var rootCmd = &cobra.Command{
	Use:   "readwise-daily",
	Short: "readwise-daily print daily highlights to the command line",
	Run: func(cmd *cobra.Command, args []string) {
		highlight := getHighlight()
		styleTitle := lipgloss.NewStyle().Padding(2).Align(lipgloss.Center)
		style := lipgloss.NewStyle().Bold(true)
		// pterm.DefaultHeader.Println(highlight.Title)
		// pterm.DefaultBasicText.Println(highlight.Text)
		fmt.Println(styleTitle.Render(highlight.Title))
		fmt.Println(style.Render(highlight.Author))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
