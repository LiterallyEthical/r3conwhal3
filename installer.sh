#!/bin/bash

echo -e "\033[36m      ________                            .__          .__  ________  \033[0m"
echo -e "\033[36m______\\_____  \\  ____  ____   ______  _  _|  |__ _____ |  | \\_____  \\ \033[0m"
echo -e "\033[36m\\_  __ \\_(__  <_/ ___\\/  _ \\ /    \\ \\/ \\/ /  |  \\\\__  \\ |  |   _(__  < \033[0m"
echo -e "\033[36m |  | \\/       \\  \\__(  <_> )   |  \\     /|   Y  \\/ _ \\|  |__/       \\ \033[0m"
echo -e "\033[36m |__| /______  /\\___  >____/|___|  /\\/\\_/ |___| (____  /____/______  /\033[0m" 
echo -e "\033[36m             \\/     \\/           \\/            \\/    \\/            \\/ v.1.0\033[0m"
echo -e "\033[36m                                                   ~ by LiterallyEthical\033[0m"


# Associative array mapping tool names to their URLs
declare -A tools=(
    ["subfinder"]="github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest"
    ["assetfinder"]="github.com/tomnomnom/assetfinder@latest"
    ["amass"]="github.com/owasp-amass/amass/v4/...@master"
    ["httpx"]="github.com/projectdiscovery/httpx/cmd/httpx@latest"
)

# Function to check if a Go tool is installed
is_go_tool_installed() {
    if which "$1" >/dev/null 2>&1; then
        return 0  # Tool is installed
    else
        return 1  # Tool is not installed
    fi
}

# Install required Go tools if not already installed
for tool in "${!tools[@]}"; do
    if is_go_tool_installed "$tool"; then
        echo "Already installed: $tool"
    else
        echo "Installing $tool..."
        go install -v "${tools[$tool]}"
        if [ $? -eq 0 ]; then
            echo "Successfully installed: $tool"
        else
            echo "Failed to install: $tool"
        fi
    fi
done
