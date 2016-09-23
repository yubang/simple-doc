# 简易文档展示系统使用说明

1. golang编写，只提供Linux64位编译版本，别的平台请自行编译
2. 只需要运行二进制文件即可
3. 静态资源请直接放进static文件夹
4. 需要渲染的makedown文件请放入db文件夹
5. 路径与makedown文件映射关系， /AAA/BBB/CCC.md  ---->  /db/AAA/BBB/CCC.md
6. 配置文件请修改config/setting.conf
7. 配置文件字段说明，host为监听的域名，port为监听的端口，auth为网站访问是否需要登录，username为用户名，password为密码，auth_token为登录凭证（放在cookie），auth_cookie_key为cookie名字，auth_cookie_timeout为cookie有效期, static_timeout为静态资源缓存时间
8. 要修改渲染出来的页面，请修改模板文件config/template.html
9. 如果使用过程中遇到任何问题，欢迎反馈
10. 开发该项目的起因，由于平时要出API接口文档，所以业余时间造一个小轮子