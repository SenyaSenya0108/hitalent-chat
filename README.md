# hitalent-chat

### Сборка приложения

```bash
make init
```

### Запуск приложения

```bash
make up
```

### Остановка приложения

```bash
make down
```

### Создание миграции

#### Шаблон _make create-migration name="<имя_миграции>"_

```bash
make create-migration
```

### Применение миграций

```bash
make run-migrations
```

### Перед сборкой приложения заполнить .env из .env.tmpl