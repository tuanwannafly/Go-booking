# GoBooking — Kế Hoạch Triển Khai Dự Án
### Flight/Hotel Reservation Core Service (Go, Travel Tech Portfolio Project)

> Mục tiêu: xây 1 service đặt vé/phòng chống overbooking, có idempotency, có test concurrency —
> đủ chi tiết để show trong CV apply Backend Golang Intern và đủ "chuyện để kể" khi phỏng vấn.

**Timeline đề xuất:** 5 sprint x 1 tuần (giả định ~15-20h/tuần). Nếu gấp deadline nộp CV: làm xong
Sprint 0 → 2 là đã có MVP chạy được + concurrency test chứng minh được logic chống overbooking,
đủ để đưa vào CV. Sprint 3-4 hoàn thiện thêm sau.

---

## 0. Phạm Vi Dự Án

**Trong phạm vi (MVP):**
- Search chuyến bay / phòng khách sạn
- Hold ghế/phòng tạm thời (chống 2 người cùng giữ 1 chỗ)
- Tạo booking với Idempotency-Key
- Confirm / Cancel booking + refund policy
- Auto-release hold hết hạn
- Test concurrency chứng minh không overbook

**Ngoài phạm vi (không làm, để tránh lan man):**
- Payment gateway thật (chỉ mock/stub)
- Auth đầy đủ (JWT đơn giản là đủ, không cần OAuth)
- Notification email/SMS thật
- Admin dashboard UI (chỉ cần API)

---

## 1. Kiến Trúc & Tech Stack

| Thành phần | Lựa chọn | Lý do |
|---|---|---|
| Ngôn ngữ | Go 1.22+ | |
| Web framework | Gin | nhẹ, phổ biến, dễ viết middleware |
| DB chính | PostgreSQL | hỗ trợ `SELECT ... FOR UPDATE`, transaction mạnh |
| Cache/Lock | Redis + Redsync | distributed lock cho hold-seat |
| Migration | golang-migrate | version hoá schema, dễ demo trong CI |
| Test | testing + testify + `-race` | bắt buộc để chứng minh concurrency-safe |
| Container | Docker Compose | `docker compose up` chạy full stack 1 lệnh |
| CI | GitHub Actions | lint + test + build image |
| API Doc | swaggo (swagger) | tự sinh OpenAPI từ comment |
| Logging | slog (chuẩn lib Go 1.21+) | structured log, có request-id |

**Kiến trúc phân lớp (Clean Architecture rút gọn):**

```
HTTP Handler → Service (use case) → Repository → PostgreSQL/Redis
```

---

## 2. Cấu Trúc Thư Mục

```
gobooking/
├── cmd/
│   └── api/main.go
├── internal/
│   ├── config/            # load env, config struct
│   ├── domain/             # entity + interface (Flight, Seat, Booking...)
│   ├── service/            # business logic (BookingService, SearchService...)
│   ├── repository/
│   │   ├── postgres/
│   │   └── redis/
│   ├── http/
│   │   ├── handler/
│   │   ├── middleware/     # logging, recover, rate-limit, idempotency
│   │   └── router.go
│   └── worker/              # background job: release expired hold/booking
├── migrations/               # golang-migrate .sql files
├── test/
│   ├── integration/
│   └── concurrency/         # test bắn nhiều goroutine
├── docker-compose.yml
├── Dockerfile
├── .github/workflows/ci.yml
└── README.md
```

---

## 3. Database Schema

| Bảng | Cột chính | Ghi chú |
|---|---|---|
| `flights` | id, code, origin, destination, departure_time, base_price | |
| `hotels` | id, name, city, address | |
| `room_types` | id, hotel_id, name, base_price | |
| `inventory_units` | id, resource_type(`flight_seat`/`hotel_room`), resource_id, unit_code, status(`available/held/booked/cancelled`), held_until, **version int** | bảng lõi cho locking |
| `bookings` | id, user_id, status, **idempotency_key**, total_amount, expires_at, created_at | |
| `booking_items` | id, booking_id, inventory_unit_id, price | |
| `idempotency_keys` | key(PK), request_hash, response_body, status, created_at | lưu response để trả lại khi client retry |
| `payments` | id, booking_id, status, amount, provider_ref | mock payment |
| `users` | id, email, password_hash, role | auth tối giản |

