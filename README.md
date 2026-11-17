# Clayrok

A simple terminal chatbot powered by the Grok (xAI) API. It keeps the conversation history, so Grok remembers earlier messages in a session.

## Stack

- [Go](https://go.dev/)
- [huh](https://github.com/charmbracelet/huh) for the input prompt and loading spinner
- [godotenv](https://github.com/joho/godotenv) for loading `.env` files
- [xAI API](https://docs.x.ai/) (OpenAI-compatible `chat/completions` endpoint)

## Setup

1. Get an API key from [console.x.ai](https://console.x.ai/).
2. Create a `.env` file (see `.env.example`):

   ```env
   XAI_API_KEY=your-xai-api-key
   GROK_MODEL=grok-3
   ```

   `GROK_MODEL` is optional and defaults to `grok-3`.

Clayrok looks for `.env` in two places, in this order:

1. The current directory
2. `~/.config/clayrok/.env`

Variables already set in your shell are never overwritten. For running `clayrok` from anywhere, put your key in `~/.config/clayrok/.env`.

## Run

From the project folder:

```sh
go run .
```

## Install

From the project folder:

```sh
go install .
```

This builds the `clayrok` binary into your Go bin directory (`$GOBIN`, or `$GOPATH/bin`). Make sure that directory is on your `PATH`, then run it from anywhere:

```sh
clayrok
```

After changing the code, run `go install .` again to update the command.

## Usage

Type a message and press Enter. Type `exit` or `quit`, or press Ctrl+C, to leave.

