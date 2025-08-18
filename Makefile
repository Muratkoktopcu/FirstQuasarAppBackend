# ==== Proje Ayarları ====
APP_NAME := myfirstbackend
MAIN_FILE := main.go
BIN_DIR := bin
BIN_PATH := $(BIN_DIR)/$(APP_NAME)

# Go build bayrakları (küçük, statik binary vs. istersen düzenleyebilirsin)
GO_BUILD_FLAGS := -ldflags="-s -w"

# Varsayılan hedef
.PHONY: all
all: build

# ==== Geliştirme ====
.PHONY: run
run:
	go run $(MAIN_FILE)

# Kod değişince otomatik yeniden başlatma (reflex veya air varsa)
# reflex yüklüyse: go install github.com/cespare/reflex@latest
.PHONY: dev
dev:
	@command -v reflex >/dev/null 2>&1 && reflex -r '\.go$$' -- sh -c 'go run $(MAIN_FILE)' || (echo "reflex yok. Kur: go install github.com/cespare/reflex@latest" && exit 1)

# ==== Derleme / Test ====
.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	go build $(GO_BUILD_FLAGS) -o $(BIN_PATH) $(MAIN_FILE)

.PHONY: clean
clean:
	@rm -rf $(BIN_DIR)

.PHONY: test
test:
	go test ./...

.PHONY: cover
cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

# ==== Kod Kalitesi ====
.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

# Eğer golangci-lint kuruluysa kullan (opsiyonel)
.PHONY: lint
lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || (echo "golangci-lint yok. Kur: https://golangci-lint.run/usage/install/"; exit 0)

# ==== Swagger ====
# Not: 'swag' yüklü değilse: go install github.com/swaggo/swag/cmd/swag@latest
.PHONY: swagger
swagger:
	swag init -g $(MAIN_FILE) -o ./docs

.PHONY: swagger-clean
swagger-clean:
	@rm -rf ./docs
	mkdir -p ./docs
	swag init -g $(MAIN_FILE) -o ./docs

# ==== Docker (opsiyonel: docker-compose.yml varsa) ====
.PHONY: docker-up
docker-up:
	docker compose up --build

.PHONY: docker-down
docker-down:
	docker compose down

# ==== Yardım ====
.PHONY: help
help:
	@echo ""
	@echo "Kullanışlı komutlar:"
	@echo "  make run            -> Uygulamayı çalıştır"
	@echo "  make dev            -> Kod değişince otomatik çalıştır (reflex gerekir)"
	@echo "  make build          -> Bin klasörüne derle"
	@echo "  make test           -> Testleri çalıştır"
	@echo "  make cover          -> Coverage raporu (terminalde)"
	@echo "  make fmt / make vet -> Kod format / statik analiz"
	@echo "  make lint           -> golangci-lint varsa çalıştır"
	@echo "  make swagger        -> Swagger dokümantasyonu oluştur/güncelle"
	@echo "  make swagger-clean  -> docs klasörünü temizleyip Swagger'ı yeniden üret"
	@echo "  make docker-up      -> Docker ile ayağa kaldır (compose)"
	@echo "  make docker-down    -> Docker'ı kapat"
	@echo ""