**Thiết kế locking có chủ đích (điểm để kể khi phỏng vấn):**
- `flight_seat`: dùng **pessimistic lock** (`SELECT ... FOR UPDATE` + Redis distributed lock)
  vì tần suất tranh chấp cực cao trong thời gian ngắn (mở bán vé giờ vàng).
- `hotel_room`: dùng **optimistic lock** (cột `version`, update kèm `WHERE version = $x`, retry
  nếu conflict) vì tranh chấp rải đều theo khoảng ngày, không dồn burst như vé máy bay.
- → Đây chính là câu trả lời "khi nào dùng pessimistic, khi nào dùng optimistic" mà interviewer
  hay hỏi, và bạn có sẵn ví dụ thực tế thay vì trả lời lý thuyết suông.

---

## 4. Git Workflow

- Nhánh: `main` (luôn chạy được, có tag release) + `feature/<sprint>-<mô-tả-ngắn>`
- Mỗi user story = 1 branch = 1 Pull Request (tự review, mô phỏng quy trình thật để có story kể
  trong CV: "áp dụng git-flow, PR review, CI gate trước khi merge")
- Commit convention (Conventional Commits): `feat:`, `fix:`, `test:`, `docs:`, `refactor:`, `chore:`
- Tag theo sprint: `v0.1.0` (sau Sprint 1), `v0.2.0` (Sprint 2)... để CV/README có mốc rõ ràng

---

## 5. Sprint Planning

### Sprint 0 — Nền tảng (2-3 ngày)
**Mục tiêu:** repo chạy được, có DB + Redis qua Docker Compose, CI xanh.

| ID | User Story | Tasks | Branch | Điểm |
|---|---|---|---|---|
| US-01 | Là dev, tôi muốn `docker compose up` chạy được Go app + Postgres + Redis, để dev môi trường đồng nhất | Dockerfile multi-stage, docker-compose.yml, healthcheck | `feature/s0-docker-setup` | 3 |
| US-02 | Là dev, tôi muốn schema DB được version hoá, để không phải chạy SQL tay | Cài golang-migrate, viết migration cho 8 bảng ở mục 3 | `feature/s0-db-migration` | 3 |
| US-03 | Là dev, tôi muốn CI tự lint + test mỗi lần push | GitHub Actions: `golangci-lint`, `go test ./...` | `feature/s0-ci-pipeline` | 2 |
| US-04 | Là dev, tôi muốn app có config qua env + health check | `internal/config`, `GET /healthz`, `GET /readyz`, slog setup | `feature/s0-base-server` | 3 |

**DoD Sprint 0:** clone repo → `docker compose up` → `curl /healthz` trả 200.

---

### Sprint 1 — Search & Inventory (1 tuần)
**Mục tiêu:** có dữ liệu chuyến bay/phòng, search được.

| ID | User Story | Tasks | Branch | Điểm |
|---|---|---|---|---|
| US-05 | Là admin, tôi muốn seed dữ liệu flight + seat map, để có data test | Script seed, tạo `inventory_units` cho từng seat | `feature/s1-flight-seed` | 3 |
| US-06 | Là admin, tôi muốn seed hotel + room_type + room inventory | Tương tự, `resource_type=hotel_room` | `feature/s1-hotel-seed` | 3 |
| US-07 | Là khách, tôi muốn tìm chuyến bay theo origin/destination/date | `GET /flights/search?from=&to=&date=`, filter + pagination | `feature/s1-flight-search` | 5 |
| US-08 | Là khách, tôi muốn tìm phòng theo city/checkin/checkout | `GET /hotels/search`, tính số đêm, tổng giá tạm tính | `feature/s1-hotel-search` | 5 |

**AC mẫu (US-07):**
- Given có 3 chuyến bay SGN→HAN, When gọi search với date đúng, Then trả về đúng 3 kết quả kèm giá
- Trả 400 nếu thiếu param bắt buộc

---

### Sprint 2 — Hold & Locking (trọng tâm, 1 tuần)
**Mục tiêu:** chống overbooking, có test concurrency chứng minh.

