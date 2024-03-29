# gomake

Sebuah aplikasi cli untuk memulai proyek Go barumu. Menggunakan git dan makefile.

## Penggunaan

```console
gomake "LitFill/program"
```

perintah ini akan membuat direktori `program`, menjalankan `go mod init LitFill/program`, membuat `main.go`, menginisialisasi repo git, dan membuat `Makefile`.

kemudian kamu bisa menggunakan Makefile seperti ini:

```console
make        # membangun program untuk linux
make win    # membangun program untuk windows
make run    # menjalankan program
make help   # menampilkan pesan bantuan
```

README.md in English : [README.md](./README.md)
