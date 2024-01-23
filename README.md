# Утилита для анализа access.log nginx

Утилита предназначена для анализа access.log веб-сервера nginx.

На данный момент работает только с конфигурацией log_format такого вида:

```nginx
log_format  main  '$request_time - $remote_addr $time_local $request_length - '
                  '"$request" $status $upstream_cache_status - [$body_bytes_sent] '
                  '$hostname -, $upstream_cache_status, $upstream_status';
```

## Запуск
```bash
./ntaz <command> <flag>
```

## Команды

### mbps
Вычитывает лог указанного выше формата и выводит отдаваемый трафик с разбивкой посекундно в мбит/с.

#### Поддерживамые флаги:
| Флаг          | Описание                              |
| ------------- | ------------------------------------- |
| **path**      | Путь к access.log файлу.<br> Не обязателен.<br> По умолчанию /var/log/nginx/access.log |









