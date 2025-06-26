# Микрофреймворк для диалоговых Telegram ботов

Специализированный фреймворк для создания Telegram ботов с многошаговыми диалогами, где логика описывается в YAML файлах, а бизнес-логика остается в коде.

## 🎯 Какие задачи решает

### ✅ Идеально подходит для:
- **Многошаговые опросники и анкеты** - сбор данных пользователей через последовательность вопросов
- **Системы поддержки** - создание заявок в техподдержку с категоризацией
- **FAQ боты** - интерактивные справочники с кнопочной навигацией  
- **Лид-генерация** - сбор контактов и заявок от потенциальных клиентов
- **Онбординг пользователей** - знакомство новых пользователей с сервисом
- **Заявки на услуги** - структурированный сбор информации

**YAML отвечает за:**
- Тексты сообщений и кнопок
- Последовательность шагов в диалогах
- Именование данных

**Go код отвечает за:**
- Условия запуска диалогов
- Обработку собранных данных  
- Интеграции (email, БД, API)

## 🚀 Быстрый старт

### Установка и запуск

```bash
# Клонируем проект
git clone <repo-url>
cd yamlconf

# Устанавливаем зависимости
go mod download

# Настраиваем переменные окружения
export BOT_TOKEN="your_telegram_bot_token"
export SMTP_HOST="smtp.gmail.com"
export SMTP_PORT=465
export SMTP_USERNAME="your_email@gmail.com"
export SMTP_PASSWORD="your_app_password"
export SMTP_DEFAULT_EMAIL_FROM="your_email@gmail.com"
export SMTP_DEFAULT_EMAIL_TO="admin@yourcompany.com"

# Запускаем бота
go run .
```

## 📝 Добавление простых ответов (без состояний)

Для статических сообщений и меню используйте файлы в `responses/`:

### 1. Создайте YAML файл с ответами

```yaml
# responses/products.yaml
products.catalog:
    text: "📦 Наш каталог продуктов"
    buttons:
        - text: "💎 Премиум"
          code: products.premium
        - text: "🔥 Хиты продаж"
          code: products.bestsellers
        - text: "← Назад в меню"
          code: main.menu

products.premium:
    text: |
        💎 Премиум продукты:
        • Товар 1 - 1000₽
        • Товар 2 - 1500₽
        • Товар 3 - 2000₽
    buttons:
        - text: "← Назад к каталогу"
          code: products.catalog
```

### 2. Добавьте кнопку в главное меню

```yaml
# responses/main.yaml
main.menu:
    text: Разделы
    buttons:
        - text: Каталог
          code: products.catalog  # Ссылка на новый раздел
        - text: Техподдержка
          code: support.faq.intro
```

**Готово!** Простые диалоги работают автоматически через систему кнопок.

## 🔄 Добавление многошаговых диалогов (со состояниями)

Для сложных диалогов с сохранением данных используйте flows:

### 1. Опишите flow в YAML

```yaml
# flow/config.yaml
order:
    id: order
    steps:
        idle:
            next: ask_product
        ask_product:
            message:
                text: "Какой продукт вас интересует?"
                answers:
                    - text: "💎 Премиум"
                    - text: "🔥 Базовый"
                    - text: "⭐ Эконом"
            code: product
            action: collect_text
            next: ask_quantity
        ask_quantity:
            message:
                text: "Сколько единиц нужно?"
            code: quantity
            action: collect_text
            next: ask_contact
        ask_contact:
            message:
                text: "Оставьте контакт для связи"
                answers:
                    - text: "Отправить номер"
                      request_contact: true
            code: phone
            action: collect_contact
            next: complete
        complete:
            message:
                text: "✅ Заявка принята! Мы свяжемся с вами в течение часа."
            action: send_order_request
```

### 2. Зарегистрируйте flow в коде

