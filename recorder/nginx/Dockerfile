FROM nginx:1.27-alpine

# 自前設定ファイルを追加
# default.conf は削除しない
COPY default.conf /etc/nginx/conf.d/

# ここで default.conf に触らなければ共存する
