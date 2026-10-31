package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type livery struct {
	name, description, game, aircraft, cnURL, globalURL, version string
}

func main() {
	db, err := sql.Open("sqlite", "file:liveries.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	items := []livery{
		{"[Xplane] FlightFactor 777F China Cargo Airlines STD livery Pack - B-222K [GE]", "China Cargo Airlines STD livery, GE90-110B1", "Xplane 11/12", "FlightFactor 777F B-222K", "https://livery.cn-sy1.rains3.com/%5B200F%20GE%5D%20CHINA%20CARGO%20AIRLINES%20B-222K.zip", "https://drive.google.com/file/d/1VQ3uLFdzKte0h0G0DiA3f5gjowJR20VB/view?usp=sharing", ""},
		{"[Xplane] FlightFactor 777F China Cargo Airlines STD livery Pack - B-223A [GE]", "China Cargo Airlines STD livery, GE90-115B", "Xplane 11/12", "FlightFactor 777F B-223A", "https://livery.cn-sy1.rains3.com/%5B200F%20GE%5D%20CHINA%20CARGO%20AIRLINES%20B-223A.zip", "https://drive.google.com/file/d/12sSTBv3bIP_BXBn-IrJlo0tVXzaQh-yI/view?usp=sharing", ""},
		{"[Xplane] FlightFactor 777F China Cargo Airlines STD livery Pack - B-227M [GE]", "China Cargo Airlines STD livery, GE90-115B", "Xplane 11/12", "FlightFactor 777F B-227M", "https://livery.cn-sy1.rains3.com/%5B200F%20GE%5D%20CHINA%20CARGO%20AIRLINES%20B-227M.zip", "https://drive.google.com/file/d/1t0nrf-1rJtzEIzvFN_TYACOVqKJvwb0A/view?usp=sharing", ""},
		{"[Xplane] FlightFactor 777F Air China Cargo B-222S [GE]", "Air China Cargo B-222S B777-FFT, GE90-115B", "Xplane 11/12", "FlightFactor 777F B-222S", "https://livery.cn-sy1.rains3.com/%5B200F%20GE%5D%20AIR%20CHINA%20CARGO%20B-222S.zip", "https://drive.google.com/file/d/1upbCqGx9v7rNRTGM_JrKZv4sW9ZEBRIc/view?usp=sharing", ""},
		{"FlightFactor 767-200 Japan Transocean Air 1994-STD JA8231 [PW]", "Japan Transocean Air 1994-STD", "Xplane 11/12", "FlightFactor 767-200 JA8231", "https://livery.cn-sy1.rains3.com/762%20PW%20JTA%20%28JA8231%29.zip", "https://drive.google.com/file/d/16RuNQ6qWIKi-OAi8y5EsBK3btU4473kW/view?usp=sharing", ""},
		{"[Xplane] X-Crafts E-Jets Family E175 Fuji Dream Airlines JA12FJ", "Fuji Dream Airlines JA12FJ; extended winglets should be enabled in EFB settings", "Xplane 11/12", "X-Crafts E-Jets Family E175 JA12FJ", "https://livery.cn-sy1.rains3.com/%5B175%5DFuji%20Dream%20Airlines%20JA12FJ.zip", "https://drive.google.com/file/d/1fulsM7IgV_bMUoducbnHOER9l3xdMHy6/view?usp=sharing", ""},
		{"[Xplane] FlightFactor 777F Air China Cargo B-2097 [GE] [NO NATIONAL FLAG]", "Air China Cargo B-2097 B777-FFT, no national flag; V1.1 engine textures fix", "Xplane 11/12", "FlightFactor 777F B-2097", "https://livery.cn-sy1.rains3.com/%5B200F%20GE%5D%20AIR%20CHINA%20CARGO%20B-2097.zip", "https://drive.google.com/file/d/19OMEjDNiYQiemt_Qr1BngS6-QKFzE1Kp/view?usp=sharing", "V1.1"},
		{"[Xplane] Toliss A320neo WestAir STD livery - B-306T normal v2.5", "WestAir A320neo STD livery, PW1100 engines", "Xplane 11/12", "Toliss A320neo B-306T normal", "https://cn-sy1.rains3.com/livery/[PWG]%20West%20Air%20B-306T.zip", "https://drive.google.com/file/d/1Kk7osn8Q5eEYIMpiWEXCZgmr3y8NNkdD/view?usp=sharing", "2.5"},
		{"[Xplane] Toliss A320neo WestAir STD livery - B-306T white nose v2.5", "WestAir A320neo STD livery, white nose, PW1100 engines", "Xplane 11/12", "Toliss A320neo B-306T white nose", "https://cn-sy1.rains3.com/livery/[PWG]%20West%20Air%20%28White%20nose%29%20B-306T.zip", "https://drive.google.com/file/d/1ChJgrkumIi12uFiVVz8_92d9zyYjlXuE/view?usp=sharing", "2.5"},
		{"[Xplane] Toliss A320neo China Eastern airline STD livery - B-322J", "China Eastern airline STD livery, LEAP engines", "Xplane 11/12", "Toliss A320neo B-322J", "https://cn-sy1.rains3.com/livery/[LEAP]%20CES%20B-322J.zip", "https://drive.google.com/file/d/1wZQCKB5rvceRjmYCsJxpPhCcK6P1cym2/view?usp=sharing", "3.0"},
		{"[Xplane] Toliss A320neo China Eastern airline STD livery - B-32EP", "China Eastern airline STD livery, LEAP engines", "Xplane 11/12", "Toliss A320neo B-32EP", "https://cn-sy1.rains3.com/livery/[LEA]%20CES%20B-32EP.zip", "https://drive.google.com/file/d/1uujq0YkMlvuiHTIvP8pa0VQoE0bUtjul/view?usp=sharing", "3.0"},
		{"[P3D] PMDG 747-4B5F Grandstar Cargo B-2427 ZIP", "Grandstar Cargo; supports Prepar 3D v4/5", "Prepar 3D v4/5", "PMDG 747-4B5F B-2427", "https://cn-sy1.rains3.com/livery/texture.B2427.zip", "", ""},
		{"[P3D] PMDG 747-4B5F Grandstar Cargo B-2427 PTP", "Grandstar Cargo paint package; supports Prepar 3D v4/5", "Prepar 3D v4/5", "PMDG 747-4B5F B-2427", "https://cn-sy1.rains3.com/livery/B-2427.ptp", "", ""},
		{"[P3D] PMDG 747-40BERF China Cargo Airline B-2425", "China Cargo Airline; supports Prepar 3D v4/5", "Prepar 3D v4/5", "PMDG 747-40BERF B-2425", "https://cn-sy1.rains3.com/livery/texture.B2425.zip", "", ""},
		{"[P3D] PMDG 737-84P Hainan Airlines OLD B-2157", "Hainan Airlines old livery; supports PMDG 737-800NGXu for Prepar 3D v4/5", "Prepar 3D v4/5", "PMDG 737-84P B-2157", "https://cn-sy1.rains3.com/livery/texture.B-2157.zip", "", ""},
		{"[P3D] iFly 737-MAX8 China Eastern Yunnan B-1380", "China Eastern Yunnan livery; supports iFly 737-MAX8 for Prepar 3D v5", "Prepar 3D v5", "iFly 737-MAX8 B-1380", "https://cn-sy1.rains3.com/livery/iFly%20737-MAX8%20P3D%20v5%20CHINA%20EASTERN%20YUNNAN%20B-1380.zip", "", ""},
		{"[P3D] iFly 737-MAX8 China Eastern Yunnan B-1383", "China Eastern Yunnan livery; available through iFly's livery manager", "Prepar 3D v5", "iFly 737-MAX8 B-1383", "https://cn-sy1.rains3.com/livery/iFly%20737-MAX8%20P3D%20v5%20CHINA%20EASTERN%20YUNNAN%20B-1383.zip", "", ""},
		{"[P3D] iFly 737-MAX8 China Eastern Yunnan B-1385", "China Eastern Yunnan livery; supports iFly 737-MAX8 for Prepar 3D v5", "Prepar 3D v5", "iFly 737-MAX8 B-1385", "https://cn-sy1.rains3.com/livery/iFly%20737-MAX8%20P3D%20v5%20CHINA%20EASTERN%20YUNNAN%20B-1385.zip", "", ""},
		{"[MFS2020] PMDG NG3 737-881 ANA EXPO 2025 OSAKA JA58AN", "ANA EXPO 2025 OSAKA; install using PMDG OC2", "MFS 2020", "PMDG NG3 737-881 JA58AN", "https://cn-sy1.rains3.com/livery/JA58AN.zip", "https://drive.google.com/file/d/1NUmpuup7ZVawVRhv2zh-lcyX3Khh42iY/view?usp=sharing", ""},
		{"[MFS2020] PMDG 737-89P China Eastern STD Livery B-5838", "China Eastern STD livery; supports PMDG 737-800NG3 for MSFS", "MFS 2020", "PMDG 737-89P B-5838", "https://cn-sy1.rains3.com/livery/B-5838_ku3OX.zip", "", "1.5F"},
	}

	stmt, err := db.Prepare(`INSERT INTO liveries (name, description, game, aircraft, download_url_cn, download_url_global, version)
SELECT ?, ?, ?, ?, ?, ?, ?
WHERE NOT EXISTS (SELECT 1 FROM liveries WHERE name = ? AND download_url_cn = ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	inserted := 0
	for _, item := range items {
		result, err := stmt.Exec(item.name, item.description, item.game, item.aircraft, item.cnURL, item.globalURL, item.version, item.name, item.cnURL)
		if err != nil {
			log.Fatalf("insert %q: %v", item.name, err)
		}
		count, _ := result.RowsAffected()
		inserted += int(count)
	}
	fmt.Printf("imported %d new liveries, processed %d records\n", inserted, len(items))
}
