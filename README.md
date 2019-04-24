# Cli-v2
Cli v4 (using Go)

### Test environment
url - https://app.qordobatest.com


### List of supported commands:
- version
- init
    Init configuration
    - from stdin
    - from file. If in command line there will be file -> it will be used as an source, that will be parsed and used internally
- push
Push all local  files to the Qordoba project.  source
Supported flags:
    - --files
    (Optional). Lists the file paths to upload.
    - --version
    (Optional). Sets the version tag. Updates the file with version tag. Uploads a file with the
    version tag
    - --file-path pushs entire (relative) file paths. Please review the commands documentation further
    down.
- download
Default file download command will give you two things  A)only the completed files B) will give you all the files (all locals and audiences
      without source file) 
      - -c --current to pull the current state of the files 
      - -a --audiences Option to work only on specific (comma-separated) languages. example: `qor pull -a en-us, de-de`
      - -s --source file option to download the update source file
      - -o original file option to download the original file (note if the customer using -s and -o in the same command rename the file original to
      filename-original.xxx) 
      --skip skip downloading if file exists.
- ls
  Show present files (show 50 only). Support standard and json output format
  - Standard:
  ```
  +--------+-------------------+----------+-----+-----------+---------------------+----------+
  |   ID   |       NAME        | VERSION  | TAG | #SEGMENTS |     UPDATED ON      |  STATUS  |
  +--------+-------------------+----------+-----+-----------+---------------------+----------+
  | 430005 | basic-srt.srt     | v1       |     |           | 2019-04-20 03:00:00 | DISABLED |
  | 430025 | core.csv          |          |     |           | 2019-04-22 03:00:00 | ENABLED  |
  | 430015 | core.qordoba-pot2 |          |     |           | 2019-04-22 21:13:30 | ENABLED  |
  ```
  - json
- delete
  Delete file from server (support file versions)  
   
  
  
  
