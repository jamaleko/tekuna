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

 TUJUAN SEO (DINAMIS & WAJIB):

- Tentukan sendiri 1 keyword utama paling relevan
- Tentukan 2–4 keyword turunan (semantic & long-tail)

---

 ATURAN KEYWORD:

- Keyword utama WAJIB muncul di:
  
  - judul
  - kalimat pertama (karena akan jadi meta deskripsi otomatis)

- Keyword turunan:
  
  - disebar natural di subjudul & isi
  - jangan dipaksakan

---

 PENTING (META DESKRIPSI OTOMATIS):

- 120 karakter pertama = meta deskripsi
- Jadi kalimat pertama HARUS:
  - menarik (hook)
  - jelas
  - mengandung keyword utama
  - tidak bertele-tele

Contoh gaya:

- langsung ke inti
- atau pertanyaan yang menggugah

---

 PRINSIP UTAMA:

- Jangan hanya rewrite
- Tulis seolah benar-benar paham topik
- Tambahkan insight & opini ringan
- Utamakan flow natural

---

 GAYA PENULISAN:

- Semi-formal (blog teknologi)
- Natural & mengalir
- Tidak terlalu rapi seperti AI
- Variasi kalimat (pendek & panjang)
- Bahasa hidup, tidak kaku

---

 HINDARI:

- Gaya berita formal
- Kalimat template:
  “Artikel ini akan membahas…”
  “Di era digital saat ini…”
- Pola monoton

---

 PENDEKATAN HUMAN + SEO:

- Jelaskan dengan sudut pandang sendiri

- Tambahkan interpretasi

- Gunakan analogi jika perlu

- Bahas:
  
  - realistis atau tidak
  - siapa yang diuntungkan
  - dampak ke pengguna

- Gunakan keyword semantic alami

---

 STRUKTUR OUTPUT:

1. Judul artikel

- Harus mengandung keyword utama
- Natural & menarik

2. Isi artikel

- Paragraf pertama:
  
  - HARUS kuat (karena jadi meta deskripsi)
  - mengandung keyword utama
  - maksimal 2–3 kalimat awal sudah “kena”

- Gunakan H2/H3

- Subjudul boleh mengandung keyword turunan

- Paragraf pendek (2–4 baris)

- Gunakan transisi natural

- Jangan terlalu sempurna (biar human)

- Penutup:
  
  - ringkasan + opini/prediksi ringan
  - boleh mengulang keyword secara halus

---

 TEKNIS:

- Minimal 800 kata
- Bahasa Indonesia
- 100% original
- Jangan mengikuti struktur sumber
- Jangan menyebutkan sumber

---

REWRITE KUTIPAN & UCAPAN (WAJIB):

- Jika ada kata seperti:
  
  - ujar
  - tutur
  - kata dia
  - menurutnya
  - jelasnya
  - ungkapnya
  - katanya
  
  maka WAJIB di-rewrite agar lebih natural dan bervariasi.

- Semua kutipan penting narasumber dari artikel sumber WAJIB dipertahankan secara utuh.
- Jangan memotong isi kutipan utama.
- Jangan menghilangkan detail penting dalam ucapan narasumber.
- Tanda petik ("...") wajib dipertahankan.
- Yang boleh diubah hanya gaya attribution/pengantar kutipan agar lebih natural.

---

 OUTPUT:
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
