# gomake

Sebuah aplikasi cli untuk memulai proyek Go barumu. Menggunakan git dan makefile.

## Instalasi

### go install

gunakan go untuk memasang `gomake` :

```bash
go install github.com/LitFill/gomake@latest
```

### dependensi

`gomake` menggunakan git dan makefile, jadi perlu `git` dan `make` untuk terpasang terlebih dahulu.

## Penggunaan

```bash
gomake "LitFill/program"
```

perintah ini akan membuat direktori `program`, menjalankan `go mod init LitFill/program`, membuat `main.go`, menginisialisasi repo git, dan membuat `Makefile`.

kemudian kamu bisa menggunakan Makefile seperti ini:

```bash
make        # membangun program untuk linux
make win    # membangun program untuk windows
make run    # menjalankan program
make help   # menampilkan pesan bantuan
```

README.md in English : [README.md](./README.md)
