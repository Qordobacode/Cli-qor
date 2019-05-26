# Cli-v2
Command line interface for remote work with Qordoba (https://qordoba.com/) 

# Usage

```
Usage:
  qor [COMMAND] [OPTIONS]

Application Commands:
    - init       Init configuration from stdin or from file if path to it would be provided as an argument
    
    - push       Push all local  files to the Qordoba project.  source
        Supported flags:
            -f, --files            Lists the file paths to upload (Optional)   
            -v, --version          Sets the version (Optional)
            -p, --file-path        Pushs entire folder. 
            
    - download    Download files from remote. Default file download command will give you two things  A)only the completed files B) will give you all the files (all locals and audiences
        Supported flags:
            -c, --current          Download the current state of the files 
            -a, --audiences        Specific (comma-separated) languages. example: `qor download -a en-us,de-de`
            -s, --source           Download source files
            -o, --original         Download original files (note if the customer using -s and -o in the same command rename the file original to filename-original.xxx) 
            --skip                 Skip downloading if file exists.
            
    - ls         Show present files (show 50 only). Support standard and json output format
        Standard:
        +--------+-------------------+----------+-----+-----------+---------------------+----------+
        |   ID   |       NAME        | VERSION  | TAG | #SEGMENTS |     UPDATED ON      |  STATUS  |
        +--------+-------------------+----------+-----+-----------+---------------------+----------+
        | 430005 | basic-srt.srt     | v1       |     |           | 2019-04-20 03:00:00 | DISABLED |
        | 430025 | core.csv          |          |     |           | 2019-04-22 03:00:00 | ENABLED  |
        | 430015 | core.qordoba-pot2 |          |     |           | 2019-04-22 21:13:30 | ENABLED  |
        Supported flag:
           --json                   Print output in json format
      
    - delete      Delete file from server (support file versions)  
    
    - status Status per project or file (Support file versions) 
    
    - add-key functionality
    
    - update-value functionality
   
    Help Options:
      -h, --help            Show this help message    
      -v, --version         Print CLI version
      --verbose             Switch to verbose mode
```

# Build

```bash
git clone https://github.com/Qordobacode/Cli-v2.git
cd Cli-v2
./build
bin/tf --version
```

# Install using Homebrew

Homebrew support is in progress

# Roadmap
- Update value by key command
- Create content command
- Pull value by key command
- Score per file command


  
  
  
