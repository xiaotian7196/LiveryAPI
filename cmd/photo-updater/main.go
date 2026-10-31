package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

const rawBase = "https://raw.githubusercontent.com/xiaotian7196/LiveryAPI/photo/photo/images/"

type photoUpdate struct {
	pattern string
	file    string
}

func main() {
	db, err := sql.Open("sqlite", "file:liveries.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	updates := []photoUpdate{
		{"China Cargo Airlines STD livery Pack - B-", "xplane-777f-china-cargo-airlines-2026.png"},
		{"FlightFactor 777F Air China Cargo B-222S", "xplane-777f-air-china-cargo-b222s-2026.png"},
		{"FlightFactor 767-200 Japan Transocean Air", "xplane-767-200-japan-transocean-air-ja8231-2026.png"},
		{"X-Crafts E-Jets Family E175 Fuji Dream Airlines", "xplane-e175-fuji-dream-airlines-ja12fj-2026.png"},
		{"FlightFactor 777F Air China Cargo B-2097", "xplane-777f-air-china-cargo-b2097-2026.png"},
		{"Toliss A320neo WestAir STD livery", "xplane-a320neo-westair-b306t-2024.png"},
		{"Toliss A320neo China Eastern airline STD livery", "xplane-a320neo-china-eastern-b322j-b32ep-2024.jpg"},
		{"PMDG 747-4B5F Grandstar Cargo B-2427", "p3d-pmdg-747-4b5f-grandstar-cargo-b2427-2024.jpg"},
		{"PMDG 747-40BERF China Cargo Airline B-2425", "p3d-pmdg-747-40berf-china-cargo-b2425-2024.png"},
		{"PMDG 737-84P Hainan Airlines OLD B-2157", "p3d-pmdg-737-84p-hainan-b2157-2024.png"},
		{"iFly 737-MAX8 China Eastern Yunnan", "p3d-ifly-737-max8-china-eastern-yunnan-2024.png"},
		{"PMDG NG3 737-881 ANA EXPO 2025 OSAKA JA58AN", "mfs2020-pmdg-737-881-ana-ja58an-2024.png"},
		{"PMDG 737-89P China Eastern STD Livery B-5838", "mfs2020-pmdg-737-89p-china-eastern-b5838-2024.png"},
	}

	updated := 0
	for _, item := range updates {
		result, err := db.Exec(`UPDATE liveries SET image_url = ? WHERE name LIKE ?`, rawBase+item.file, "%"+item.pattern+"%")
		if err != nil {
			log.Fatalf("update %q: %v", item.pattern, err)
		}
		count, _ := result.RowsAffected()
		updated += int(count)
		fmt.Printf("%s: %d records\n", item.file, count)
	}
	fmt.Printf("updated %d livery records\n", updated)
}
