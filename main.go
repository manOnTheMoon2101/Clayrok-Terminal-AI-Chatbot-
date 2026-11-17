package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/joho/godotenv"
)

const apiURL = "https://api.x.ai/v1/chat/completions"

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
	Error any `json:"error"`
}

func ask(client *http.Client, key, model string, msgs []message) (string, error) {
	body, err := json.Marshal(chatRequest{Model: model, Messages: msgs})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned %s: %s", resp.Status, strings.TrimSpace(string(data)))
	}
	var out chatResponse
	if err := json.Unmarshal(data, &out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", errors.New("API returned no choices")
	}
	return out.Choices[0].Message.Content, nil
}

func loadEnv() {
	_ = godotenv.Load()
	if dir, err := os.UserConfigDir(); err == nil {
		_ = godotenv.Load(filepath.Join(dir, "clayrok", ".env"))
	}
}

func main() {
	loadEnv()

	key := os.Getenv("XAI_API_KEY")
	if key == "" || key == "your-xai-api-key-here" {
		fmt.Fprintln(os.Stderr, "Set XAI_API_KEY in .env or ~/.config/clayrok/.env")
		os.Exit(1)
	}
	model := os.Getenv("GROK_MODEL")
	if model == "" {
		model = "grok-3"
	}

	client := &http.Client{Timeout: 120 * time.Second}
	history := []message{{Role: "system", Content: "You are a helpful assistant."}}

	fmt.Printf("Clayrok is chatting with %s. Type 'exit' or press Ctrl+C to quit.\n\n", model)

	for {
		var input string
		err := huh.NewInput().
			Title("You").
			Value(&input).
			Run()
		if err != nil {
			break
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		if input == "exit" || input == "quit" {
			break
		}
		fmt.Println("You:", input)

		history = append(history, message{Role: "user", Content: input})

		var reply string
		var reqErr error
		action := func() { reply, reqErr = ask(client, key, model, history) }
		if err := spinner.New().Title("Thinking...").Action(action).Run(); err != nil {
			break
		}
		if reqErr != nil {
			fmt.Fprintln(os.Stderr, "Error:", reqErr)
			history = history[:len(history)-1]
			continue
		}

		history = append(history, message{Role: "assistant", Content: reply})
		fmt.Printf("\nClayrok: %s\n\n", reply)
	}
}
