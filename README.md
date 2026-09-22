# AISwer Assistant

AISwer Assistant adalah aplikasi menu bar (system tray) macOS yang memungkinkan Anda untuk dengan cepat menangkap tangkapan layar (screenshot) dari soal ujian (pilihan ganda) dan langsung menggunakan AI (Google Gemini) untuk mencari jawaban yang paling tepat.

Aplikasi ini berjalan di latar belakang tanpa mengganggu layar Anda.

## Prasyarat

- **Sistem Operasi**: macOS (direkomendasikan)
- **Go**: Versi 1.27 atau lebih baru
- **API Key**: Google Gemini API Key dari Google AI Studio

## Cara Build & Install

Aplikasi ini dilengkapi dengan script build otomatis yang akan mem-package aplikasi ke dalam format macOS App Bundle (`.app`). Format ini memungkinkan aplikasi berjalan rapi di latar belakang tanpa memunculkan jendela terminal saat diklik.

1. Clone repositori ini:
   ```bash
   git clone <url-repositori-anda>
   cd "AISwer"
   ```

2. Jalankan script build macOS:
   ```bash
   chmod +x build_mac.sh
   ./build_mac.sh
   ```

3. Jika build berhasil, Anda akan menemukan file **`AISwer.app`** di dalam folder `dist/`.

## Cara Menjalankan Aplikasi

1. Buka folder `dist/` dan *double-click* aplikasi **`AISwer.app`**.
2. **Perizinan (Permissions):**
   - Saat pertama kali dibuka, aplikasi akan memunculkan pop-up yang meminta izin **Accessibility**. Izin ini diperlukan agar aplikasi dapat mendengarkan kombinasi tombol keyboard (hotkey) secara global.
   - Buka `System Settings -> Privacy & Security -> Accessibility`, lalu centang/aktifkan untuk `AISwer.app`.
   - Setelah Accessibility diizinkan, aplikasi secara otomatis akan melakukan satu screenshot *dummy* di latar belakang untuk memicu permintaan izin **Screen Recording**. Ini wajib diizinkan agar aplikasi bisa membaca soal dari layar Anda.
3. Setelah semua izin diberikan, Anda akan melihat ikon aplikasi (atau tanda `-`) muncul di menu bar macOS Anda di sudut kanan atas.

## Pengaturan (Configuration)

Anda dapat melakukan konfigurasi langsung melalui ikon aplikasi di menu bar:
- **Set API Key...**: Klik menu ini untuk memasukkan atau mengubah Google Gemini API Key Anda.
- **Model**: Pilih model Gemini yang ingin Anda gunakan. Daftar model ini diambil langsung secara dinamis dari server Gemini API.
- **Quit**: Keluar dari aplikasi.

*Catatan: Segala perubahan API Key atau pilihan model akan langsung disimpan ke dalam file `.env` dan diaplikasikan tanpa perlu merestart aplikasi.*

## Cara Menggunakan

1. Temukan soal pilihan ganda di layar komputer Anda.
2. Tekan kombinasi tombol **`Cmd + Shift + X`** (kombinasi *hotkey* default).
3. Aplikasi akan mengambil screenshot dari layar Anda dan mengirimkannya ke Gemini AI.
4. Ikon menu bar akan berubah menjadi `...` (menandakan AI sedang memproses).
5. Setelah selesai, teks ikon pada menu bar akan berubah menjadi huruf jawaban (misalnya: **`A`**, **`B`**, **`C`**).

## Troubleshooting

- **Muncul `! Perms` di menu bar**: Anda belum memberikan izin Accessibility. Buka System Settings macOS Anda untuk mengizinkannya.
- **Muncul `! Key` di menu bar**: API Key Anda belum diatur atau tidak valid. Silakan klik "Set API Key..." dan masukkan key yang valid.
- **Muncul `! AI` di menu bar**: Terjadi kesalahan saat menghubungi server Gemini (misalnya masalah jaringan atau token API tidak valid).
- **Muncul `! Capture` di menu bar**: Aplikasi gagal mengambil screenshot. Pastikan izin "Screen Recording" telah diberikan di System Settings.
