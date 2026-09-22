# Product Requirements Document (PRD)
**Project Name:** Exam OCR Assistant (Auto-Answer MCQ)
**Platform:** Desktop (macOS, Windows, Linux)
**Language:** Go (Golang)
**AI Engine:** Google Gemini Pro Vision (Rekomendasi utama) / Google Cloud Vision + LLM

---

## 1. Ringkasan Eksekutif (Executive Summary)
Exam OCR Assistant adalah aplikasi desktop yang berjalan di latar belakang (background) untuk membantu pengguna menjawab soal ujian pilihan ganda secara instan. Saat pengguna menemui soal yang sulit, mereka dapat menekan kombinasi tombol (hotkey) tertentu di keyboard. Aplikasi akan otomatis mengambil tangkapan layar (screenshot), menganalisis soal menggunakan AI (Gemini), dan menampilkan jawaban pilihan ganda (misal: A, B, C, atau D) melalui pop-up kecil yang tidak mencolok di pojok kanan bawah layar.

## 2. Tujuan (Objective)
- Memberikan solusi instan dan akurat untuk soal pilihan ganda tanpa mengharuskan pengguna berpindah aplikasi (seamless experience).
- Memanfaatkan kemampuan Multimodal AI (OCR + Reasoning) secara optimal dalam hitungan detik.
- Menyediakan UI yang minimalis dan tidak mengganggu layar utama pengguna.

## 3. Fitur Utama (Key Features)
### 3.1. Background Execution & System Tray
Aplikasi berjalan tanpa jendela utama (headless) saat dimulai. Aplikasi akan memunculkan icon di System Tray (Windows) atau Menu Bar (macOS) untuk menandakan bahwa aplikasi sedang aktif. Melalui icon ini, pengguna dapat melakukan pengaturan atau mematikan aplikasi.

### 3.2. Global Hotkey Listener
Aplikasi secara terus-menerus mendengarkan input keyboard secara global (meskipun aplikasi lain sedang fokus). Ketika kombinasi tombol spesifik ditekan (misal: `Cmd + Shift + X` di Mac atau `Ctrl + Shift + X` di Windows), proses capture akan dimulai.

### 3.3. Screen Capture (Screenshot)
Ketika hotkey ditekan, aplikasi akan secara diam-diam (tanpa suara/animasi besar) mengambil tangkapan layar dari monitor utama atau area tertentu.

### 3.4. AI OCR & Reasoning
Gambar hasil tangkapan layar akan diubah ke format yang sesuai (misal base64 atau byte array) dan dikirim ke API AI.
- **Rekomendasi Utama:** Menggunakan **Google Gemini (Gemini 1.5 Flash / Pro)** karena memiliki kemampuan vision (mengenali teks dari gambar) sekaligus reasoning (menjawab soal) dalam satu pemanggilan API, sehingga lebih cepat dan murah.
- **Prompt yang digunakan:** *"Baca soal pilihan ganda dari gambar ini beserta pilihan jawabannya. Pilih satu jawaban yang paling tepat dan kembalikan HANYA huruf pilihannya saja (contoh: A, B, C, D, atau E)."*

### 3.5. Minimalist Overlay / Pop-up Notification
Setelah mendapat balasan dari AI, aplikasi akan menampilkan pop-up berukuran kecil (toast/overlay) di pojok kanan bawah layar.
- Pop-up ini bersifat "Always on Top" agar menimpa aplikasi ujian/browser.
- Pop-up berisi teks singkat, misalnya: **"C"**.
- Pop-up akan menghilang secara otomatis (auto-dismiss) setelah 3-5 detik.

## 4. Arsitektur & Teknologi (Tech Stack)
Aplikasi akan dikembangkan menggunakan bahasa **Go (Golang)**. Berikut adalah komponen yang disarankan:

1. **Bahasa Pemrograman:** Go 1.21+
2. **Screenshot:** `github.com/kbinani/screenshot` (Lintas platform untuk menangkap layar).
3. **Global Hotkey:** `golang.design/x/hotkey` (Mendukung macOS, Windows, Linux secara stabil).
4. **GUI / Pop-up:** 
   - Opsi 1: `github.com/ncruces/zenity` untuk dialog sistem (sederhana).
   - Opsi 2: `fyne.io/fyne/v2` (Library GUI Go asli, bisa membuat window kecil, borderless, dan always-on-top).
   - Opsi 3: `github.com/wailsapp/wails/v2` (Jika butuh desain pop-up berbasis HTML/CSS yang sangat cantik dan modern).
5. **AI Integration:** Google Gen AI SDK untuk Go (`github.com/google/genai-go` atau REST API standar).

## 5. User Flow (Alur Pengguna)
1. **Start:** Pengguna menjalankan aplikasi `exam-ocr`. Aplikasi berjalan di latar belakang (muncul di menu bar).
2. **Encounter Question:** Pengguna sedang mengerjakan soal ujian di web browser.
3. **Trigger:** Pengguna menekan kombinasi `Cmd+Shift+X`.
4. **Processing:**
   - Aplikasi menangkap layar (screenshot).
   - Aplikasi mengirim gambar ke Google Gemini API.
   - *Optional:* Muncul indikator loading kecil ("Sedang berpikir...").
5. **Result:** Gemini mengembalikan string "C".
6. **Display:** Jendela kecil muncul di pojok kanan bawah: **"C"**.
7. **End:** Jendela tertutup otomatis setelah 3 detik. Pengguna mengklik jawaban "C" pada soal ujian.

## 6. Syarat dan Izin Sistem (System Requirements & Permissions)
Khusus untuk pengguna **macOS**:
1. **Screen Recording Permission:** Wajib diizinkan di System Settings > Privacy & Security > Screen Recording agar aplikasi Go dapat mengambil tangkapan layar.
2. **Accessibility Permission:** Wajib diizinkan agar aplikasi dapat menangkap input keyboard (global hotkey) saat tidak sedang aktif.

## 7. Fase Pengembangan (Milestones)
- **Fase 1: Proof of Concept (PoC)**
  - Membuat script Go sederhana yang bisa di-trigger dengan hotkey, lalu mengambil screenshot dan menyimpannya di file lokal.
- **Fase 2: AI Integration**
  - Menyambungkan screenshot tersebut secara langsung ke Gemini API dan mem-print hasil jawaban di terminal.
- **Fase 3: Overlay UI**
  - Membangun pop-up window (borderless, always-on-top) di pojok kanan bawah untuk menampilkan hasil.
- **Fase 4: Packaging & Polishing**
  - Membuat aplikasi menjadi executable background app (tanpa console terminal yang terbuka). Menambahkan opsi system tray.

## 8. Out of Scope (Di Luar Cakupan Saat Ini)
- Crop manual: Saat ini aplikasi mengambil layar penuh. Memilih area spesifik (seperti Snipping Tool) mungkin akan menambah interaksi yang tidak diinginkan (memperlambat alur).
- Mode Offline: Mengingat kebutuhan reasoning soal yang rumit, mode offline murni belum didukung.
