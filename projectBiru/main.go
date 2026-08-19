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

var db *sql.DB

// Menyimpan list kontrakan per kota
var databaseKontrakan = map[string][]Kontrakan{
	"jakarta": {
		{
			ID:        "jakarta-ancol",
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
			},
		},
	},
	"bogor": {
		{
			ID:        "bogor-mbah-rusdi",
			Nama:      "Kontrakan Mbah Rusdi",
			Kota:      "Bogor",
			Harga:     "Rp 1.500.000",
			Foto:      "(G2)Dafkon.png",
			Galeri:    []string{"(G2)Dafkon.png", "(G3)Dafkon.png", "(G4)Dafkon.png"},
			Deskripsi: "Udara sejuk khas Puncak Cisarua.",
			Lokasi:    "Jl. Cisarua Raya No. 45, Bogor",
			Fasilitas: "WiFi, Saung Santai, Dapur Bersama.",
			Kebijakan: "Dilarang membuat keributan.",
			Reviews:   []Review{},
		},
	},
	"depok": {
		{
			ID:        "depok-biru",
			Nama:      "Kontrakan Biru",
			Kota:      "Depok",
			Harga:     "Rp 2.000.000",
			Foto:      "(G3)Dafkon.png",
			Galeri:    []string{"(G3)Dafkon.png", "(G0)Dafkon.png"},
			Deskripsi: "Kontrakan strategis dekat UI.",
			Lokasi:    "Jl. Margonda Raya No. 88, Depok",
			Fasilitas: "Kasur, Lemari, WiFi",
			Kebijakan: "Tamu wajib lapor.",
			Reviews:   []Review{},
		},
	},
	"tangerang": {
		{
			ID:        "tangerang-alam-sutera",
			Nama:      "Kontrakan Alam Sutera",
			Kota:      "Tangerang",
			Harga:     "Rp 2.800.000",
			Foto:      "(G4)Dafkon.png",
			Galeri:    []string{"(G4)Dafkon.png", "(G5)Dafkon.png"},
			Deskripsi: "Kontrakan strategis dekat Pasar BSD.",
			Lokasi:    "Jl. BSD Raya, Tangerang",
			Fasilitas: "WiFi, Parkir Motor",
			Kebijakan: "Dilarang membuat keributan.",
			Reviews:   []Review{},
		},
	},
	"bekasi": {
		{
			ID:        "bekasi-ibu-raju",
			Nama:      "Kontrakan Ibu Raju",
			Kota:      "Bekasi",
			Harga:     "Rp 1.800.000",
			Foto:      "(G5)Dafkon.png",
			Galeri:    []string{"(G5)Dafkon.png", "(GED)Dafkon.png"},
			Deskripsi: "Dekat Summarecon Bekasi.",
			Lokasi:    "Jl. Pasar Proyek, Bekasi",
			Fasilitas: "Kipas, Lemari, WiFi",
			Kebijakan: "Dilarang merokok.",
			Reviews:   []Review{},
		},
	},
}

func main() {
	connStr := "user=postgres password=postgres dbname=kontrakan_db sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil || db.Ping() != nil {
		log.Println("Peringatan DB tidak terhubung, menggunakan memori lokal.")
	}

	fs := http.FileServer(http.Dir("./"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/dafkon.html", dafkonPageHandler)
	http.HandleFunc("/dafkot.html", dafkotPageHandler)
	http.HandleFunc("/api/kontrakan/", kontrakanHandler)
	http.HandleFunc("/detail", detailHandler)
	http.HandleFunc("/tambah-kontrakan", tambahKontrakanHandler)

	fmt.Println("Server KontrakanKece jalan di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func getAllKontrakanList() []Kontrakan {
	var list []Kontrakan
	for _, items := range databaseKontrakan {
		list = append(list, items...)
	}
	return list
}

func dafkonPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("dafkon.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, getAllKontrakanList())
}

func dafkotPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("dafkot.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, databaseKontrakan)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		http.ServeFile(w, r, "kontrakan.html")
		return
	}

	// Cegah detailkamar.html dibaca sebagai file statis biasa
	if path == "/detailkamar.html" {
		http.Redirect(w, r, "/detail?"+r.URL.RawQuery, http.StatusSeeOther)
		return
	}

	if strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".jpg") {
		http.ServeFile(w, r, "."+path)
		return
	}
	http.NotFound(w, r)
}

func tambahKontrakanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		nama := r.FormValue("nama")
		kota := strings.ToLower(r.FormValue("kota"))
		harga := r.FormValue("harga")
		alamat := r.FormValue("alamat")
		fasilitas := r.FormValue("fasilitas")
		kebijakan := r.FormValue("kebijakan")

		idBaru := strings.ToLower(strings.ReplaceAll(nama, " ", "-"))

		kontrakanBaru := Kontrakan{
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

		// Tambahkan ke map memori
		databaseKontrakan[kota] = append(databaseKontrakan[kota], kontrakanBaru)

		// Simpan ke PostgreSQL jika terhubung
		if db != nil && db.Ping() == nil {
			query := `INSERT INTO kontrakan (nama_kontrakan, kota, harga_per_bulan, alamat_lengkap, deskripsi, fasilitas, kebijakan, galeri) 
			          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
			db.Exec(query, nama, kota, harga, alamat, kontrakanBaru.Deskripsi, fasilitas, kebijakan, pq.Array(kontrakanBaru.Galeri))
		}

		// Redirect langsung ke daftar kontrakan
		http.Redirect(w, r, "/dafkon.html", http.StatusSeeOther)
		return
	}
	http.ServeFile(w, r, "daftar-kontrakan.html")
}

func getKontrakanData(idOrKota string) Kontrakan {
	idOrKota = strings.ToLower(idOrKota)
	for _, list := range databaseKontrakan {
		for _, k := range list {
			if strings.ToLower(k.ID) == idOrKota || strings.ToLower(k.Kota) == idOrKota {
				return k
			}
		}
	}
	return databaseKontrakan["jakarta"][0]
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
	idQuery := strings.ToLower(r.URL.Query().Get("id"))
	if idQuery == "" {
		idQuery = strings.ToLower(r.URL.Query().Get("kota"))
	}

	// Mengambil data kontrakan dari memori
	data := getKontrakanData(idQuery)

	// GUNAKAN template.ParseFiles, BUKAN http.ServeFile
	tmpl, err := template.ParseFiles("detailkamar.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send data struct Kontrakan ke template HTML
	tmpl.Execute(w, data)
}

func kontrakanHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
		if len(data.Reviews) == 0 {
			html += "<p style='color: #666;'>Belum ada ulasan untuk kontrakan ini.</p>"
		} else {
			for _, rev := range data.Reviews {
				html += fmt.Sprintf(`
					<div style="border: 1px solid #ddd; padding: 12px; margin-top: 10px; border-radius: 8px; background: #fff;">
						<strong>%s</strong> <span style="color: #ffb400;">%s</span>
						<p style="margin-top: 5px; color: #444;">"%s"</p>
					</div>
				`, rev.Nama, rev.Rating, rev.Komentar)
			}
		}
		fmt.Fprint(w, html)
	case "lokasi":
		fmt.Fprintf(w, "<h3>Lokasi %s</h3><p>📍 %s</p>", data.Nama, data.Lokasi)
	case "fasilitas":
		fmt.Fprintf(w, "<h3>Fasilitas %s</h3><p>✨ %s</p>", data.Nama, data.Fasilitas)
	case "kebijakan":
		fmt.Fprintf(w, "<h3>Kebijakan Kontrakan</h3><p>📋 %s</p>", data.Kebijakan)
	default:
		fmt.Fprintf(w, "<h3>Deskripsi %s</h3><p>%s</p>", data.Nama, data.Deskripsi)
	}
}