```go
// routes.go - добавьте в функцию registerFlows()
func (a *App) registerFlows() {
    // ... существующие flows ...
    
    // Новый flow для заказов
    orderFlow := a.flowRegistry.CreateFlow("order")
    orderFlow.InitConditionFunc = func(c tele.Context) bool {
        // Запускается по кнопке
        return c.Callback() != nil && c.Callback().Data == "order.start"
    }
}
```

### 3. Добавьте обработчик завершения

```go
// routes.go - добавьте в registerRoutes()
a.bot.Handle("send_order_request", a.onOrderComplete)

// handlers.go - добавьте новый обработчик
func (a *App) onOrderComplete(c tele.Context) error {
    if c.Get("session") == nil {
        return nil
    }
    
    session := c.Get("session").(*state.Session)
    
    // Формируем уведомление
    notification := fmt.Sprintf("Новый заказ от пользователя %d:\n", session.UserID)
    for key, value := range session.Data {
        switch v := value.(type) {
        case string:
            notification += fmt.Sprintf("%s: %s\n", key, v)
        case *tele.Contact:
            if v != nil {
                notification += fmt.Sprintf("%s: %s\n", key, v.PhoneNumber)
            }
        }
    }
    
    // Отправляем email асинхронно
    go func() {
        if err := a.mailer.Send(notification, "Новый заказ"); err != nil {
            log.Printf("Failed to send order email: %v", err)
        }
    }()
    
    return nil
}
```

### 4. Добавьте кнопку запуска

```yaml
# responses/main.yaml
main.menu:
    text: Разделы
    buttons:
        - text: "🛒 Сделать заказ"
          code: order.start  # Это активирует flow
        - text: Каталог
          code: products.catalog
```

## 🛠️ Доступные действия (Actions)

### Стандартные действия:
- `send_message` - отправить сообщение (по умолчанию)
- `collect_text` - собрать текстовый ответ  
- `collect_contact` - собрать контакт пользователя

### Кастомные действия:
Создавайте свои обработчики для специфичной логики:

```go
// Регистрируем кастомное действие
a.bot.Handle("send_to_crm", a.onSendToCRM)

func (a *App) onSendToCRM(c tele.Context) error {
    session := c.Get("session").(*state.Session)
    // Отправляем данные в CRM систему
    return a.crmClient.CreateLead(session.Data)
}
```

## 📁 Структура проекта

```
yamlconf/
├── main.go              # Точка входа
├── routes.go            # Регистрация flows и обработчиков  
├── handlers.go          # Бизнес-логика обработчиков
├── middleware.go        # Поиск и запуск flows
├── flow/               
│   ├── config.yaml      # Конфигурация многошаговых диалогов
│   ├── flow.go          # Структуры данных flows
│   ├── fsm.go           # Конечный автомат
│   └── registry.go      # Реестр flows
├── responses/           # Статические ответы и меню
│   ├── main.yaml        # Главное меню
│   ├── support.yaml     # FAQ и поддержка
│   └── errors.yaml      # Сообщения об ошибках
├── configurator/        # Загрузка YAML конфигураций
├── email/              # Отправка email уведомлений
├── state/              # Управление сессиями пользователей
└── storage/            # Хранение данных пользователей
```

## 🔧 Настройка окружения

### Обязательные переменные:
```bash
BOT_TOKEN=your_telegram_bot_token

# SMTP для email уведомлений  
SMTP_HOST=smtp.gmail.com
SMTP_PORT=465
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password
SMTP_DEFAULT_EMAIL_FROM=your_email@gmail.com
SMTP_DEFAULT_EMAIL_TO=admin@yourcompany.com
```

### Получение токена бота:
1. Найдите @BotFather в Telegram
2. Отправьте `/newbot`
3. Следуйте инструкциям
4. Скопируйте полученный токен

## 📚 Примеры использования

Посмотрите на встроенные flows для понимания возможностей:

- **`greeting`** - онбординг новых пользователей с сбором контактов
- **`support`** - система подачи заявок в техподдержку  
- **`help`** - сбор запросов на помощь по ключевым словам
