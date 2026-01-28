# Hello World - Rust Project

一个简单的 Rust Hello World 项目。

## 项目说明

这是一个使用 Rust 编程语言创建的基础 Hello World 应用程序。

## 代码

```rust
fn main() {
    println!("Hello, world!");
}
```

## 如何运行

### 前置要求

- Rust 工具链（rustc 和 cargo）

### 运行步骤

```bash
# 克隆项目
git clone <repository-url>
cd repo

# 运行项目
cargo run

# 或者先编译再运行
cargo build
./target/debug/hello_world

# 发布版本构建
cargo build --release
./target/release/hello_world
```

## 项目结构

```
.
├── Cargo.toml      # Rust 项目配置文件
├── src/
│   └── main.rs     # 主程序入口
└── README.md       # 项目说明文档
```

## 学习资源

- [Rust 官方文档](https://www.rust-lang.org/learn)
- [Rust 程序设计语言](https://doc.rust-lang.org/book/)
- [Cargo 使用指南](https://doc.rust-lang.org/cargo/)