| ID | User Story | Tasks | Branch | Điểm |
|---|---|---|---|---|
| US-09 | Là khách, tôi muốn giữ chỗ 1 ghế trong 10 phút trước khi thanh toán | Redsync lock + `SELECT FOR UPDATE`, set `held_until` (như snippet đã có) | `feature/s2-seat-hold-pessimistic` | 8 |
| US-10 | Là khách, tôi muốn giữ chỗ 1 phòng bằng optimistic lock | Update kèm `WHERE version=$x`, retry tối đa 3 lần nếu conflict | `feature/s2-room-hold-optimistic` | 8 |
| US-11 | Là hệ thống, tôi muốn tự động giải phóng hold hết hạn | Worker `time.Ticker` 30s quét `held_until < now`, set lại `available` | `feature/s2-auto-release-worker` | 5 |
| US-12 | Là dev, tôi muốn chứng minh không overbook khi 100 request cùng lúc | 100 goroutine gọi `HoldSeat` cùng 1 seat, assert chỉ 1 thành công, chạy với `go test -race` | `feature/s2-concurrency-test` | 5 |

**AC mẫu (US-12) — đây là phần nên chụp lại kết quả test để bỏ vào README:**
```
=== RUN   TestHoldSeat_ConcurrentRequests
    concurrency_test.go: 100 goroutines, 1 succeeded, 99 rejected (ErrSeatUnavailable)
--- PASS: TestHoldSeat_ConcurrentRequests (0.42s)
PASS
ok      gobooking/test/concurrency    0.431s
```

---

### Sprint 3 — Booking Flow & Idempotency (1 tuần)
**Mục tiêu:** hoàn thiện vòng đời booking.

| ID | User Story | Tasks | Branch | Điểm |
|---|---|---|---|---|
| US-13 | Là khách, tôi muốn tạo booking mà không bị double-charge nếu app gửi lại request | Middleware đọc header `Idempotency-Key`, check `idempotency_keys`, trả cached response nếu key đã xử lý | `feature/s3-idempotency-middleware` | 8 |
| US-14 | Là khách, tôi muốn confirm booking sau khi hold thành công | `POST /bookings/{id}/confirm`: hold→booked, tạo payment mock | `feature/s3-confirm-booking` | 5 |
| US-15 | Là khách, tôi muốn hủy booking và biết mình được hoàn bao nhiêu | `POST /bookings/{id}/cancel` + policy: hủy >24h hoàn 100%, 6-24h hoàn 50%, <6h không hoàn | `feature/s3-cancel-refund-policy` | 5 |
| US-16 | Là hệ thống, tôi muốn tự expire booking chưa confirm sau 10 phút | Mở rộng worker Sprint 2, set status `expired`, trả inventory về `available` | `feature/s3-booking-expiry` | 3 |

**AC mẫu (US-13):**
- Given đã gọi `POST /bookings` với key `abc123` thành công, When gọi lại với cùng key + cùng body,
  Then trả đúng response cũ, KHÔNG tạo booking thứ 2
- Given cùng key nhưng body khác, Then trả lỗi `422 Idempotency-Key conflict`

---

### Sprint 4 — Hardening & Polish (1 tuần)
**Mục tiêu:** sẵn sàng demo/deploy, tài liệu đầy đủ.

| ID | User Story | Tasks | Branch | Điểm |
|---|---|---|---|---|
| US-17 | Là hệ thống, tôi muốn giới hạn request theo IP để chống spam hold | Rate limit middleware (token bucket, `golang.org/x/time/rate`) | `feature/s4-rate-limiting` | 3 |
| US-18 | Là dev, tôi muốn log có request-id để trace | slog middleware, inject `X-Request-ID` | `feature/s4-structured-logging` | 3 |
| US-19 | Là dev khác, tôi muốn xem API doc mà không cần đọc code | swaggo annotation + `/swagger/index.html` | `feature/s4-openapi-doc` | 3 |
| US-20 | Là dev, tôi muốn integration test full flow (search→hold→book→confirm) | `test/integration`, dùng testcontainers hoặc docker-compose test | `feature/s4-integration-tests` | 5 |
| US-21 | Là recruiter đọc CV, tôi muốn README rõ kiến trúc & trade-off | README: sơ đồ, lý do chọn Go, so sánh pessimistic/optimistic lock, hướng dẫn chạy demo | `feature/s4-readme-docs` | 3 |

---

### Sprint 5 — Stretch Goals (tuỳ thời gian còn lại)
Chỉ làm nếu vẫn còn thời gian trước deadline nộp hồ sơ:

- Prometheus metrics (`/metrics`) + Grafana dashboard đơn giản
- JWT auth thật + role `admin`/`user`
- Load test bằng `k6` hoặc `vegeta` giả lập "flash sale", đo throughput/latency, bỏ số liệu vào README
- Circuit breaker khi gọi payment provider giả lập (link tư duy với Project 2 - Aggregator)

