package main

import (
	"context"
	"log"
	"os"

	"github.com/seymahandekli/git-gen/pkg/gitgen"
	"github.com/urfave/cli/v3"
)

func main() {
	var openAiKey string
	var promptModel string
	var maxTokens int64
	var sourceRef string
	var destinationRef string

	cmd := &cli.Command{
		Name:  "git-gen",
		Usage: "Generate commit messages and perform code reviews using ChatGPT",

		Commands: []*cli.Command{
			{
				Name:  "commit",
				Usage: "Generates a commit message",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "apikey",
						Usage:       "OpenAI API key",
						Sources:     cli.EnvVars("OPENAI_API_KEY"),
						Destination: &openAiKey,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "source",
						Usage:       "Source Ref",
						Value:       "HEAD~",
						Destination: &sourceRef,
					},
					&cli.StringFlag{
						Name:        "dest",
						Usage:       "Destination Ref",
						Value:       "",
						Destination: &destinationRef,
					},
					&cli.IntFlag{
						Name:        "maxtokens",
						Usage:       "Maximum tokens to generate",
						Value:       3500,
						Destination: &maxTokens,
					},
					&cli.StringFlag{
						Name:        "model",
						Usage:       "OpenAI Model",
						Value:       "gpt-4o",
						Destination: &promptModel,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					config := gitgen.DefaultConfig()
					config.OpenApiKey = openAiKey
					config.PromptMaxTokens = maxTokens
					config.PromptModel = promptModel
					config.SourceRef = sourceRef
					config.DestinationRef = destinationRef

					result, err := gitgen.Do(gitgen.PromptCommitMessage, config)

					if err != nil {
						return err
					}

					log.Println(result)
					return nil
				},
			},
			{
				Name:  "review",
				Usage: "Performs a code review",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "apikey",
						Usage:       "OpenAI API key",
						Sources:     cli.EnvVars("OPENAI_API_KEY"),
						Destination: &openAiKey,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "source",
						Usage:       "Source Ref",
						Value:       "HEAD~",
						Destination: &sourceRef,
					},
					&cli.StringFlag{
						Name:        "dest",
						Usage:       "Destination Ref",
						Value:       "",
						Destination: &destinationRef,
					},
					&cli.IntFlag{
						Name:        "maxtokens",
						Usage:       "Maximum tokens to generate",
						Value:       3500,
						Destination: &maxTokens,
					},
					&cli.StringFlag{
						Name:        "model",
						Usage:       "OpenAI Model",
						Value:       "gpt-4o",
						Destination: &promptModel,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					config := gitgen.DefaultConfig()
					config.OpenApiKey = openAiKey
					config.PromptMaxTokens = maxTokens
					config.PromptModel = promptModel
					config.SourceRef = sourceRef
					config.DestinationRef = destinationRef

					result, err := gitgen.Do(gitgen.PromptCodeReview, config)

					if err != nil {
						return err
					}

					log.Println(result)
					return nil
				},
			},
			{
				Name:  "register",
				Usage: "Registers itself to the running system",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					// if strings.Contains(existingPaths, fmt.Sprintf(":%s", workingDirectory)) {
					// 	fmt.Printf("%s path is already registered.\n", workingDirectory)

					// 	return nil
					// }

					// fmt.Printf("setting %s to PATH\n", workingDirectory)
					// newPaths := fmt.Sprintf("%s:%s", existingPaths, workingDirectory)

					// os.Setenv("PATH", newPaths)
					// fmt.Println("done.")

					err := gitgen.RegisterToPath()

					return err
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
		return
	}
}
