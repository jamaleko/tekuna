package main

import (
 "context"
 "os"

 openai "github.com/sashabaranov/go-openai"
)

func GenerateAIArticle(source string) (string, error) {

 config :=
  openai.DefaultConfig(
   os.Getenv("GROQ_API_KEY"),
  )

 config.BaseURL =
  "https://api.groq.com/openai/v1"

 client :=
  openai.NewClientWithConfig(config)

 systemPrompt := `
Kamu adalah penulis blog teknologi Indonesia.

Tugas:
Rewrite artikel menjadi artikel baru yang natural, human-like, dan SEO friendly.

OUTPUT WAJIB:

JUDUL:
[judul baru]

ISI:
[isi artikel HTML]

ATURAN KRITIS:
- Semua jenis tanda kutip wajib dipertahankan:
  - "..."
  - “...”
- Isi kutipan narasumber tidak boleh diubah
- Jangan menyebut sumber asli
- Jangan gunakan markdown
- Gunakan HTML <h2> dan <p>
- Minimal 800 kata
- Bahasa Indonesia
- Output hanya boleh berisi format:
  JUDUL:
  ISI:

ATURAN SEO:
- Tentukan keyword utama sendiri
- Keyword utama wajib muncul di:
  - judul
  - kalimat pertama
- Kalimat pertama maksimal 120 karakter dan menarik

GAYA PENULISAN:
- Semi formal
- Natural
- Tidak kaku
- Variasi panjang kalimat
- Tambahkan insight ringan

HINDARI:
- "Artikel ini akan membahas"
- "Di era digital saat ini"

REWRITE KUTIPAN:
- Attribution seperti:
  - ujar
  - tutur
  - jelasnya
  - katanya
  boleh diubah agar lebih natural
- Tetapi isi kutipan di dalam tanda petik wajib dipertahankan utuh

Gunakan struktur HTML seperti:
<h2>Subjudul</h2>
<p>Isi paragraf</p>
`

 userPrompt :=
  "Pahami artikel berikut lalu rewrite:\n\n" +
   source

 resp, err := client.CreateChatCompletion(
  context.Background(),
  openai.ChatCompletionRequest{

   Model: "llama-3.3-70b-versatile",

   Messages: []openai.ChatCompletionMessage{

    {
     Role: "system",

     Content: systemPrompt,
    },

    {
     Role: "user",

     Content: userPrompt,
    },
   },

   Temperature: 0.3,
  },
 )

 if err != nil {
  return "", err
 }

 result :=
  resp.Choices[0].Message.Content

 return result, nil
}
