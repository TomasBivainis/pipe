<h1 align="center">
  pvm
</h1>

<p align="center">
  <strong>Focus on coding.</strong><br>
  <em>Let pvm handle Python environments and pip.</em>
</p>

<p align="center">
  <a href="https://github.com/TomasBivainis/pvm/actions/workflows/test.yml">
    <img src="https://github.com/TomasBivainis/pvm/actions/workflows/test.yml/badge.svg" alt="Build Status" />
  </a>
  <a href="https://img.shields.io/github/v/release/TomasBivainis/pvm">
    <img src="https://img.shields.io/github/v/release/TomasBivainis/pvm" alt="Latest Release" />
  </a>
  <a href="https://img.shields.io/github/go-mod/go-version/TomasBivainis/pvm">
    <img src="https://img.shields.io/github/go-mod/go-version/TomasBivainis/pvm" alt="Go Version" />
  </a>
  <a href="https://github.com/TomasBivainis/pvm/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="MIT License" />
  </a>
</p>

---

## 🚀 What is `pvm`?

**pvm** is a fast and reliable CLI tool written in Go that abstracts away Python virtual environment and package management.

Say goodbye to fiddling with `venv`, `pip`, and `requirements.txt`. Just run `pvm` and get back to coding.

---

## 🤔 Why use `pvm`?

- 🧠 No more remembering `pip` and `venv` syntax.
- 🔁 Consistent, repeatable Python environments across systems.
- 🖥️ Cross-platform: macOS, Linux, Windows.
- 👶 No Python expertise needed—ideal for new developers.

---

## ✨ Features

- ✅ `pvm init` — Create a virtual environment and a `requirements.txt`.
- 📦 `pvm install <package>` — Install pip packages _and_ update `requirements.txt`.
- ❌ `pvm uninstall <package>` — Clean removal of packages and their entries.
- 🔄 Reproducible environments without external tools.

---

## 🛠️ Installation

### Option 1: Download a Release

1. Visit the [Releases](https://github.com/TomasBivainis/pvm/releases) page.
2. Download the appropriate binary for your OS and extract it:

   #### On Linux/macOS:

   ```sh
   tar -xzf pvm-linux-amd64.tar.gz
   ```

3. Move it somewhere in your `PATH`:

   #### On Linux/macOS:

   ```sh
   sudo mv pvm /usr/local/bin/
   ```

   #### On Windows (PowerShell):

   - Move `pvm.exe` to a directory, e.g. `C:\Tools\pvm\`
   - Add that folder to your System PATH:

     1. Open **System Properties** → **Environment Variables**
     2. Edit **Path** → Add: `C:\Tools\pvm\`

   - Restart your terminal or run:

     ```powershell
     refreshenv
     ```

---

### Option 2: Build from Source

```sh
git clone https://github.com/TomasBivainis/pvm.git
cd pvm
go build -o pvm
```

Then add the `pvm` binary to your `PATH` as described above.

---

## 🚦 Usage

```sh
pvm init
pvm install requests flask
pvm uninstall flask
```

- `pvm init` — Initializes a Python project with a virtual environment and `requirements.txt`.
- `pvm install <package>...` — Installs one or more pip packages and updates `requirements.txt`.
- `pvm uninstall <package>...` — Uninstalls packages and removes them from `requirements.txt`.

---

## 📄 License

MIT License. See [LICENSE](https://github.com/TomasBivainis/pvm/blob/main/LICENSE) for details.

---

## 🙌 Contributing

Contributions are welcome! Feel free to open issues, suggest features, or submit a pull request.