---

## 6. API Spec Đầy Đủ

| Method | Path | Auth | Idempotent | Mô tả |
|---|---|---|---|---|
| GET | `/healthz`, `/readyz` | - | - | health check |
| GET | `/flights/search` | - | - | tìm chuyến bay |
| GET | `/hotels/search` | - | - | tìm phòng |
| POST | `/flights/{id}/seats/{seatId}/hold` | user | - | giữ ghế tạm |
| POST | `/hotels/rooms/{roomId}/hold` | user | - | giữ phòng tạm |
| POST | `/bookings` | user | ✅ header `Idempotency-Key` | tạo booking |
| POST | `/bookings/{id}/confirm` | user | - | xác nhận + mock payment |
| POST | `/bookings/{id}/cancel` | user | - | hủy + tính refund |
| GET | `/bookings/{id}` | user | - | xem chi tiết booking |
| GET | `/swagger/index.html` | - | - | API doc |

---

## 7. Testing Strategy

| Loại | Công cụ | Mục tiêu |
|---|---|---|
| Unit test | `testing` + `testify` | service layer, business rule (refund policy...) |
| Concurrency test | goroutine + `-race` | US-12: chống overbook |
| Integration test | docker-compose test env | full flow search→hold→book→confirm→cancel |
| Load test (stretch) | k6/vegeta | đo throughput lúc nhiều hold cùng lúc |

Chạy chuẩn trước khi merge bất kỳ branch nào:
```bash
go vet ./...
golangci-lint run
go test -race -cover ./...
```

---

## 8. CI/CD Pipeline (GitHub Actions, tóm tắt)

```yaml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres: { image: postgres:16, ports: ["5432:5432"] }
      redis: { image: redis:7, ports: ["6379:6379"] }
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - run: go vet ./...
      - run: golangci-lint run
      - run: go test -race -cover ./...
      - run: docker build -t gobooking:ci .
```

---

## 9. Tài Liệu Cần Chuẩn Bị Cho CV/Phỏng Vấn

1. **README.md** — kiến trúc, cách chạy demo bằng 1 lệnh, sơ đồ thư mục
2. **ADR (Architecture Decision Record)** ngắn 1 file — ghi lại lý do chọn pessimistic vs
   optimistic lock cho 2 loại resource khác nhau (chính là phần "senior thinking" nhà tuyển dụng
   muốn thấy)
3. Screenshot/log kết quả test concurrency (US-12) — bằng chứng trực quan, dễ đưa vào CV/portfolio
4. Postman collection hoặc file `.http` mẫu cho từng endpoint

---

## 10. Câu Chuyện Để Kể Khi Phỏng Vấn

Khi được hỏi "kể về 1 project bạn tự hào", trình bày theo mạch:

1. **Vấn đề:** hệ thống đặt vé/phòng dễ bị overbook khi nhiều người cùng đặt 1 chỗ trong tích tắc
2. **Giải pháp:** kết hợp Redis distributed lock (chặn ở tầng ứng dụng, nhanh) + Postgres
   `SELECT FOR UPDATE` (chặn ở tầng dữ liệu, đảm bảo tuyệt đối) cho ghế máy bay; dùng optimistic
   locking cho phòng khách sạn vì mức độ tranh chấp thấp hơn, tránh lock DB không cần thiết
3. **Idempotency:** xử lý trường hợp client mobile mạng chập chờn gửi lại cùng 1 request đặt vé,
   tránh tạo 2 booking / trừ tiền 2 lần
4. **Bằng chứng:** test 100 goroutine cùng tranh 1 ghế với `go test -race`, chỉ 1 request thắng —
   không chỉ code chạy được mà còn *chứng minh được* nó đúng

---

## 11. Checklist Timeline Tổng

- [ ] Tuần 1: Sprint 0 + Sprint 1 (setup + search)
- [ ] Tuần 2: Sprint 2 (locking + concurrency test) ← **ưu tiên cao nhất nếu thiếu thời gian**
- [ ] Tuần 3: Sprint 3 (booking flow + idempotency)
- [ ] Tuần 4: Sprint 4 (hardening + docs) — sẵn sàng đưa vào CV
- [ ] Tuần 5 (nếu còn thời gian): Sprint 5 stretch goals
