# Bookmarks Archive

Простой сервис для архивации веб-страниц с возможностью добавления заметок и получения уникального ID для бумажных блокнотов.

## Запуск

1.  Убедитесь, что у вас установлен Go.
2.  Запустите сервер:
    ```bash
    go run cmd/server/main.go
    ```
3.  Сервер будет доступен по адресу `http://localhost:8080`.

## Использование

### Веб-интерфейс
Откройте `http://localhost:8080` в браузере, чтобы просматривать сохраненные статьи и редактировать заметки.

### iOS Shortcut (Команда)

Чтобы сохранять статьи с iPhone, создайте новую команду в приложении "Команды" (Shortcuts):

1.  **Входные данные**: Установите "Receive Any" (Получать что угодно) из "Share Sheet" (Листа поделиться).
2.  **Действие 1**: "Get URLs from Shortcut Input" (Получить URL из входных данных).
3.  **Действие 2**: "Ask for Input" (Запросить входные данные) с текстом "Комментарий" (опционально).
4.  **Действие 3**: "Get Contents of URL" (Получить содержимое URL).
    *   **URL**: `http://<ВАШ_IP_АДРЕС>:8080/api/bookmarks` (Замените `<ВАШ_IP_АДРЕС>` на IP вашего компьютера в локальной сети, например `192.168.1.5`).
    *   **Method**: POST
    *   **Request Body**: JSON
    *   **Add new field**:
        *   Key: `url`, Value: (Variable from Action 1 - URL)
        *   Key: `comment`, Value: (Variable from Action 2 - Provided Input)
5.  **Действие 4**: "Get Dictionary Value" (Получить значение словаря).
    *   Get `Value` for key `id` from (Contents of URL).
6.  **Действие 5**: "Show Notification" (Показать уведомление).
    *   Text: "Saved as # (Dictionary Value)".

Теперь вы можете нажать "Поделиться" на любой странице в Safari, выбрать эту команду, ввести комментарий, и статья сохранится в вашем архиве, а вы получите ID.

## Структура проекта

*   `cmd/server`: Точка входа приложения.
*   `internal/models`: Структуры данных.
*   `internal/storage`: Работа с SQLite.
*   `internal/archiver`: Логика скачивания и очистки HTML (использует `go-readability`).
*   `web/templates`: HTML шаблоны.
