<h1 align="center">
  pvm
</h1>

<p align="center">
  <strong>Focus on coding</strong><br>
  <em>pvm will help with the tooling.</em>
</p>

<p align="center">
  <a href="https://github.com/TomasBivainis/pvm/actions/workflows/test.yml">
    <img src="https://github.com/TomasBivainis/pvm/actions/workflows/test.yml/badge.svg" alt="Build Status" />
  </a>
  <a href="https://github.com/TomasBivainis/pvm/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT License" />
  </a>
</p>

## 🚀 What is pvm?

**pvm** is a simple, fast, and reliable CLI tool written in Go that abstracts away Python virtual environment and pip management.  
With pvm, you can initialize, install, and uninstall Python packages in isolated environments—no more manual venv or requirements.txt headaches.

## 🤔 Why use pvm?

- No more remembering pip/venv commands.
- Consistent, repeatable Python environments.
- Works the same on every OS.

## ✨ Features

- **Easy project initialization:** `pvm init` sets up a virtual environment and requirements.txt for you.
- **Effortless package management:** `pvm install <package>` and `pvm uninstall <package>` handle pip and requirements.txt automatically.
- **Cross-platform:** Works on Linux, macOS, and Windows.
- **No Python knowledge required:** Focus on coding, not on tooling.

## 🛠️ Installation

Download the latest release from [Releases](https://github.com/TomasBivainis/pvm/releases) or build from source:

```sh
git clone https://github.com/TomasBivainis/pvm.git
cd pvm
go build -o pvm
```

## 🚦 Usage

```sh
pvm init
pvm install requests flask
pvm uninstall flask
```

- `pvm init` — Initializes a new Python project with a virtual environment and requirements.txt.
- `pvm install <package>...` — Installs one or more pip packages and updates requirements.txt.
- `pvm uninstall <package>...` — Uninstalls packages and removes them from requirements.txt.
-

## 📄 License

MIT License. See [LICENSE](https://github.com/TomasBivainis/pvm/blob/main/LICENSE) for details.
