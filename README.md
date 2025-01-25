# NTAZ

Утилита для анализа access.log веб-сервера nginx.

## WIP
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
Принимает либо входящий пайп (stdin) с нужным файлом, либо путь к файлу в качестве аргумента.
Например:
```bash
./ntaz mbps /logs/access.log
```
```bash
cat access.log | ./ntaz mbps
```








