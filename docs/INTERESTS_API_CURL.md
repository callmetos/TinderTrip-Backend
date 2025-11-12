# Interests API - cURL Commands

## 🔓 Public Endpoints

### 1. Get All Interests (Global)
```bash
# Get all interests
curl -X 'GET' \
  'https://api.tindertrip.phitik.com/api/v1/interests' \
  -H 'accept: application/json'

# Filter by category
curl -X 'GET' \
  'https://api.tindertrip.phitik.com/api/v1/interests?category=cafe' \
  -H 'accept: application/json'
```

---

## 🔒 User Interests (Authentication Required)

### 2. Get User Interests (with selection status)
```bash
curl -X 'GET' \
  'https://api.tindertrip.phitik.com/api/v1/users/interests' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer YOUR_TOKEN'

# Filter by category
curl -X 'GET' \
  'https://api.tindertrip.phitik.com/api/v1/users/interests?category=activity' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer YOUR_TOKEN'
```

### 3. Get User Selected Interests Only
```bash
curl -X 'GET' \
  'https://api.tindertrip.phitik.com/api/v1/users/interests/selected' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer YOUR_TOKEN'
```

### 4. Update User Interests (PUT - Bulk Replace)
```bash
curl -X 'PUT' \
  'https://api.tindertrip.phitik.com/api/v1/users/interests' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer YOUR_TOKEN' \
  -H 'Content-Type: application/json' \
  -d '{
  "interest_codes": [
    "coffee",
    "bubble_tea",
    "reading",
    "running",
    "football"
  ]
}'
```

**Note:** ไม่มี POST endpoint แยก - ใช้ PUT (bulk replace) แทน

---

## 🔒 Event Interests (Authentication Required)

Event interests จัดการผ่าน event endpoints

### 5. Get Event Interests
```bash
curl -X 'GET' \
  'https://api.tindertrip.phitik.com/api/v1/events/EVENT_ID' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer YOUR_TOKEN'
```

**Response จะมี `interests` array ใน event object**

### 6. Create Event with Interests (POST)
```bash
curl -X 'POST' \
  'https://api.tindertrip.phitik.com/api/v1/events' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer YOUR_TOKEN' \
  -H 'Content-Type: application/json' \
  -d '{
  "title": "Coffee & Reading Meetup",
  "description": "Join us for coffee and reading",
  "event_type": "meal",
  "start_at": "2025-12-15T10:00:00+07:00",
  "end_at": "2025-12-15T12:00:00+07:00",
  "capacity": 10,
  "budget_min": 200,
  "budget_max": 500,
  "currency": "THB",
  "address_text": "Central World",
  "interest_codes": [
    "coffee",
    "reading",
    "chilling"
  ]
}'
```

### 7. Update Event Interests (PUT)
```bash
curl -X 'PUT' \
  'https://api.tindertrip.phitik.com/api/v1/events/EVENT_ID' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer YOUR_TOKEN' \
  -H 'Content-Type: application/json' \
  -d '{
  "interest_codes": [
    "coffee",
    "bubble_tea",
    "reading",
    "running"
  ]
}'
```

**Note:** 
- ไม่มี POST/PUT endpoint แยกสำหรับ event interests
- ใช้ `interest_codes` ใน CreateEvent และ UpdateEvent
- PUT จะ replace interests ทั้งหมด (bulk replace)

---

## 📋 Interest Codes Reference

### Restaurant (28)
`fast_food`, `noodles`, `grill`, `pasta`, `dim_sum`, `indian_food`, `salads`, `japanese_food`, `izakaya`, `muu_kra_ta`, `street_food`, `pork`, `pizza`, `vegan`, `chinese_food`, `sushi`, `fine_dining`, `halal`, `burger`, `korean_food`, `buffet`, `ramen`, `bbq`, `meat`, `healthy_food`, `shabu_sukiyaki_hot_pot`, `omakase`, `seafood`

### Cafe (7)
`bubble_tea`, `bingsu`, `matcha`, `bakery_cake`, `ice_cream`, `pancakes`, `coffee`

### Activity (33)
`chilling`, `painting`, `baking`, `investing`, `fan_meet`, `shopping`, `wakeboard`, `laser_tag`, `superstition`, `bb_gun`, `travel`, `photography`, `temple`, `night_market`, `park`, `amusement_park`, `movies`, `karaoke`, `running`, `art_gallery`, `archery`, `scuba_diving`, `reading`, `skateboard`, `walking`, `volunteer`, `boardgame`, `paintball`, `museum`, `flower_arrangement`, `kpop`, `concert`, `aquarium`

### Pub & Bar (16)
`wine`, `ratchathewi`, `khaosan_road`, `thonglor`, `thai_music`, `edm_music`, `heartbroken`, `beer_tower`, `jazz`, `cocktail_bar`, `live_music`, `rca_plaza`, `kpop_music`, `pubs_bars`, `rooftop`, `alcohol`

### Sport (15)
`football`, `sports`, `snooker`, `rock_climbing`, `golf`, `boxing`, `fitness`, `basketball`, `bowling`, `tennis`, `badminton`, `volleyball`, `racquet`, `table_tennis`, `yoga`
