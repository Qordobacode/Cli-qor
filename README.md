# Cli-v2
Command line interface for remote work with Qordoba (https://qordoba.com/).

This CLI may be installed via
- Using homebrew
- Using binaries from release assets
- from sources 

# Install using Homebrew

1. Install homebrew tap 
   ```
   brew tap qordobacode/qor
   ```
2. Install qor app   
   ```
   brew install qordobacode/qor/qor
   ```

# Using release binaries
From [release page](https://github.com/Qordobacode/Cli-qor/releases) it is possible to download binaries.
   

# Build from source

```bash
git clone https://github.com/Qordobacode/Cli-v2.git
cd Cli-v2
go build -o qor main.go
./qor --version
```
