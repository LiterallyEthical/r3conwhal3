<div align="center">
  <h1>r3conwhal3</h1>
</div>

<p align="center">
  <img src="assets/images/r3conwhal3.png" alt="r3conwhal3 Logo" class="img-circle" width=200 height=200>
</p>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#features">Features</a> •
  <a href="#disclaimer">Disclaimer</a> •
</p>

`r3conwhale` aims to develop a multifunctional recon chain for web applications, intelligently interpreting collected data, and optimizing performance and resource consumption through a concurrency-based approach.

# Installation

- `r3conwhal3` requires go >= 1.21.1+ to install and paths correctly set ($GOPATH, $GOROOT).

Run the following command to get the repo:

```
go install -v github.com/LiterallyEthical/r3conwhal3/cmd/r3conwhal3@latest
```

Run the following command to install dependencies

```
wget "https://raw.githubusercontent.com/LiterallyEthical/r3conwhal3/main/installer.sh"
chmod +x installer.sh
./installer.sh
```

OR

```
git clone https://github.com/LiterallyEthical/r3conwhal3
cd r3conwhal3/
chmod +x installer.sh
./installer.sh
```

<div align="center">
  
| :exclamation:  **Disclaimer**  |
|:-------------------:|
| **This project is in active development**. Expect breaking changes with releases. |

</div>

# Usage

## Options

| Flag             | Description                                                       |
| :--------------- | :---------------------------------------------------------------- |
| -A, --all        | Perform all passive & active recon process                        |
| -a, --active     | Perform active recon process (DNS bruteforce & DNS permutation)   |
| -c, --config-dir | Path to directory which config.env exists (default "embedded")    |
| -d, --domain     | Target domain to enumerate                                        |
| -o, --out-dir    | Directory to keep all output (default "$HOME/r3conwhal3/results") |
| -p, --passive    | Perform passive subdomain enumeration process                     |
| -h, --help       | Show help menu                                                    |

<div align="center">

|                                                   :exclamation: **Disclaimer**                                                    |
| :-------------------------------------------------------------------------------------------------------------------------------: |
| See the [**wiki**](https://github.com/LiterallyEthical/r3conwhal3/wiki) for running the **r3conwhal3** with custom configuration. |

</div>

## Example Usage

### Running the scan with default options

```
r3conwhal3 -d <domain-name>
```

### Running the scan with custom options

```
r3conwhal3  -d <domain> [-c <path-to-config-dir>] [-outDir <path-to-out-dir>]
```

<div align="center">

|                                                                                                                   :exclamation: **Disclaimer**                                                                                                                   |
| :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: |
| [**subkille3r**](https://github.com/LiterallyEthical/r3conwhal3/blob/main/pkg/subkill3r/subkill3r.go) is in active development. By default it's excluded from the recon chain but if you want to use it, simply supply a wordlist that contains subdomain names. |

</div>

# Features

## <div style="position: relative; display: flex; align-items: flex-end;"><img src="assets/images/inspector_gadget.ico" alt="Your Icon" width="60" height="60"> Passive Subdomain Enumeration

| ID  | Tool                                                                      | Role                                                  |
| :-: | :------------------------------------------------------------------------ | :---------------------------------------------------- |
|  1  | [subfinder](https://github.com/projectdiscovery/subfinder)                | discovering subdomains                                |
|  2  | [assetfinder](https://github.com/tomnomnom/assetfinder)                   | discovering more subdomains                           |
|  3  | [amass](https://github.com/owasp-amass/amass)                             | discovering more subdomains                           |
|  4  | [subkill3r](https://github.com/LiterallyEthical/r3conwhal3/pkg/subkill3r) | discovering more subdomains (still under development) |

## Active Subdomain Enumeration

| ID  | Tool                                           | Role                                 |
| :-: | :--------------------------------------------- | :----------------------------------- |
|  1  | [puredns](https://github.com/d3mondev/puredns) | subdomain resolving and bruteforcing |
|  2  | [gotator](https://github.com/Josue87/gotator)  | DNS permutations                     |

## Web Operations

| ID  | Tool                                                           | Role                                                |
| :-: | :------------------------------------------------------------- | :-------------------------------------------------- |
|  1  | [httpx](https://github.com/projectdiscovery/httpx/tree/v1.3.7) | filtering live domains from the gathered subdomains |

# Disclaimer

Usage of this program for attacking targets without consent is illegal. It is the user's responsibility to obey all applicable laws. The developer assumes no liability and is not responsible for any misuse or damage caused by this program. Please use responsibly.
