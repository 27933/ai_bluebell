## 1. 数据库初始化

- [x] 1.1 创建 `./sql/` 目录
- [x] 1.2 创建 `./sql/init.sql` 文件，包含所有建表语句（users, articles, tags, article_tags, categories, article_categories, comments, likes, article_visits, trending_articles）
- [x] 1.3 执行 SQL 脚本初始化数据库

## 2. 接口测试

- [x] 2.1 进入 bluebell 容器（使用 golang:1.21-alpine 容器挂载项目目录）
- [x] 2.2 运行 tests/api 下的所有测试
- [x] 2.3 验证测试结果，确保核心功能正常
