package main

import (
	"context"
	"log"
	"os"

	"github.com/seymahandekli/git-gen/pkg/gitgen"
	"github.com/seymahandekli/git-gen/pkg/models"
	"github.com/urfave/cli/v3"
)

func main() {
	var openAiKey string
	var sourceRef string
	var destinationRef string
	var promptModel string
	var ollamaAiModel string
	var maxTokens int64

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
					},
					&cli.StringFlag{
						Name:        "source",
						Usage:       "Source Ref",
						Value:       "HEAD",
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
					&cli.StringFlag{
						Name:        "ollamaai",
						Usage:       "OllamaAI Model",
						Value:       "llama3",
						Destination: &ollamaAiModel,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var result string
					var err error
					var runtime models.Model

					if openAiKey != "" {
						config := gitgen.NewConfig(
							gitgen.WithOpenAiKey(openAiKey),
							gitgen.WithSourceRef(sourceRef),
							gitgen.WithDestinationRef(destinationRef),
							gitgen.WithPromptModel(promptModel),
							gitgen.WithPromptMaxTokens(maxTokens),
						)
						runtime = models.NewOpenAi()
						result, err = gitgen.Do(gitgen.PromptCommitMessage, *config, runtime)

						if err != nil {
							return err
						}
						log.Println(result)
						return nil

					}

					config := gitgen.NewConfig(
						gitgen.WithSourceRef(sourceRef),
						gitgen.WithDestinationRef(destinationRef),
						gitgen.WithOllamaAiModel(ollamaAiModel),
					)
					runtime = models.NewOllamaAi()
					result, err = gitgen.Do(gitgen.PromptCommitMessage, *config, runtime)
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
					},
					&cli.StringFlag{
						Name:        "source",
						Usage:       "Source Ref",
						Value:       "HEAD",
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
					&cli.StringFlag{
						Name:        "ollamaai",
						Usage:       "OllamaAI Model",
						Value:       "llama3",
						Destination: &ollamaAiModel,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var result string
					var err error
					var runtime models.Model

					if openAiKey != "" {
						config := gitgen.NewConfig(
							gitgen.WithOpenAiKey(openAiKey),
							gitgen.WithSourceRef(sourceRef),
							gitgen.WithDestinationRef(destinationRef),
							gitgen.WithPromptModel(promptModel),
							gitgen.WithPromptMaxTokens(maxTokens),
						)
						runtime = models.NewOpenAi()
						result, err = gitgen.Do(gitgen.PromptCodeReview, *config, runtime)

						if err != nil {
							return err
						}
						log.Println(result)
						return nil

					}
					config := gitgen.NewConfig(
						gitgen.WithSourceRef(sourceRef),
						gitgen.WithDestinationRef(destinationRef),
						gitgen.WithOllamaAiModel(ollamaAiModel),
					)
					runtime = models.NewOllamaAi()
					result, err = gitgen.Do(gitgen.PromptCodeReview, *config, runtime)
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
