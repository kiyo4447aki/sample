# セットアップ手順

## 依存パッケージのインストール

```bash
sudo apt update
sudo apt install -y python3 python3-venv python3-pip libusb-1.0-0-dev \
  libjpeg-dev build-essential pkg-config
```

```bash
pip install -r requirements.txt
```

## libuvc ビルド&インストール

[libuvc リポジトリ参照](https://github.com/groupgets/libuvc)

## ユーザー作成

```bash
sudo useradd -r -s /usr/sbin/nologin thermal || true
sudo mkdir -p /var/lib/thermalwatcher/spool
sudo chown -R thermal:thermal /opt/thermalwatcher /var/lib/thermalwatcher /var/log
```

## udev ルールの作成

```bash
touch /etc/udev/rules.d/purethermal3.rules
nano /etc/udev/rules.d/purethermal3.rules
```

### purethermal3.rules

```plaintext
SUBSYSTEM=="usb", ATTRS{product}=="PureThermal 3*", MODE="0666"
```

### 適用

```bash
sudo udevadm control --reload
sudo udevadm trigger
```

## 環境変数設定

## systemd ユニット登録

```bash
sudo cp thermalwatcher.service /etc/systemd/system/
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now thermalwatcher
```
