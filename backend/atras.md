# atlas 使用法

## 現在スキーマの取得

```bash
atlas schema inspect -u "${DSN}" --format '{{ sql . }}' > old.sql
```

## スキーマを適用

```bash
atlas schema apply -u "${DSN}" --to file://tables.sql --dev-url "docker://postgres/15/dev"
```
