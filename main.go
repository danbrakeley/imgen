package main

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	now := time.Now()

	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("error loading .env file: %v\n", err)
		os.Exit(1)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if len(apiKey) == 0 {
		fmt.Println("error reading OPENAI_API_KEY: empty or missing")
		os.Exit(2)
	}

	configFileName := "imgen.toml"
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}
	cfg, err := LoadFromToml(configFileName)
	if err != nil {
		fmt.Printf("error loading %s: %v\n", configFileName, err)
		os.Exit(3)
	}

	client := openai.NewClient(apiKey)
	ctx := context.Background()

	fmt.Println("Making request...")

	resp, err := client.CreateImage(ctx, openai.ImageRequest{
		Model:          cfg.Model,
		Prompt:         cfg.Prompt,
		N:              1,
		Size:           cfg.Size,
		Quality:        cfg.Quality,
		Style:          cfg.Style,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
	})
	if err != nil {
		fmt.Printf("error calling openai: %v\n", err)
		os.Exit(10)
	}

	if len(resp.Data) != 1 {
		fmt.Printf("expected 1 image, got %d\n", len(resp.Data))
		os.Exit(11)
	}

	dec, err := base64.StdEncoding.DecodeString(resp.Data[0].B64JSON)
	if err != nil {
		fmt.Printf("error decoding base64: %v\n", err)
		fmt.Printf("\nRaw Response:\n")
		fmt.Println("-- Revised Prompt: ", resp.Data[0].RevisedPrompt)
		fmt.Println("-- B64JSON: ", resp.Data[0].B64JSON)
		os.Exit(20)
	}

	h := md5.Sum(dec)

	filename := fmt.Sprintf("%s [%s].png",
		now.Format("2006-01-02-1504"),
		hex.EncodeToString(h[:]),
	)

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}

	fmt.Println("Wrote image to file: ", filename)
	fmt.Println("Revised Prompt: ", resp.Data[0].RevisedPrompt)
}
