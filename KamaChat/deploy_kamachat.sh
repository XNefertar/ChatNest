#!/bin/bash

# 更新系统软件包
sudo apt update && sudo apt upgrade -y

# 安装 MySQL
echo "Installing MySQL..."
sudo apt install mysql-server -y

# 配置 MySQL 安全
sudo mysql_secure_installation

# 启动并启用 MySQL 服务
sudo systemctl start mysql
sudo systemctl enable mysql

# 自动创建数据库
echo "Creating database 'kama_chat_server'..."
sudo mysql -u root -p <<EOF
CREATE DATABASE kama_chat_server;
EOF

# 安装 Redis
echo "Installing Redis..."
sudo apt install redis-server -y

# 配置 Redis
# nano /etc/redis/redis.conf  # 修改 bind 127.0.0.1 改为 bind 0.0.0.0（如果需要外部访问）

# 启动并启用 Redis 服务
sudo systemctl start redis
sudo systemctl enable redis

# 卸载旧版本 Node.js 和 npm，如果不是纯净版的Ubuntu的话
echo "Uninstalling previous versions of Node.js and npm..."
sudo apt remove --purge -y nodejs npm

echo "Installing Node Version Manager (nvm)..."

apt install curl

curl -fsSL https://deb.nodesource.com/setup_16.x | sudo -E bash -

apt install -y nodejs

# 安装 Go
echo "Installing Go..."
wget https://mirrors.aliyun.com/golang/go1.20.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.20.linux-amd64.tar.gz

cp -r /usr/local/go/bin/* /usr/bin
# 设置 Go 环境变量
echo "Configuring Go environment..."

export PATH=$PATH:/usr/local/go/bin

 # 设置 Go 环境变量（避免重复写入 ~/.bashrc）
 if ! grep -q "export GOPATH=$HOME/go" ~/.bashrc; then
     echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
 fi
 source ~/.bashrc


# # 配置 Go 代理
echo "Configuring Go proxy..."
go env -w GOPROXY=https://goproxy.cn,direct


cd ~/project/KamaChat

# 安装 ssl 模块
echo "Installing ssl..."
apt-get install openssl
apt-get install libssl-dev

# # 创建根密钥，生成证书签名请求 (CSR)，创建根证书
openssl genrsa -out /etc/ssl/private/root.key 2048
openssl req -new -key /etc/ssl/private/root.key -out /etc/ssl/certs/root.csr
openssl x509 -req -in /etc/ssl/certs/root.csr -out /etc/ssl/certs/root.crt -signkey /etc/ssl/private/root.key -CAcreateserial -days 3650

# 生成服务器密钥，生成服务器证书签名请求 (CSR)，创建服务器证书扩展文件
openssl genrsa -out /etc/ssl/private/server.key 2048
openssl req -new -key /etc/ssl/private/server.key -out /etc/ssl/certs/server.csr
nano v3.ext

# # 使用根证书为服务器证书签名
openssl x509 -req -in /etc/ssl/certs/server.csr -CA /etc/ssl/certs/root.crt -CAkey /etc/ssl/private/root.key -CAcreateserial -out /etc/ssl/certs/server.crt -days 500 -sha256 -extfile v3.ext

# 安装 Apache2
apt update && apt install apache2

# 配置虚拟主机
mkdir -p /etc/apache2/sites-enabled
nano /etc/apache2/sites-enabled/000-default.conf


# 启用所需模块
a2enmod ssl
a2enmod rewrite

# 检查配置文件语法
apache2ctl configtest

# 重启 Apache 服务
systemctl start apache2
systemctl enable apache2



cd ~/project/KamaChat

# 安装 Vue.js 开发环境
echo "Installing Vue.js development environment..."
apt install npm -y

# 将 Node.js 的安装路径加入环境变量
export PATH="/root/.nvm/versions/node/v16.20.2/bin:$PATH"

# 将 Node.js 的安装路径加入 .bashrc
echo "export PATH=\$PATH:/root/.nvm/versions/node/v16.20.2/bin" >> ~/.bashrc

# 刷新 .bashrc
source ~/.bashrc

# 方案1：使用 npm 安装 Yarn
# npm install -g yarn

# 方案2：使用cnpm 安装 Yarn
npm install -g cnpm --registry=https://registry.npmjs.org
cnpm install -g yarn

# 安装 Vue CLI
cnpm install -g @vue/cli

# 重新安装项目依赖
cd ~/project/KamaChat/web/chat-server

yarn cache clean
rm -rf node_modules

yarn install # 会把package.json中所有依赖配置好的

#打包项目成dist，放到/var/www/html/，此时就可以通过云服务器的公网ip看到前端页面了
mkdir -p /var/www/html
rm -rf /var/www/html/* 
rm -rf /root/project/KamaChat/web/chat-server/dist
yarn build
cp -r /root/project/KamaChat/web/chat-server/dist/* /var/www/html # 改成自己的项目路径
chmod -R 755 /var/www/html
chown -R www-data:www-data /var/www/html


# 配置turn服务器
echo "Installing coturn..."
apt install coturn
mkdir -p /etc/coturn
nano /etc/coturn/coturn.conf

systemctl start coturn
systemctl enable coturn

cd /root/project/KamaChat/cmd/kama_chat_server
go build -o kama_chat_backend main.go

nano /etc/systemd/system/kama_chat_backend.service

systemctl daemon-reload
# 启动服务
systemctl start kama_chat_backend


# 启动服务
# systemctl start kama_chat_backend

# 启用服务（开机自启）
systemctl enable kama_chat_backend


# # 查看服务状态
# systemctl status kama_chat_backend

# # 停止服务
# systemctl stop kama_chat_backend


# 输出完成信息
echo "Deployment complete!"