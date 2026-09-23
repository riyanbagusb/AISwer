# Product Requirements Document (PRD)
**Project Name:** AISwer
**Platform:** Desktop (macOS, Windows, Linux)
**Language:** Go (Golang)
**AI Engine:** Google Gemini Pro / Flash

---

## 1. Ringkasan Eksekutif (Executive Summary)
AISwer adalah aplikasi desktop berbasis *menu bar* (system tray) yang berjalan di latar belakang (background) untuk membantu pengguna menjawab soal ujian pilihan ganda secara instan. Saat pengguna menemui soal yang sulit, mereka dapat menekan kombinasi tombol (hotkey) tertentu di keyboard. Aplikasi akan otomatis mengambil tangkapan layar (screenshot), menganalisis soal menggunakan AI (Gemini), dan menampilkan jawaban pilihan ganda (misal: A, B, C, atau D) secara langsung pada indikator *menu bar*.

## 2. Tujuan (Objective)
- Memberikan solusi instan dan akurat untuk soal pilihan ganda tanpa mengharuskan pengguna berpindah aplikasi (seamless experience).
- Memanfaatkan kemampuan Multimodal AI (OCR + Reasoning) secara optimal dalam hitungan detik.
- Menyediakan UI yang 100% terintegrasi dengan OS (melalui Menu Bar/System Tray) sehingga tidak mengganggu layar utama pengguna sama sekali.

## 3. Fitur Utama (Key Features)
### 3.1. Menu Bar Execution & Settings
Aplikasi berjalan tanpa jendela utama (headless) dan langsung bertengger di Menu Bar (macOS) atau System Tray (Windows/Linux). Melalui icon ini, pengguna dapat:
- **Set API Key**: Mengubah Google Gemini API Key kapan saja tanpa perlu restart aplikasi.
- **Pilih Model**: Memilih model AI yang ingin digunakan (daftar model diambil secara dinamis dari API Gemini).
- **Melihat Hasil**: Hasil jawaban AI akan ditampilkan langsung menggantikan icon/teks menu bar.

### 3.2. Global Hotkey Listener
Aplikasi secara terus-menerus mendengarkan input keyboard secara global (meskipun aplikasi lain sedang fokus). Ketika kombinasi tombol spesifik ditekan (default: `Cmd + Shift + X`), proses capture akan dimulai.

### 3.3. Screen Capture (Screenshot)
Ketika hotkey ditekan, aplikasi akan secara diam-diam mengambil tangkapan layar penuh dari monitor.

### 3.4. AI OCR & Reasoning
Gambar hasil tangkapan layar dikirim ke Google Gemini API.
- **Prompt yang digunakan:** *"Baca soal pilihan ganda dari gambar ini beserta pilihan jawabannya. Pilih satu jawaban yang paling tepat dan kembalikan HANYA huruf pilihannya saja (contoh: A, B, C, D, atau E)."*

### 3.5. Smart Permissions Handling
Pada macOS, aplikasi secara otomatis memeriksa apakah izin **Accessibility** dan **Screen Recording** sudah diberikan. Jika belum, aplikasi akan:
- Memunculkan pop-up permintaan izin asli dari sistem macOS (secara otomatis di awal).
- Menampilkan teks peringatan `! Perms` di Menu Bar.
- Melakukan polling secara diam-diam di latar belakang, dan langsung mengaktifkan fitur hotkey begitu izin diberikan oleh pengguna.

### 3.6. API Tier & Rate Limiting
Pengguna dapat memilih antara **Free Tier** dan **Paid Tier**.
- **Free Tier**: Aplikasi menerapkan *rate limit* secara ketat per model (misal: 15 RPM, 250K TPM, 500 RPD untuk model Flash Lite). Model "Pro" tidak dapat digunakan pada tier ini.
- **Paid Tier**: Tidak ada batasan *rate limit* dari sisi aplikasi, dan semua model *Text-out* termasuk varian "Pro" dapat digunakan secara bebas.

## 4. Arsitektur & Teknologi (Tech Stack)
Aplikasi dikembangkan menggunakan bahasa **Go (Golang)** dengan pendekatan *clean code* terdesentralisasi:

1. **Bahasa Pemrograman:** Go 1.21+
2. **System Tray/Menu Bar:** `github.com/getlantern/systray`
3. **Screenshot:** `github.com/kbinani/screenshot`
4. **Global Hotkey:** `golang.design/x/hotkey`
5. **AI Integration:** Google Gen AI SDK untuk Go (`github.com/google/generative-ai-go`)
6. **Environment Config:** `github.com/joho/godotenv`

## 5. User Flow (Alur Pengguna)
1. **Start:** Pengguna menjalankan `AISwer.app`. Aplikasi meminta izin macOS (jika belum ada) lalu diam di Menu Bar dengan tanda `-`.
2. **Trigger:** Pengguna menekan kombinasi `Cmd+Shift+X`.
3. **Processing:**
   - Aplikasi menangkap layar.
   - Menu Bar berubah menjadi `...` (menandakan AI sedang berpikir).
   - Aplikasi mengirim gambar ke Google Gemini API.
4. **Result:** Gemini mengembalikan huruf "C".
5. **Display:** Menu Bar berubah teksnya menjadi **"C"**.
6. **Reset:** Setelah beberapa detik, Menu Bar kembali menjadi tanda `-`.

## 6. Syarat dan Izin Sistem (System Requirements & Permissions)
Khusus untuk pengguna **macOS**:
1. **Accessibility Permission:** Wajib diizinkan agar aplikasi dapat menangkap input keyboard (global hotkey). Aplikasi menggunakan `C.AXIsProcessTrustedWithOptions` untuk meminta izin ini secara otomatis.
2. **Screen Recording Permission:** Wajib diizinkan agar aplikasi dapat mengambil tangkapan layar.

## 7. Fase Pengembangan (Status Saat Ini)
- [x] **Fase 1: Proof of Concept (PoC)** - Hotkey & Capture.
- [x] **Fase 2: AI Integration** - Sambungan ke Gemini Vision.
- [x] **Fase 3: System Tray UI** - Pembuatan menu interaktif, pengaturan API Key, dan daftar model dinamis.
- [x] **Fase 4: Packaging & Refactoring** - Clean code modular, smart permissions, dan packaging menjadi `.app` bundle untuk macOS.
- [x] **Fase 5: API Tier & Rate Limiting** - Pemilihan tier gratis/berbayar, perlindungan limit (RPM, TPM, RPD) per model, dan penyaringan model khusus *Text-out*.
