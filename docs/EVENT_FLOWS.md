# 📋 Event Flows - ทุกเคสที่เป็นไปได้

เอกสารนี้สรุป flow ทั้งหมดของ Event ที่เป็นไปได้ในระบบ TinderTrip

---

## 📊 สารบัญ

1. [Event Lifecycle (Status)](#1-event-lifecycle-status)
2. [Event Types](#2-event-types)
3. [Event Member Lifecycle](#3-event-member-lifecycle)
4. [Event Swipe Flow](#4-event-swipe-flow)
5. [Event Operations (CRUD)](#5-event-operations-crud)
6. [Event Image/Photo Operations](#6-event-imagephoto-operations)
7. [Event Discovery & Suggestions](#7-event-discovery--suggestions)
8. [Error Cases & Edge Cases](#8-error-cases--edge-cases)

---

## 1. Event Lifecycle (Status)

### 1.1 Event Status Types

| Status | Description | Can Update? | Can Delete? |
|--------|-------------|-------------|-------------|
| `published` | Event ที่เผยแพร่แล้ว (default) | ✅ Yes (Creator only) | ✅ Yes (Creator only) |
| `cancelled` | Event ที่ถูกยกเลิก | ✅ Yes (Creator only) | ✅ Yes (Creator only) |
| `completed` | Event ที่เสร็จสิ้นแล้ว | ❌ No | ❌ No |

### 1.2 Status Transitions

```
┌─────────────┐
│  published  │ ← (Default เมื่อสร้าง event)
└──────┬──────┘
       │
       ├───> PUT /events/{id} {status: "cancelled"}
       │     └──> [cancelled]
       │
       ├───> POST /events/{id}/complete (Creator only)
       │     └──> [completed]
       │           └──> สร้าง UserEventHistory สำหรับทุก confirmed members
       │
       └───> DELETE /events/{id} (Creator only)
             └──> [soft delete] (deleted_at set)
```

### 1.3 Status Flow Cases

#### Case 1: Create Event
```
POST /events
→ Status: published (default)
→ Creator automatically added as member (confirmed)
→ Chat room created automatically
```

#### Case 2: Cancel Event
```
PUT /events/{id}
Body: {status: "cancelled"}
→ Status changed: published → cancelled
→ Members can still see event but status is cancelled
```

#### Case 3: Complete Event (Creator only)
```
POST /events/{id}/complete
→ Status changed: published → completed
→ Creates UserEventHistory for all confirmed members
→ Event is locked (cannot update/delete)
```

#### Case 4: Delete Event (Soft Delete)
```
DELETE /events/{id}
→ Sets deleted_at timestamp
→ Event hidden from queries (WHERE deleted_at IS NULL)
→ Related data preserved (members, swipes, history)
```

---

## 2. Event Types

### 2.1 Available Event Types

| Type | Value | Description |
|------|-------|-------------|
| Meal | `meal` | มื้ออาหาร |
| Day Trip | `daytrip` / `one_day_trip` | ทริปไปเที่ยววันเดียว |
| Overnight | `overnight` | ทริปค้างคืน |
| Activity | `activity` | กิจกรรมอื่นๆ |
| Other | `other` | อื่นๆ |

**Note:** ใน API request ใช้ `one_day_trip` แต่ใน DB/model เก็บเป็น `daytrip`

### 2.2 Event Type Validation

```go
// CreateEventRequest
event_type: "required,oneof=meal one_day_trip overnight"

// UpdateEventRequest  
event_type: "omitempty,oneof=meal daytrip overnight activity other"
```

---

## 3. Event Member Lifecycle

### 3.1 Member Roles

| Role | Description | Status on Creation |
|------|-------------|-------------------|
| `creator` | ผู้สร้าง event | `confirmed` (auto) |
| `participant` | ผู้เข้าร่วม | `pending` (default) |

### 3.2 Member Status Types

| Status | Description | Can Confirm? | Can Leave? |
|--------|-------------|--------------|------------|
| `pending` | รอการยืนยัน | ✅ Yes | ❌ No |
| `confirmed` | ยืนยันแล้ว | ❌ No | ✅ Yes |
| `declined` | ปฏิเสธ | ❌ No | ❌ No |
| `left` | ออกจาก event | ❌ No | ❌ No |
| `kicked` | ถูกไล่ออก | ❌ No | ❌ No |

### 3.3 Member Flow Cases

#### Case 1: Creator Joins (Automatic)
```
POST /events (Create Event)
→ Creator automatically added as:
  - Role: creator
  - Status: confirmed
  - JoinedAt: now()
  - ConfirmedAt: now()
```

#### Case 2: User Joins via Join Endpoint
```
POST /events/{id}/join
→ Creates EventMember:
  - Role: participant
  - Status: pending
  - JoinedAt: now()
  
Possible Errors:
- 400: Event not found
- 409: User is already a member
```

#### Case 3: User Swipes Like
```
POST /events/{id}/swipe
Body: {direction: "like"}
→ Creates EventSwipe (like)
→ Creates EventMember (if not exists):
  - Role: participant
  - Status: pending
```

#### Case 4: User Confirms Participation
```
POST /events/{id}/confirm
→ Updates EventMember:
  - Status: pending → confirmed
  - ConfirmedAt: now()
  
Possible Errors:
- 404: Event not found
- 404: Member not found
- 409: Event is full (capacity check)
```

#### Case 5: User Cancels Participation
```
POST /events/{id}/cancel
→ Updates EventMember:
  - Status: pending → declined
  
Possible Errors:
- 404: Event not found
- 404: Member not found
```

#### Case 6: User Leaves Event
```
POST /events/{id}/leave
→ Updates EventMember:
  - Status: confirmed → left
  - LeftAt: now()
  
Possible Errors:
- 404: Event not found
- 404: User is not a member
```

### 3.4 Member Status Transitions

```
┌──────────┐
│ pending  │ (Default when joining)
└────┬─────┘
     │
     ├───> POST /events/{id}/confirm
     │     └──> [confirmed]
     │           └──> Can leave: POST /events/{id}/leave
     │                 └──> [left]
     │
     └───> POST /events/{id}/cancel
           └──> [declined]
```

### 3.5 Capacity Management

```go
// When confirming participation
if event.Capacity != nil {
    confirmedCount := COUNT(members WHERE status = 'confirmed')
    if confirmedCount >= event.Capacity {
        return error("event is full")
    }
}
```

---

## 4. Event Swipe Flow

### 4.1 Swipe Directions

| Direction | Value | Action |
|-----------|-------|--------|
| Like | `like` | สนใจ / ชอบ event → สร้าง member (pending) |
| Pass | `pass` | ไม่สนใจ / ข้าม event |

### 4.2 Swipe Flow Cases

#### Case 1: Swipe Like (First Time)
```
POST /events/{id}/swipe
Body: {direction: "like"}
→ Creates EventSwipe (like)
→ Creates EventMember (pending)
```

#### Case 2: Swipe Like (Already Swiped)
```
POST /events/{id}/swipe
Body: {direction: "like"}
→ Updates EventSwipe (like)
→ No new member created (already exists)
```

#### Case 3: Swipe Pass
```
POST /events/{id}/swipe
Body: {direction: "pass"}
→ Creates/Updates EventSwipe (pass)
→ No member created/updated
```

#### Case 4: Change Mind (Pass → Like)
```
POST /events/{id}/swipe {direction: "pass"}
→ EventSwipe: pass

POST /events/{id}/swipe {direction: "like"}
→ EventSwipe: pass → like
→ Creates EventMember (pending) if not exists
```

### 4.3 Swipe vs Join

| Action | Endpoint | Creates Swipe? | Creates Member? | Member Status |
|--------|----------|----------------|-----------------|---------------|
| Swipe Like | `POST /events/{id}/swipe {like}` | ✅ Yes | ✅ Yes | `pending` |
| Join | `POST /events/{id}/join` | ❌ No | ✅ Yes | `pending` |
| Swipe Pass | `POST /events/{id}/swipe {pass}` | ✅ Yes | ❌ No | N/A |

**Note:** Swipe like และ Join ทำหน้าที่คล้ายกัน แต่ Swipe มี record ของ swipe action ด้วย

---

## 5. Event Operations (CRUD)

### 5.1 Create Event

#### Endpoint: `POST /events`

**JSON Mode:**
```json
{
  "title": "Weekend Brunch",
  "description": "Let's enjoy brunch together!",
  "event_type": "meal",
  "address_text": "123 Sukhumvit Rd",
  "lat": 13.7563,
  "lng": 100.5018,
  "start_at": "2024-12-25T10:00:00Z",
  "end_at": "2024-12-25T12:00:00Z",
  "capacity": 6,
  "budget_min": 300,
  "budget_max": 500,
  "currency": "THB",
  "category_ids": ["uuid1", "uuid2"],
  "tag_ids": ["uuid3"]
}
```

**Multipart Mode:**
- All fields as form data
- `file`: Cover image (optional)
- `files[]`: Event photos (multiple, optional)

**What Happens:**
1. Event created with status `published`
2. Creator added as member (confirmed)
3. Chat room created automatically
4. Categories/Tags linked
5. Images uploaded (if provided)

**Possible Errors:**
- 400: Invalid request data
- 401: Not authenticated

---

### 5.2 Read Events

#### Endpoint: `GET /events`

**Query Parameters:**
- `page` (default: 1)
- `limit` (default: 10)
- `event_type` (filter: meal, daytrip, overnight)
- `status` (filter: published, cancelled, completed)

**Returns:** Paginated list of events

---

#### Endpoint: `GET /events/{id}`

**Returns:** Single event with full details

**Possible Errors:**
- 400: Invalid event ID
- 404: Event not found

---

#### Endpoint: `GET /events/joined`

**Query Parameters:**
- `page` (default: 1)
- `limit` (default: 10)
- `status` (filter: pending, confirmed, declined)

**Returns:** Events where user is a member

**Possible Errors:**
- 401: Not authenticated

---

#### Endpoint: `GET /public/events`

**Returns:** Public events (no auth required)
- Only `published` events
- No user-specific data (swipes, membership)

---

### 5.3 Update Event

#### Endpoint: `PUT /events/{id}`

**Body (JSON):**
```json
{
  "title": "Updated Title",
  "status": "cancelled",
  "capacity": 10
}
```

**Rules:**
- ✅ Only creator can update
- ✅ Can update status: `published` ↔ `cancelled`
- ❌ Cannot update if status is `completed`
- ✅ Can update any field except ID/CreatorID

**Possible Errors:**
- 400: Invalid request
- 401: Not authenticated
- 403: Not authorized (not creator)
- 404: Event not found

---

### 5.4 Delete Event

#### Endpoint: `DELETE /events/{id}`

**Rules:**
- ✅ Only creator can delete
- ✅ Soft delete (sets `deleted_at`)
- ✅ Event hidden from queries
- ✅ Related data preserved

**Possible Errors:**
- 400: Invalid event ID
- 401: Not authenticated
- 403: Not authorized (not creator)
- 404: Event not found

---

## 6. Event Image/Photo Operations

### 6.1 Cover Image

#### Create/Update Cover on Event Creation
```
POST /events (multipart)
Form: file = cover_image.jpg
→ Uploads to: event_covers/{date}/{uuid}.jpg
→ Sets event.cover_image_url
```

#### Update Cover Image
```
PUT /events/{id}/cover (multipart)
Form: file = new_cover.jpg
→ Uploads to: event_covers/{date}/{uuid}.jpg
→ Updates event.cover_image_url
```

#### Serve Cover Image
```
GET /images/events/{event_id}
→ Returns event cover image
→ Requires authentication
```

---

### 6.2 Event Photos (Gallery)

#### Add Photos on Event Creation
```
POST /events (multipart)
Form: files[] = [photo1.jpg, photo2.jpg]
→ Uploads to: event_photos/{date}/{uuid}.jpg
→ Creates EventPhoto records
```

#### Add Photos to Existing Event
```
POST /events/{id}/photos (multipart)
Form: files[] = [photo1.jpg, photo2.jpg]
→ Uploads photos
→ Appends to EventPhoto table
→ Only creator can add
```

**Possible Errors:**
- 400: No files provided
- 400: Invalid file format
- 403: Not authorized (not creator)

---

## 7. Event Discovery & Suggestions

### 7.1 Get Event Suggestions

#### Endpoint: `GET /events/suggestions`

**Query Parameters:**
- `page` (default: 1)
- `limit` (default: 20)

**Algorithm:**
1. Match user tags with event tags
2. Calculate match score
3. Return events sorted by score
4. Exclude events user already swiped/joined

**Returns:**
```json
{
  "events": [
    {
      "event": {...},
      "match_score": 0.85,
      "matched_tags": [...]
    }
  ],
  "total": 50,
  "page": 1,
  "limit": 20
}
```

---

### 7.2 Discovery Flow

```
User logs in
  ↓
GET /events/suggestions
  ↓
Shows events with match scores
  ↓
User swipes (like/pass)
  ↓
POST /events/{id}/swipe
  ↓
If like → Auto join as pending member
```

---

## 8. Error Cases & Edge Cases

### 8.1 Common Error Cases

| Scenario | HTTP Status | Error Message |
|----------|-------------|---------------|
| Event not found | 404 | "Event not found" |
| Not authenticated | 401 | "User not authenticated" |
| Not authorized | 403 | "You don't have permission..." |
| Already a member | 409 | "User is already a member" |
| Not a member | 404 | "You are not a member..." |
| Event is full | 409 | "Event has reached its capacity" |
| Invalid event ID | 400 | "Invalid event ID" |

---

### 8.2 Edge Cases

#### Edge Case 1: User Joins Full Event
```
Event capacity: 5
Confirmed members: 5

POST /events/{id}/join
→ ✅ Success (creates as pending)

POST /events/{id}/confirm
→ ❌ 409 Conflict: "Event is full"
```

#### Edge Case 2: Creator Tries to Leave
```
Creator joins automatically as confirmed

POST /events/{id}/leave (by creator)
→ ✅ Success (status → left)
→ But creator still owns event
```

#### Edge Case 3: Complete Event with No Members
```
Event with only creator

POST /events/{id}/complete
→ ✅ Success
→ Status → completed
→ History created only for creator
```

#### Edge Case 4: Swipe on Own Event
```
User creates event
→ Creator auto-joined as confirmed

POST /events/{id}/swipe {direction: "like"}
→ ✅ Success (creates swipe record)
→ ❌ No new member (already exists)
```

#### Edge Case 5: Update Completed Event
```
Event status: completed

PUT /events/{id} {...}
→ ❌ Fails (cannot update completed event)
```

#### Edge Case 6: Delete Completed Event
```
Event status: completed

DELETE /events/{id}
→ ❌ Fails (cannot delete completed event)
```

---

## 9. Summary Flow Diagram

### Complete Event Lifecycle

```
1. CREATE EVENT
   POST /events
   → Status: published
   → Creator: confirmed member
   → Chat room: created

2. DISCOVERY
   GET /events/suggestions
   → Shows matched events

3. SWIPE/JOIN
   POST /events/{id}/swipe {like}
   → Member: pending
   
   OR
   
   POST /events/{id}/join
   → Member: pending

4. CONFIRM PARTICIPATION
   POST /events/{id}/confirm
   → Member: pending → confirmed

5. EVENT HAPPENS
   (Real world event)

6. COMPLETE EVENT (Creator)
   POST /events/{id}/complete
   → Status: completed
   → History: created for all confirmed members

ALTERNATIVE PATHS:
- Cancel: PUT /events/{id} {status: "cancelled"}
- Decline: POST /events/{id}/cancel
- Leave: POST /events/{id}/leave
- Delete: DELETE /events/{id}
```

---

## 10. API Endpoints Summary

### Event CRUD
- `POST /events` - Create event
- `GET /events` - List events (filtered)
- `GET /events/{id}` - Get event details
- `PUT /events/{id}` - Update event (creator only)
- `DELETE /events/{id}` - Delete event (creator only)

### Event Participation
- `POST /events/{id}/join` - Join event
- `POST /events/{id}/leave` - Leave event
- `POST /events/{id}/confirm` - Confirm participation
- `POST /events/{id}/cancel` - Cancel participation

### Event Interaction
- `POST /events/{id}/swipe` - Swipe on event (like/pass)
- `GET /events/joined` - Get user's joined events
- `GET /events/suggestions` - Get personalized suggestions

### Event Status
- `POST /events/{id}/complete` - Complete event (creator only)

### Event Images
- `PUT /events/{id}/cover` - Update cover image
- `POST /events/{id}/photos` - Add photos to gallery

### Public Access
- `GET /public/events` - Get public events
- `GET /public/events/{id}` - Get public event details

---

**Last Updated:** 2025-10-29

