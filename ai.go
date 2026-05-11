package main

import (
	"context"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func GenerateAIArticle(source string) (string, error) {

	config := openai.DefaultConfig(os.Getenv("GROQ_API_KEY"))

	config.BaseURL = "https://api.groq.com/openai/v1"

	client := openai.NewClientWithConfig(config)

	prompt := `Pahami artikel di bawah ini:

` + source + `

---

Tulis ulang menjadi artikel blog teknologi semi-formal, natural, SEO-friendly, minimal 800 kata, bahasa Indonesia.

WAJIB:
- Judul menarik
- Kalimat pertama kuat
- Natural seperti manusia
- Jangan seperti AI
- Gunakan H2/H3
- Tambahkan insight ringan
- Variasi kalimat
- Jangan kaku
- Jangan copy mentah
- Pertahankan informasi penting
- Jangan sebut sumber

OUTPUT:
Judul
Isi Artikel`

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "llama3-70b-8192",

			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: prompt,
				},
			},

			Temperature: 0.9,
		},
	)

	if err != nil {
		return "", err
	}

	result := resp.Choices[0].Message.Content

	return result, nil
}
