# 管理界面访问指南

## 访问地址
- 管理界面URL: http://your-domain:9147/admin/
- 需要认证才能访问

## 认证方式

### 方法1：使用浏览器插件
安装浏览器插件（如 ModHeader）添加以下请求头：
```
Authorization: Bearer YOUR_ADMIN_TOKEN
```

### 方法2：使用代理工具
可以使用 nginx 反向代理自动添加认证头：

```nginx
location /admin/ {
    proxy_pass http://localhost:9147;
    proxy_set_header Authorization "Bearer YOUR_ADMIN_TOKEN";
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

### 方法3：创建一个简单的HTML页面
创建一个 `admin-access.html` 文件：

```html
<!DOCTYPE html>
<html>
<head>
    <title>Admin Access</title>
</head>
<body>
    <h1>Shop Bot Admin Access</h1>
    <form id="adminForm">
        <label>Admin Token: <input type="password" id="token" required></label>
        <button type="submit">Access Admin Panel</button>
    </form>
    
    <script>
    document.getElementById('adminForm').onsubmit = function(e) {
        e.preventDefault();
        var token = document.getElementById('token').value;
        
        // 存储token到localStorage
        localStorage.setItem('adminToken', token);
        
        // 创建一个带认证的请求
        fetch('/admin/', {
            headers: {
                'Authorization': 'Bearer ' + token
            }
        }).then(response => {
            if (response.ok) {
                // 如果认证成功，重定向到管理界面
                window.location.href = '/admin/';
            } else {
                alert('Invalid token!');
            }
        });
    };
    </script>
</body>
</html>
```

## 获取Admin Token
Admin Token 在 `.env` 文件中配置：
```
ADMIN_TOKEN=YOUR_SECURE_TOKEN_HERE
```

## 管理界面功能
- **Dashboard**: 查看统计信息和最近订单
- **Products**: 管理商品和库存
- **Orders**: 查看和管理订单
- **Recharge Cards**: 生成和管理充值卡
- **Message Templates**: 自定义消息模板