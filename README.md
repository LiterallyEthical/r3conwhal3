<h1 style="text-align: center;">r3conwhal3</h1>

<div style="text-align:center;">
  <img src="assets/images/r3conwhal3.png" alt="R3conKill3r Logo" style="border-radius: 50%; width: 200px; height: 200px;">
</div>

<p align="center">
  <a href="#💿-installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#features">Features</a> •
  <a href="#disclaimer">Disclaimer</a> •
</p>



`r3conwhale` aims to develop a multifunctional  recon chain for web applications, intelligently interpreting collected data, and optimizing  performance and resource consumption through a concurrency-based approach.

# 💿 Installation

- `r3conwhal3` requires go >= 1.21.1+ to install and paths correctly set ($GOPATH, $GOROOT). 

Run the following command to get the repo:

```
go install -v github.com/LiterallyEthical/r3conwhal3/cmd/r3conwhal3@latest
```

Run the following command to install dependencies

```
git clone https://github.com/LiterallyEthical/r3conwhal3
cd r3conwhal3/
chmod +x installer.sh
./installer.sh
```
OR

```
wget "https://github.com/LiterallyEthical/r3conwhal3/installer.sh"
chmod +x installer.sh
./installer.sh
```



# Usage

## Options

| Flag | Description |
|:---------|:---------|
| -domain| Target domain to enumerate |
| -dir-name | Directory to keep all output (default "$HOME/r3conwhal3/results") |
| -file-name | File to store gathered subdomains (default "subdomains.txt") |
| -help | Show help menu | 

## Example Usage

### Running the full scan with default options
```
r3conwhal3 -domain <domain-name>
```

### Running the full scan with custom options

```
r3conwhal3 <domain-name> [<-dir-name>] [<-file-name>] 
```


# Features


## <div style="position: relative; display: flex; align-items: flex-end;"><img src="assets/images/inspector_gadget.ico" alt="Your Icon" width="60" height="60"> Passive Subdomain Enumeration 

| ID | Tool | Role |
|:---------:|:---------|:---------|
| 1 | [subfinder](https://github.com/projectdiscovery/subfinder)  | discovering subdomains
| 2 | [assetfinder](https://github.com/tomnomnom/assetfinder)  | discovering more subdomains
| 3 | [amass](https://github.com/owasp-amass/amass)  | discovering more subdomains
| 4 | [subkill3r](https://github.com/LiterallyEthical/r3conwhal3/pkg/subkill3r)  | discovering more subdomains (still under development) 
| 5 | [httpx](https://github.com/projectdiscovery/httpx/tree/v1.3.7)  | filtering live domains from the gathered subdomains 

# Disclaimer

Usage of this program for attacking targets without consent is illegal. It is the user's responsibility to obey all applicable laws. The developer assumes no liability and is not responsible for any misuse or damage caused by this program. Please use responsibly.




