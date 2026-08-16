package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/lib/pq"
)

type Review struct {
	Nama     string
	Rating   string
	Komentar string
}

type Kontrakan struct {
	ID        string
	Nama      string
	Kota      string
	Harga     string
	Foto      string
	Galeri    []string
	Deskripsi string
	Lokasi    string
	Fasilitas string
	Kebijakan string
	Reviews   []Review
}

// Global variable untuk koneksi database PostgreSQL
var db *sql.DB

var databaseKontrakan = map[string]Kontrakan{
	"jakarta": {
		ID:        "jakarta",
		Nama:      "Kontrakan Ancol",
		Kota:      "Jakarta",
		Harga:     "Rp 3.200.000",
		Foto:      "(G1)Dafkon.png",
		Galeri:    []string{"(G1)Dafkon.png", "(G0)Dafkon.png", "(G2)Dafkon.png", "(G3)Dafkon.png", "(G4)Dafkon.png"},
		Deskripsi: "Kontrakan bersih, aman, dan nyaman di pusat kota Jakarta.",
		Lokasi:    "Jl. Lodan Raya No. 12, Ancol, Jakarta Utara",
		Fasilitas: "AC, WiFi 100Mbps, Kasur Springbed, Kamar Mandi Dalam.",
		Kebijakan: "Jam malam 23.00 WIB, Dilarang membawa hewan peliharaan.",
		Reviews: []Review{
			{Nama: "Budi Santoso", Rating: "⭐⭐⭐⭐⭐", Komentar: "Dekat stasiun dan Dufan, mantap!"},
			{Nama: "Siti Rahma", Rating: "⭐⭐⭐⭐☆", Komentar: "Akses jalan lumayan oke, lingkungan tenang."},
		},
	},
	"bogor": {
		ID:        "bogor",
		Nama:      "Kontrakan Mbah Rusdi",
		Kota:      "Bogor",
		Harga:     "Rp 1.500.000",
		Foto:      "(G2)Dafkon.png",
		Galeri:    []string{"(G2)Dafkon.png", "(G3)Dafkon.png", "(G4)Dafkon.png", "(G5)Dafkon.png", "(GED)Dafkon.png"},
		Deskripsi: "Udara sejuk khas Puncak Cisarua, lingkungan sangat tenang dan asri.",
		Lokasi:    "Jl. Cisarua Raya No. 45, Bogor, Jawa Barat",
		Fasilitas: "WiFi, Saung Santai, Dapur Bersama, Area Parkir Luas.",
		Kebijakan: "Dilarang membuat keributan di atas jam 22.00 WIB.",
		Reviews: []Review{
			{Nama: "Fauzi R.", Rating: "⭐⭐⭐⭐⭐", Komentar: "Udaranya adem banget, cocok buat istirahat!"},
		},
	},
	"depok": {
		ID:        "depok",
		Nama:      "Kontrakan Biru",
		Kota:      "Depok",
		Harga:     "Rp 2.000.000",
		Foto:      "(G3)Dafkon.png",
		Galeri:    []string{"(G3)Dafkon.png", "(G0)Dafkon.png", "(G1)Dafkon.png", "(G4)Dafkon.png", "(G5)Dafkon.png"},
		Deskripsi: "Kontrakan strategis khusus mahasiswa & pekerja dekat Kampus UI.",
		Lokasi:    "Jl. Margonda Raya No. 88, Depok, Jawa Barat",
		Fasilitas: "Kasur, Lemari, Meja Belajar, WiFi Cepat, Parkir Motor.",
		Kebijakan: "Tamu menginap wajib lapor ke pemilik kontrakan.",
		Reviews: []Review{
			{Nama: "Asep C.", Rating: "⭐⭐⭐⭐⭐", Komentar: "Jalan kaki ke stasiun Pondok Cina cuma 5 menit!"},
		},
	},
	"tangerang": {
		ID:        "tangerang",
		Nama:      "Kontrakan Alam Sutera",
		Kota:      "Tangerang",
		Harga:     "Rp 2.800.000",
		Foto:      "(G4)Dafkon.png",
		Galeri:    []string{"(G4)Dafkon.png", "(G5)Dafkon.png", "(GED)Dafkon.png", "(GED)Dafkon.png", "(G1)Dafkon.png"},
		Deskripsi: "Kontrakan strategis khusus mahasiswa & pekerja dekat Pasar BSD.",
		Lokasi:    "Jl. BSD Raya, Tangerang, Banten",
		Fasilitas: "Kasur, Lemari, Meja Belajar, WiFi Cepat, Parkir Motor. Dapur, 8 Kamar",
		Kebijakan: "Dilarang membawa pacar, dilarang mabuk dan berkumpul lebih dri jam 11.",
		Reviews: []Review{
			{Nama: "Jeffrey Yanto.", Rating: "⭐⭐⭐", Komentar: "Bau rokok! saya di lantai 2 jadi ga tenang tidurnya gegara kebauan!"},
		},
	},
	"bekasi": {
		ID:        "bekasi",
		Nama:      "Kontrakan Ibu Raju",
		Kota:      "Bekasi",
		Harga:     "Rp 1.800.000",
		Foto:      "(G5)Dafkon.png",
		Galeri:    []string{"(G5)Dafkon.png", "(GED)Dafkon.png", "(G0)Dafkon.png", "(G1)Dafkon.png", "(G2)Dafkon.png"},
		Deskripsi: "Kontrakan sangat dekat dengan Sumarecon Bekasi dan dekat dengan Tol Bekasi",
		Lokasi:    "Jl. Pasar Proyek, Bekasi, Jawa Barat",
		Fasilitas: "Kasur, Lemari, Meja Belajar, WiFi Cepat, Kipas",
		Kebijakan: "Dilarang mabuk, merokok dan sebagainya. dilarang mengganggu penghuni lain.",
		Reviews: []Review{
			{Nama: "Tatang.", Rating: "⭐⭐⭐⭐", Komentar: "Asik banget! Kontrakannya dekat banget sama Summarecon Bekasi!"},
		},
	},
}

