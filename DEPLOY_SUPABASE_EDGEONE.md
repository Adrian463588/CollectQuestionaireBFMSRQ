# Panduan Deployment: Supabase (Database) & Tencent EdgeOne Pages (Frontend)

Mendeploy aplikasi ini memerlukan dua pijakan utama: **Supabase** untuk melayani _Database PostgreSQL_ aplikasi Anda secara awan (cloud), dan **Tencent EdgeOne Pages** untuk menayangkan antarmuka web _Next.js_. 

Karena backend (API) khusus Anda dikembangkan dengan bahasa **Golang (Fiber)**, Supabase akan berperan sebagai penyedia Database. Server Golang Anda tetap harus di-hosting di penyedia layanan aplikasi seperti **Heroku, Railway, atau Render** agar dapat mengeksekusi logika skoring, yang mana backend tersebut akan "menghubungkan diri" ke _Database Supabase_.

Berikut adalah langkah-langkah implementasinya:

---

## TAHAP 1: Konfigurasi Database di SUPABASE

Supabase memberikan _managed PostgreSQL_ gratis yang sangat ideal untuk aplikasi skoring kuesioner Anda.

1. **Buat Akun dan Proyek Baru**
   - Kunjungi [Supabase.com](https://supabase.com) dan buat akun.
   - Klik **New Project**, pilih organisasi Anda, lalu beri nama proyek (contoh: `Kuesioner BFM SRQ`).
   - Buat sebuah **Database Password** yang sangat kuat (catat *password* ini baik-baik).
   - Pilih *Region* terdekat (contoh: Singapore) dan klik **Create New Project**. Tunggu beberapa menit hingga database siap.

2. **Dapatkan *Connection String* Database**
   - Masuk ke menu **Project Settings** (ikon roda gigi di kiri bawah) ➔ **Database**.
   - Gulir ke bagian **Connection String** ➔ pilih tab **URI**.
   - Salin *URI* tersebut. Formatnya akan terlihat seperti ini:
     `postgresql://postgres.[namaproyek]:[YOUR-PASSWORD]@aws-0-[region].pooler.supabase.com:6543/postgres`
   - Ganti `[YOUR-PASSWORD]` dengan *password* yang Anda buat di Langkah 1.

3. **Migrasi / Buat Tabel Database (Otomatis)**
   - Karena Anda telah merancang arsitektur Go yang mengeksekusi Auto-Migrasi (melalui GORM `db.Migrate()`), tahap pembuatan tabel ke Supabase akan otomatis dilakukan oleh Backend Golang Anda ketika ia dijalankan. Tidak perlu repot mengeksekusi SQL manual.

---

## TAHAP 2: Peluncuran (Deploy) REST API Backend 

*Catatan: Supabase adalah database, sehingga Server API Golang Anda tetap butuh *hosting* aplikasi seperti Heroku (sebagaimana kita setup sebelumnya pada `Procfile`), Render.com, atau Railway.app.*

**Jika menggunakan Heroku / Railway:**
1. Masukkan *Repository GitHub* Anda. 
2. Pada tab **Environment Variables** (atau *Settings > Config Vars* di Heroku), tambahkan:
   - `SERVER_PORT` = `8080` (untuk heroku cukup sediakan port jika dinamis, server Go kita akan menangkapnya otomatis).
   - `DB_HOST` = (Host pooler dari URI Supabase Anda, misal: `aws-0-ap-southeast-1.pooler.supabase.com`)
   - `DB_PORT` = `6543`
   - `DB_USER` = `postgres.[namaproyek]`
   - `DB_PASSWORD` = `[Password_Supabase_Anda]`
   - `DB_NAME` = `postgres`
   - `DATABASE_URL` = (Opsi alternatif: *Paste URI lengkap* koneksi Supabase Anda).
3. Lakukan proses Deploy. Aplikasi Go akan otomatis menembak server Supabase Anda, membangun tabel SQL dari nol, dan API Anda pun resmi terhubung! (Salin *URL publik backend* Anda setelah ini selesai).

---

## TAHAP 3: Peluncuran (Deploy) Frontend di TENCENT EDGE ONE PAGES

Tencent EdgeOne (TEO) Pages menyediakan sarana hosting antarmuka statis (*static sites*) dan kerangka kerja *Next.js* yang cepat di berbagai wilayah secara global.

1. **Masuk ke Konsol Tencent Cloud EdgeOne**
   - Lakukan registrasi/login ke akun [Tencent Cloud EdgeOne Pages](https://console.tencentcloud.com/teo/pages).
   
2. **Buat Proyek Baru Terhubung GitHub**
   - Klik **Create Project** / **New Pages**.
   - Pada opsi sumber kode, pilih **Connect to GitHub**.
   - Otorisasi GitHub Anda, lalu pilih repository: `Adrian463588/CollectQuestionaireBFMSRQ`.
   - Pilih *branch* `main`.

3. **Penyetelan Konfigurasi Build Next.js**
   Agar EdgeOne meluncurkan aplikasi dengan benar, atur parameternya sebagai berikut:
   - **Root Directory**: `frontend` *(Penting! Karena kode UI kita ada di folder frontend).*
   - **Framework Preset**: Pilih **Next.js** (atau biarkan sistem mendeteksi otomatis).
   - **Build Command**: `npm run build`
   - **Output Directory**: `.next`

4. **Sisipkan Environment Variables (Keamanan)**
   Sebelum menekan tombol Deploy, tambahkan variabel krusial pengait frontend ke backend:
   - Cari menu **Environment Variables**.
   - Tambahkan variabel pertama: 
     - **Key**: `NEXT_PUBLIC_ADMIN_PASSWORD`
     - **Value**: `Admin123` *(Atau kata sandi acak yang aman sesuai mau Anda)*
   - Tambahkan variabel kedua:
     - **Key**: `NEXT_PUBLIC_API_URL`
     - **Value**: `https://<url-backend-anda-dari-tahap-2>/api` *(Pastikan menyertakan prefix /api dan jangan gunakan slash di akhir)*

5. **Deploy & Selesai**
   - Tekan **Deploy / Save**. EdgeOne akan mengunduh dependensi (NPM) dan melakukan proses pembangunan/build (`npm run build`).
   - Tunggu proses ini hingga menampikan status ✅ **Success/Deployed**.
   - EdgeOne Pages akan memberikan URL gratis secara instan (Biasanya `xxx.pages.dev` atau serupa dari layanan EdgeOne).
   
### 🚀 Cek Hasil Akhirnya!
1. Buka domain publik yang diberikan oleh EdgeOne.
2. Isi sebuah tes kuesioner baru di halaman pendaftaran.
3. Klik tombol "Lihat Grafik Interpretasi". 
4. Layar seharusnya secara ketat meminta input Password (seperti yang telah Anda set pada `NEXT_PUBLIC_ADMIN_PASSWORD` sebelumnya), dan seluruh log respon kuesioner aman tersimpan pada Supabase!
