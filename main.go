package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	gpt3 "github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetResponse(ctx context.Context, client gpt3.Client, question string) {
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			question,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0.5),
	}, func(res *gpt3.CompletionResponse) {

		fmt.Print(res.Choices[0].Text)
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	fmt.Printf("\n ")
}

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func main() {
	log.SetOutput(new(NullWriter))
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	apiKey := viper.GetString("OPEN_AI_API_KEY")

	if apiKey == "" {
		panic("OPEN_AI_API_KEY not found")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with GPT-3 in your terminal",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)
			quit := false

			for !quit {
				fmt.Println("\n")
				fmt.Println("=> Ask/say something ('quit' to exit):")
				if !scanner.Scan() {
					break
				}

				question := scanner.Text()
				switch question {
				case "quit":
					quit = true

				default:
					fmt.Println("👇")
					GetResponse(ctx, client, question)
				}
			}

		},
	}

	rootCmd.Execute()

}
