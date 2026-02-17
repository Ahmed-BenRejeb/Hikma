# Hikma حكمة

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/SQLite-003B57?style=for-the-badge&logo=sqlite&logoColor=white" alt="SQLite">
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
</p>

A beautiful terminal companion for **Arabic Poetry**, **Wisdom Quotes**, and **Prophetic Hadith**. Get inspired every time you open your terminal!

## ✨ Features

- 📜 **Arabic Poetry** — Classic verses from various eras (Abbasid, Andalusian, etc.)
- 💬 **Wisdom Quotes** — Inspirational quotes from great thinkers
- 🕌 **Prophetic Hadith** — Authentic sayings with Arabic narrator names
- 🎨 **Beautiful Display** — Colorful, formatted output that adapts to terminal width
- ⚡ **Fast & Lightweight** — Self-contained binary with embedded database
- 🐧 **Linux Native** — Follows XDG conventions (`~/.local/share/hikma/`)

## 📦 Installation

### APT (Debian/Ubuntu)

```bash
# 1. Add the repository
echo "deb [trusted=yes] https://apt.fury.io/ahmed-benrejeb/ /" | sudo tee /etc/apt/sources.list.d/hikma.list

# 2. Update and Install
sudo apt update
sudo apt install hikma
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/Ahmed-BenRejeb/Hikma.git
cd Hikma

# Build
go build -o hikma ./hikma.go

# Install (optional)
sudo mv hikma /usr/local/bin/
```

## 🚀 Usage

Simply run `hikma` to get a random piece of content:

```bash
hikma
```

### Options

| Flag | Description |
|------|-------------|
| `-p`, `--poems` | Show Arabic poetry |
| `-q`, `--quotes` | Show wisdom quotes |
| `-d`, `--hadith` | Show Prophetic Hadith |
| `-e`, `--era <name>` | Filter poems by era (e.g., `Abbasid`, `Andalusian`) |
| `-h`, `--help` | Show help message |

### Examples

```bash
# Random poetry
hikma -p

# Poetry from the Abbasid era
hikma -p -e Abbasid

# Random Hadith
hikma -d

# Random wisdom quote
hikma -q

# Completely random (poetry, quote, or hadith)
hikma
```

## 🎨 Sample Output

```
    ▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
      وَما مِن كاتِبٍ إِلّا سَيَفنى   ۞   وَيُبقي الدَهرُ ما كَتَبَت يَداهُ

      المتنبي | Abbasid
    ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄
```

## 💡 Pro Tip: Add to Shell Config

Add `hikma` to your `.bashrc` or `.zshrc` to get inspired every time you open a terminal:

```bash
# Add to ~/.bashrc or ~/.zshrc
hikma
```

## 🤝 Contributing

Contributions are welcome! Feel free to:

- Add more poetry, quotes, or hadith to the database
- Improve the display formatting
- Add new features or flags
- Report bugs or suggest improvements

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

## 🔗 Links

- **Repository**: [https://github.com/Ahmed-BenRejeb/Hikma](https://github.com/Ahmed-BenRejeb/Hikma)
- **Author**: Ahmed Ben Rejeb

---

<p align="center">
  <i>حكمة</i> — Arabic for "Wisdom"
</p>
