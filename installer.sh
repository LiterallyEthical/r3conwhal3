#!/bin/bash

echo -e "\033[36m      ________                            .__          .__  ________  \033[0m"
echo -e "\033[36m______\\_____  \\  ____  ____   ______  _  _|  |__ _____ |  | \\_____  \\ \033[0m"
echo -e "\033[36m\\_  __ \\_(__  <_/ ___\\/  _ \\ /    \\ \\/ \\/ /  |  \\\\__  \\ |  |   _(__  < \033[0m"
echo -e "\033[36m |  | \\/       \\  \\__(  <_> )   |  \\     /|   Y  \\/ _ \\|  |__/       \\ \033[0m"
echo -e "\033[36m |__| /______  /\\___  >____/|___|  /\\/\\_/ |___| (____  /____/______  /\033[0m" 
echo -e "\033[36m             \\/     \\/           \\/            \\/    \\/            \\/ v.1.0\033[0m"
echo -e "\033[36m                                                   ~ by LiterallyEthical\033[0m"

# Colors
RED='\033[0;31m'
CYAN='\033[0;36m'
NO_COLOR='\033[0m' # No Color

# Associative array mapping tool names to their URLs
declare -A tools=(
    ["subfinder"]="github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest"
    ["assetfinder"]="github.com/tomnomnom/assetfinder@latest"
    ["amass"]="github.com/owasp-amass/amass/v4/...@master"
    ["httpx"]="github.com/projectdiscovery/httpx/cmd/httpx@latest"
    ["puredns"]="github.com/d3mondev/puredns/v2@latest"
    ["gotator"]="github.com/Josue87/gotator@latest"
    ["gowitness"]="github.com/sensepost/gowitness@latest"
    ["ffuf"]="github.com/ffuf/ffuf/v2@latest"
    ["subzy"]="github.com/LukaSikic/subzy@latest"
  )

# Function to check if a tool is installed
is_tool_installed() {
    if which "$1" &>/dev/null; then
        echo -e "${CYAN}$1 is already installed.${NO_COLOR}"
        return 0
    else
        return 1
    fi
}

# Function to install massdns
install_massdns() {
    if is_tool_installed "massdns"; then
        return 0
    fi
    
    echo -e "Starting installation of massdns..."
    if git clone https://github.com/blechschmidt/massdns.git &>/dev/null && \
       cd massdns && \
       make &>/dev/null && \
       sudo make install &>/dev/null; then
        echo -e "${CYAN}massdns installed successfully.${NO_COLOR}"
    else
        echo -e "${RED}Failed to install massdns. Please check dependencies.${NO_COLOR}"
        return 1
    fi
}

# Main installation function that installs all tools
install_tools() {
    # Install massdns
    install_massdns || return 1

    # Install required Go tools if not already installed
    for tool in "${!tools[@]}"; do
        if is_tool_installed "$tool"; then
            continue
        fi

        echo "Installing $tool..."
        if go install -v "${tools[$tool]}"; then
            echo -e "${CYAN}Successfully installed: $tool${NO_COLOR}"
        else
            echo -e "${RED}Failed to install: $tool${NO_COLOR}"
        fi
    done
}

# Run the main installation function
install_tools
