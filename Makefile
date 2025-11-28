APP_NAME := bookmarks-server
DIST_DIR := dist

.PHONY: all build clean release-arm64

all: build

# Обычная сборка для текущей ОС
build:
	go build -o $(APP_NAME) cmd/server/main.go

# Очистка артефактов сборки
clean:
	rm -rf $(DIST_DIR) $(APP_NAME)

# Кросс-компиляция для Linux ARM64 и создание архива
release-arm64: clean
	mkdir -p $(DIST_DIR)/$(APP_NAME)
	# Компиляция для Linux ARM64 (CGO_ENABLED=0 для pure Go sqlite)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(DIST_DIR)/$(APP_NAME)/$(APP_NAME) cmd/server/main.go
	# Копирование шаблонов
	cp -r web $(DIST_DIR)/$(APP_NAME)/
	# Создание архива
	cd $(DIST_DIR) && tar -czvf $(APP_NAME)-linux-arm64.tar.gz $(APP_NAME)
	@echo "Архив готов: $(DIST_DIR)/$(APP_NAME)-linux-arm64.tar.gz"
