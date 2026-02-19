#!/bin/bash

# ============================================
# Frp 服务器端智能安装脚本（适用于 CentOS/RHEL）
# 改进版：支持重复运行，防火墙可选
# ============================================

set -e

echo "=================================="
echo "Frp 服务器端安装"
echo "=================================="

# 检查是否为 root 用户
if [ "$EUID" -ne 0 ]; then
    echo "❌ 请使用 root 权限运行此脚本"
    echo "使用: sudo bash install-server.sh"
    exit 1
fi

# 安装目录
INSTALL_DIR="/opt/frp"
FRP_VERSION="0.52.3"

# 检查是否已安装
if [ -f "$INSTALL_DIR/frps" ]; then
    echo "✅ 检测到 Frp 已安装"
    INSTALLED_VERSION=$($INSTALL_DIR/frps -v 2>&1 | grep -oP 'version \K[0-9.]+' || echo "未知")
    echo "   当前版本: $INSTALLED_VERSION"

    read -p "是否重新安装？(y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "⏭️  跳过安装"
        SKIP_INSTALL=true
    fi
fi

# 安装 Frp
if [ "$SKIP_INSTALL" != "true" ]; then
    mkdir -p $INSTALL_DIR
    cd /tmp

    # 检查是否已下载
    if [ -f "frp_${FRP_VERSION}_linux_amd64.tar.gz" ]; then
        echo "✅ 使用已下载的压缩包"
    else
        echo "📥 下载 Frp $FRP_VERSION..."
        wget -q --show-progress https://github.com/fatedier/frp/releases/download/v${FRP_VERSION}/frp_${FRP_VERSION}_linux_amd64.tar.gz
    fi

    # 解压
    echo "📦 解压文件..."
    tar -xzf frp_${FRP_VERSION}_linux_amd64.tar.gz

    # 备份旧配置
    if [ -f "$INSTALL_DIR/frps.toml" ]; then
        echo "💾 备份现有配置..."
        cp $INSTALL_DIR/frps.toml $INSTALL_DIR/frps.toml.backup.$(date +%Y%m%d_%H%M%S)
    fi

    # 安装文件
    cd frp_${FRP_VERSION}_linux_amd64
    mv frps $INSTALL_DIR/ 2>/dev/null || cp frps $INSTALL_DIR/

    # 只在配置文件不存在时才复制
    if [ ! -f "$INSTALL_DIR/frps.toml" ]; then
        mv frps.toml $INSTALL_DIR/
        echo "📝 已创建默认配置文件"
    else
        echo "📝 保留现有配置文件"
    fi

    # 清理临时文件
    cd /tmp
    rm -rf frp_${FRP_VERSION}_linux_amd64

    echo "✅ Frp 安装完成"
fi

# 创建日志目录
mkdir -p /var/log

# 创建 systemd 服务
echo "⚙️  配置系统服务..."
cat > /etc/systemd/system/frps.service << EOF
[Unit]
Description=Frp Server Service
After=network.target syslog.target
Wants=network.target

[Service]
Type=simple
Restart=on-failure
RestartSec=5s
ExecStart=/opt/frp/frps -c /opt/frp/frps.toml
LimitNOFILE=1048576

[Install]
WantedBy=multi-user.target
EOF

# 重载 systemd
systemctl daemon-reload
echo "✅ 系统服务已配置"

# 配置防火墙（可选）
echo ""
echo "🔥 检查防火墙..."
if command -v firewall-cmd &> /dev/null; then
    if systemctl is-active --quiet firewalld; then
        echo "检测到 FirewallD 正在运行"

        # 检查端口是否已开放
        if firewall-cmd --list-ports | grep -q "7000/tcp"; then
            echo "✅ 端口 7000 已开放"
        else
            echo "开放端口 7000..."
            firewall-cmd --permanent --add-port=7000/tcp
        fi

        if firewall-cmd --list-ports | grep -q "9090/tcp"; then
            echo "✅ 端口 9090 已开放"
        else
            echo "开放端口 9090..."
            firewall-cmd --permanent --add-port=9090/tcp
        fi

        if firewall-cmd --list-ports | grep -q "7500/tcp"; then
            echo "✅ 端口 7500 已开放"
        else
            echo "开放端口 7500（仪表板）..."
            firewall-cmd --permanent --add-port=7500/tcp
        fi

        firewall-cmd --reload
        echo "✅ 防火墙规则已更新"
    else
        echo "⚠️  FirewallD 已安装但未运行"
        echo "   端口已自动开放，无需额外配置"
    fi
else
    echo "⚠️  未检测到 FirewallD"
    echo "   请确保在云服务商控制台（阿里云安全组）开放以下端口："
    echo "   - 7000 (Frp 通信端口)"
    echo "   - 9090 (应用访问端口)"
    echo "   - 7500 (仪表板端口，可选)"
fi

# 检查配置文件
echo ""
echo "📝 检查配置文件..."
if [ -f "$INSTALL_DIR/frps.toml" ]; then
    # 检查是否使用默认 token
    if grep -q "change-this-to-your-secret-token-2024" "$INSTALL_DIR/frps.toml"; then
        echo "⚠️  警告：检测到使用默认认证令牌！"
        echo "   请修改 $INSTALL_DIR/frps.toml 中的 auth.token"
        NEED_CONFIG=true
    else
        echo "✅ 配置文件已存在自定义令牌"
    fi
else
    echo "❌ 配置文件不存在，请检查"
    NEED_CONFIG=true
fi

echo ""
echo "=================================="
echo "✅ 安装完成！"
echo "=================================="
echo ""

if [ "$NEED_CONFIG" = "true" ]; then
    echo "⚠️  下一步（必须）："
    echo ""
    echo "1. 编辑配置文件："
    echo "   vim $INSTALL_DIR/frps.toml"
    echo ""
    echo "2. 修改 auth.token 为复杂密码（重要！）"
    echo "   auth.token = \"your-secure-password-here\""
    echo ""
fi

echo "启动服务："
echo "   systemctl start frps"
echo ""
echo "设置开机自启："
echo "   systemctl enable frps"
echo ""
echo "查看状态："
echo "   systemctl status frps"
echo ""
echo "查看日志："
echo "   tail -f /var/log/frps.log"
echo ""
echo "测试连接："
echo "   netstat -tlnp | grep 7000"
echo ""

# 获取公网 IP
PUBLIC_IP=$(curl -s ifconfig.me 2>/dev/null || curl -s ip.sb 2>/dev/null || echo "YOUR_SERVER_IP")
echo "访问地址："
echo "   应用: http://$PUBLIC_IP:9090"
echo "   仪表板: http://$PUBLIC_IP:7500 (用户名: admin, 密码: admin123)"
echo ""
echo "⚠️  重要安全提示："
echo "   - 必须修改 auth.token（认证令牌）"
echo "   - 建议修改仪表板密码"
echo "   - 确保阿里云安全组已开放端口 7000 和 9090"
echo ""
