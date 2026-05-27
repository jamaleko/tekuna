package main

import (
 "context"
 "regexp"
 "strings"
 
 //"os"

 openai "github.com/sashabaranov/go-openai"
 "golang.org/x/text/unicode/norm"
)

func GenerateAIArticle(source string) (string, error) {

 config :=
  openai.DefaultConfig("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")

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

- Output WAJIB dimulai dari karakter pertama dengan:

JUDUL:

- Kata pertama di baris pertama WAJIB "JUDUL:"
- Setelah judul selesai, WAJIB ada:

ISI:

- Kata "JUDUL:" dan "ISI:" tidak boleh diubah, dihilangkan, diterjemahkan, atau diganti
- Jika output tidak diawali "JUDUL:" maka output dianggap TIDAK VALID
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
- Hindari pengulangan pola kalimat
- Hindari frasa transisi yang sama berulang kali
- Variasikan pembuka paragraf
- Tulis seperti artikel yang benar-benar ditulis manusia,
  bukan hasil parafrase AI

HINDARI:
- "Artikel ini akan membahas"
- "Di era digital saat ini"

REWRITE KUTIPAN:
- Isi kutipan di dalam tanda petik wajib dipertahankan utuh
- Attribution/pengantar kutipan wajib dibuat bervariasi dan natural
- Jangan mengulang frasa attribution yang sama terus menerus
- Hindari pengulangan seperti:
  - katanya
  - ujar dia
  - dan bahwa
  - menurutnya
  di banyak paragraf
- Gunakan variasi attribution yang lebih manusiawi dan kontekstual

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
  result = norm.NFC.String(result)
  // hapus karakter rusak
  result = strings.ReplaceAll(result, "??", "")
  result = strings.ReplaceAll(result, "�", "")
  // hapus TAB
result = strings.ReplaceAll(result, "\t", " ")

// fix enter windows
result = strings.ReplaceAll(result, "\r", "")

// FIX LIST DALAM <p>
re := regexp.MustCompile(`(?s)<p>(.*?)((\d+\..*?)+)</p>`)

result = re.ReplaceAllStringFunc(result, func(match string) string {

 clean := strings.ReplaceAll(match, "<p>", "")
 clean = strings.ReplaceAll(clean, "</p>", "")

 lines := strings.Split(clean, "\n")

 var normal []string
 var items []string

 numRe := regexp.MustCompile(`^\d+\.\s*`)

 for _, line := range lines {

  line = strings.TrimSpace(line)

  if numRe.MatchString(line) {

   line = numRe.ReplaceAllString(line, "")
   items = append(items, "<li>"+line+"</li>")

  } else {

   if line != "" {
    normal = append(normal, line)
   }
  }
 }

 if len(items) == 0 {
  return "<p>" + clean + "</p>"
 }

 return "<p>" +
  strings.Join(normal, " ") +
  "</p><ol>" +
  strings.Join(items, "") +
  "</ol>"
})
 return result, nil
}
