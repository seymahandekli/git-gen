package main

import (
	"context"
	"log"
	"os"

	"github.com/seymahandekli/git-gen/pkg/gitgen"
	"github.com/urfave/cli/v3"
)

func main() {
	var platformApiKey string
	var sourceRef string
	var destinationRef string
	var platform string
	var model string
	var maxTokens int64

	log.SetFlags(0)

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
						Usage:       "Platform API key",
						Sources:     cli.EnvVars("PLATFORM_API_KEY"),
						Destination: &platformApiKey,
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
					&cli.StringFlag{
						Name:        "platform",
						Usage:       "Platform",
						Value:       "openai",
						Destination: &platform,
					},
					&cli.StringFlag{
						Name:        "model",
						Usage:       "Model",
						Value:       "",
						Destination: &model,
					},
					&cli.IntFlag{
						Name:        "maxtokens",
						Usage:       "Maximum tokens to generate",
						Value:       3500,
						Destination: &maxTokens,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					config := gitgen.NewConfig(
						gitgen.WithPlatformApiKey(platformApiKey),
						gitgen.WithSourceRef(sourceRef),
						gitgen.WithDestinationRef(destinationRef),
						gitgen.WithPlatform(platform),
						gitgen.WithModel(model),
						gitgen.WithPromptMaxTokens(maxTokens),
					)

					result, err := gitgen.Do(gitgen.PromptCommitMessage, *config)

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
						Usage:       "Platform API key",
						Sources:     cli.EnvVars("PLATFORM_API_KEY"),
						Destination: &platformApiKey,
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
					&cli.StringFlag{
						Name:        "platform",
						Usage:       "Platform",
						Value:       "openai",
						Destination: &platform,
					},
					&cli.StringFlag{
						Name:        "model",
						Usage:       "Model",
						Value:       "",
						Destination: &model,
					},
					&cli.IntFlag{
						Name:        "maxtokens",
						Usage:       "Maximum tokens to generate",
						Value:       3500,
						Destination: &maxTokens,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					config := gitgen.NewConfig(
						gitgen.WithPlatformApiKey(platformApiKey),
						gitgen.WithSourceRef(sourceRef),
						gitgen.WithDestinationRef(destinationRef),
						gitgen.WithPlatform(platform),
						gitgen.WithModel(model),
						gitgen.WithPromptMaxTokens(maxTokens),
					)

					result, err := gitgen.Do(gitgen.PromptCodeReview, *config)

					if err != nil {
						return err
					}

					log.Println(result)
					return nil
				},
			},
			{
				Name:  "test",
				Usage: "Creating test cases",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "apikey",
						Usage:       "Platform API key",
						Sources:     cli.EnvVars("PLATFORM_API_KEY"),
						Destination: &platformApiKey,
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
					&cli.StringFlag{
						Name:        "platform",
						Usage:       "Platform",
						Value:       "openai",
						Destination: &platform,
					},
					&cli.StringFlag{
						Name:        "model",
						Usage:       "Model",
						Value:       "",
						Destination: &model,
					},
					&cli.IntFlag{
						Name:        "maxtokens",
						Usage:       "Maximum tokens to generate",
						Value:       3500,
						Destination: &maxTokens,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					config := gitgen.NewConfig(
						gitgen.WithPlatformApiKey(platformApiKey),
						gitgen.WithSourceRef(sourceRef),
						gitgen.WithDestinationRef(destinationRef),
						gitgen.WithPlatform(platform),
						gitgen.WithModel(model),
						gitgen.WithPromptMaxTokens(maxTokens),
					)

					result, err := gitgen.Do(gitgen.PromptTestCase, *config)

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
