# Test Backend Saiful Anwar

# Langkah Penggunaan

## Requirements
- Golang versi 1.24.4
- PostgreSQL

## Langkah
- Buat database di postgre dan jalankan query yang ada di file src/migrations/init.up.sql secara manual (sorry, no time for migrations setup)
- Buat file .env di root project dan sesuaikan value-nya (template dari .env.example)
- Dari folder root project jalankan
  ```console
  go mod tidy
  go run src/main.go
  ```
- App akan berjalan di port yang diisi di .env

# Desain Sistem
- Terdapat 4 tabel
  - yards, untuk menyimpan data yard.
  - blocks, untuk menyimpan data block. Terdapat field yard_id yang terkoneksi ke yards.
  - yard_plans, untuk menyimpan data planning. Terdapat field block_id yang terkoneksi ke blocks.
  - placements, untuk menyimpan data penempatan container. Container dengan size 40ft akan disimpan ke dalam 2 row di dalam tabel placements (untuk kemudahan).

# Asumsi
- Suggestion & placement hanya berlaku di area yang sudah ditambahkan ke yard_plans
- Apabila terdapat beberapa yard_plans yang cocok untuk suggestion, maka yard_plans yang diambil adalah yang pertama kali dibuat (smallest id)
- Rule placement sangat simple
  - Container bisa ditaruh di tier yang lebih tinggi walaupun di bawahnya tidak ada container lain (bisa melayang di udara)
  - Tidak re-placement secara otomatis kalau pickup tier tertentu yang tier atasnya ada container lain
- Block bisa punya ukuran (slot, row, tier) yang berbeda di dalam satu yard yang sama
- Ukuran container 40ft mengisi 2 slot
- yard_plans hanya bisa dipakai untuk container dengan type tertentu saja. Tidak ada aturan khusus lainnya
- Type dari container hanya: DRY, REEFER, OPEN_TOP, FLAT_RACK

# API Docs

## Yards

### List Yards
**GET** `/yard`

Retrieve all yards.

### Create Yard  
**POST** `/yard`

Request body:
```json
{
  "name": "Yard 2",
  "description": "The yard no 2"
}
```

### Edit Yard  
**PUT** `/yard/{id}`

Request body:
```json
{
  "name": "YRD-02",
  "description": "The yard no 1"
}
```

Example URL:
```
/yard/2
```

### Delete Yard  
**DELETE** `/yard/{id}`

Example URL:
```
/yard/1
```

---

## Blocks

### List Blocks  
**GET** `/block`

### List Blocks By Yard  
**GET** `/block/by_yard/{yard_id}`

Example:
```
/block/by_yard/1
```

### Create Block  
**POST** `/block`

Request body:
```json
{
  "yard_id": 1,
  "name": "A-1",
  "slots": 12,
  "rows": 6,
  "tiers": 5
}
```

### Edit Block  
**PUT** `/block/{id}`

Request body:
```json
{
  "yard_id": "3045f3f6-f9e4-4940-8116-62becc0cf91d",
  "name": "A-1",
  "slots": 12,
  "rows": 6,
  "tiers": 5
}
```

Example URL:
```
/block/efe159ed-50e9-4c78-af9c-586a256fafba
```

### Delete Block  
**DELETE** `/block/{id}`

Example URL:
```
/block/17bf872a-3737-4e03-8dec-5424e72319ba
```

---

## Yard Plans

### List Yard Plans  
**GET** `/yard_plan`

### List Yard Plans By Yard  
**GET** `/yard_plan/by_yard/{yard_id}`

Example:
```
/yard_plan/by_yard/1
```

### List Yard Plans By Block  
**GET** `/yard_plan/by_block/{block_id}`

Example:
```
/yard_plan/by_block/1
```

### Create Yard Plan  
**POST** `/yard_plan`

Request body:
```json
{
  "block_id": 1,
  "slot_start": 1,
  "slot_end": 2,
  "row_start": 1,
  "row_end": 2,
  "size": 20,
  "height": 10,
  "type": "DRY",
  "slot_priority": -1,
  "row_priority": 0,
  "tier_priority": 0
}
```

### Edit Yard Plan  
**PUT** `/block/{id}`  
> Note: Endpoint in Postman seems incorrect. Should be `/yard_plan/{id}`.

### Delete Yard Plan  
**DELETE** `/yard_plan/{id}`

Example:
```
/yard_plan/2
```

### Suggest Container Placement  
**POST** `/suggestion`

Request body:
```json
{
  "yard_id": 1,
  "container_id": "C001",
  "container_size": 20,
  "container_height": 8.6,
  "container_type": "DRY"
}
```

### Place Container  
**POST** `/place`

Request body:
```json
{
  "container_id": "C001",
  "container_size": 20,
  "container_height": 8.6,
  "container_type": "DRY",
  "block_id": 1,
  "slot": 1,
  "row": 1,
  "tier": 1
}
```

### Pickup Container  
**POST** `/pickup`

Request body:
```json
{
  "container_id": "C001"
}
```