func main() {
	// Inisialisasi Koneksi PostgreSQL
	var err error
	connStr := "user=postgres password=postgres dbname=kontrakan_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Peringatan DB:", err)
	} else if err = db.Ping(); err != nil {
		log.Println("Peringatan DB tidak merespon, menggunakan data bawaan:", err)
	} else {
		fmt.Println("Berhasil terhubung ke PostgreSQL kontrakan_db!")
	}

	// Serve Static Files (Gambar, CSS)
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routing Halaman Utama & HTML
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/kontrakan/", kontrakanHandler)
	http.HandleFunc("/detail", detailHandler)
	http.HandleFunc("/tambah-kontrakan", tambahKontrakanHandler)

	fmt.Println("Server KontrakanKece jalan di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Mengarahkan root http://localhost:8080/ ke kontrakan.html
	if path == "/" {
		http.ServeFile(w, r, "kontrakan.html")
		return
	}

	// Menangani request file statis (.html, .css, .png, dll)
	if strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".jpg") {
		http.ServeFile(w, r, "."+path)
		return
	}

	http.NotFound(w, r)
}

func tambahKontrakanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Gagal memproses form", http.StatusBadRequest)
			return
		}

		nama := r.FormValue("nama")
		kota := strings.ToLower(r.FormValue("kota"))
		harga := r.FormValue("harga")
		alamat := r.FormValue("alamat")
		fasilitas := r.FormValue("fasilitas")
		kebijakan := r.FormValue("kebijakan")

		idBaru := strings.ToLower(strings.ReplaceAll(nama, " ", "-"))

		// Simpan ke database jika DB terhubung
		if db != nil {
			query := `INSERT INTO kontrakan (nama_kontrakan, kota, harga_per_bulan, alamat_lengkap, deskripsi, fasilitas, kebijakan, galeri) 
			          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
			galeriSample := pq.Array([]string{"(G4)Dafkon.png", "(G5)Dafkon.png"})
			_, errDB := db.Exec(query, nama, kota, harga, alamat, "Kontrakan baru terdaftar.", fasilitas, kebijakan, galeriSample)
			if errDB != nil {
				log.Println("Gagal simpan ke DB:", errDB)
			}
		}

		// Update / simpan ke memory map
		databaseKontrakan[kota] = Kontrakan{
			ID:        idBaru,
			Nama:      nama,
			Kota:      strings.Title(kota),
			Harga:     "Rp " + harga,
			Foto:      "(G4)Dafkon.png",
			Galeri:    []string{"(G4)Dafkon.png", "(G5)Dafkon.png"},
			Deskripsi: "Kontrakan baru terdaftar di " + strings.Title(kota),
			Lokasi:    alamat,
			Fasilitas: fasilitas,
			Kebijakan: kebijakan,
			Reviews:   []Review{},
		}

		http.Redirect(w, r, "/dafkon.html", http.StatusSeeOther)
		return
	}
	http.ServeFile(w, r, "daftar-kontrakan.html")
}

// Helper function untuk mengambil data dari PostgreSQL atau fallback ke Map
func getKontrakanData(kotaQuery string) Kontrakan {
	kotaQuery = strings.ToLower(kotaQuery)

	if db != nil {
		var k Kontrakan
		var galeriArray pq.StringArray

		query := `SELECT id_kontrakan::text, nama_kontrakan, kota, 
		                 'Rp ' || harga_per_bulan::text, alamat_lengkap, deskripsi, 
		                 fasilitas, kebijakan, galeri 
		          FROM kontrakan WHERE LOWER(kota) = $1 LIMIT 1`

		err := db.QueryRow(query, kotaQuery).Scan(
			&k.ID, &k.Nama, &k.Kota, &k.Harga, &k.Lokasi,
			&k.Deskripsi, &k.Fasilitas, &k.Kebijakan, &galeriArray,
		)

		if err == nil {
			k.Galeri = []string(galeriArray)
			if len(k.Galeri) > 0 {
				k.Foto = k.Galeri[0]
			}
			if dataMap, ok := databaseKontrakan[kotaQuery]; ok {
				k.Reviews = dataMap.Reviews
			}
			return k
		}
	}

	data, ada := databaseKontrakan[kotaQuery]
	if !ada {
		data = databaseKontrakan["jakarta"]
	}
	return data
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
	kotaQuery := strings.ToLower(r.URL.Query().Get("kota"))
	data := getKontrakanData(kotaQuery)

	tmpl, err := template.ParseFiles("detailkamar.html")
	if err != nil {
		http.Error(w, "Gagal memuat template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

func kontrakanHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "URL tidak valid", http.StatusBadRequest)
		return
	}

	idKontrakan := strings.ToLower(parts[3])
	tab := strings.ToLower(parts[4])
	data := getKontrakanData(idKontrakan)

	switch tab {
	case "review", "reviews":
		html := fmt.Sprintf("<h3>Ulasan Penghuni - %s (%s)</h3>", data.Nama, data.Kota)
		for _, rev := range data.Reviews {
			html += fmt.Sprintf(`
				<div style="border: 1px solid #ddd; padding: 12px; margin-top: 10px; border-radius: 8px; background: #fff;">
					<strong>%s</strong> <span style="color: #ffb400;">%s</span>
					<p style="margin-top: 5px; color: #444;">"%s"</p>
				</div>
			`, rev.Nama, rev.Rating, rev.Komentar)
		}
		fmt.Fprint(w, html)

	case "lokasi":
		fmt.Fprintf(w, "<h3>Lokasi %s</h3><p>📍 %s</p>", data.Nama, data.Lokasi)
	case "fasilitas":
		fmt.Fprintf(w, "<h3>Fasilitas %s</h3><p>✨ %s</p>", data.Nama, data.Fasilitas)
	case "kebijakan":
		fmt.Fprintf(w, "<h3>Kebijakan Kontrakan</h3><p>📋 %s</p>", data.Kebijakan)
	case "info":
		fallthrough
	default:
		fmt.Fprintf(w, "<h3>Deskripsi %s</h3><p>%s</p>", data.Nama, data.Deskripsi)
	}
}
