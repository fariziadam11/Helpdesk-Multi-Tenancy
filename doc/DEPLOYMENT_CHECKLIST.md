# Quick Deployment Checklist

Gunakan checklist ini untuk memastikan semua langkah deployment sudah dilakukan.

## 📋 Pre-Deployment

- [ ] Server sudah siap (Linux/Ubuntu)
- [ ] Akses SSH ke server tersedia
- [ ] Domain/IP server sudah diketahui
- [ ] MySQL credentials sudah disiapkan

## 🔧 Server Setup

- [ ] Go 1.21+ terinstall
- [ ] Bun terinstall
- [ ] MySQL 8 terinstall dan running
- [ ] Nginx terinstall (opsional)
- [ ] Firewall dikonfigurasi (port 22, 80, 443)

## 🗄️ Database Setup

- [ ] Database `armmada` dibuat
- [ ] User MySQL dibuat dengan privileges
- [ ] Test koneksi database berhasil

## 🚀 Backend Deployment

- [ ] File backend di-upload ke server
- [ ] `.env` dikonfigurasi dengan benar
- [ ] Dependencies di-download (`go mod download`)
- [ ] Binary di-build (`./build.sh` atau `go build`)
- [ ] Test run backend berhasil
- [ ] Systemd service dibuat (`werk-ticketing-backend.service`)
- [ ] Service enabled dan running
- [ ] Backend accessible di `http://localhost:8080`
- [ ] Check logs tidak ada error

## 🎨 Frontend Deployment

- [ ] File frontend di-upload ke server
- [ ] `.env` dikonfigurasi (VITE_API_BASE_URL)
- [ ] Dependencies terinstall (`bun install`)
- [ ] Production build dibuat (`bun run build`)
- [ ] Build files di-copy ke `/var/www/werk-ticketing` (jika pakai Nginx)
- [ ] Nginx dikonfigurasi dengan benar
- [ ] Nginx config di-test (`nginx -t`)
- [ ] Nginx di-reload
- [ ] Frontend accessible di browser

## 🔒 Security

- [ ] JWT_SECRET diganti dengan random string yang aman
- [ ] Database password aman
- [ ] `.env` file permissions benar (600)
- [ ] SSL certificate terinstall (jika ada domain)
- [ ] Firewall rules diterapkan
- [ ] Security headers dikonfigurasi di Nginx

## ✅ Testing

- [ ] Buka frontend di browser
- [ ] Test register user baru
- [ ] Test login
- [ ] Test create ticket
- [ ] Test list tickets
- [ ] Test ticket detail dengan comments
- [ ] Test file upload/download
- [ ] Test solution accept/reject
- [ ] Check browser console tidak ada error
- [ ] Check backend logs tidak ada error

## 📊 Monitoring

- [ ] Setup log rotation
- [ ] Setup automated backup untuk database
- [ ] Setup monitoring (opsional)
- [ ] Dokumentasi credentials disimpan dengan aman

## 🎉 Go Live!

- [ ] Inform users tentang URL aplikasi
- [ ] Monitor logs untuk 24 jam pertama
- [ ] Siapkan rollback plan jika ada masalah

---

**Notes:**

- Simpan semua credentials dengan aman
- Backup `.env` file
- Dokumentasikan semua perubahan konfigurasi
- Setup automated backup untuk database
