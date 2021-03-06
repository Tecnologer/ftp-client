# ftp-client

FTP Client to create backups

## Build

- Windows
  - Native: `go build -o ftpclient.exe`
  - With Makefile: `make windows`
- Linux
  - Native: `go build -o ftpclient`
  - With Makefile: `make linux`
- Darwin
  - Native: `go build -o ftpclient`
  - With Makefile: `make darwin`

For make binary for all OS, just use `make` or `make all`

## Usage

`./ftpclient[.exe] <-host <ftp-url>> [-user <username>] [-pwd [password]] [-port <port>] [-path <path>] [-dest-path <destination_path>] [-wait] [-store]`

```txt
   -host string
        (Required) URL to the server
  -path string
        location of files in the server (default "/")
  -port int
        port to connect (default 21)
  -pwd string
        password for credentials
  -store
        store flags config into settings file
  -user string
        username for credentials, default: anonymous
  -version
        returns the current version
  -wait
        prevents the program exit on finish process
  -dest-path 
        location to save the files in local, default: "."
```

## Check the version

`./ftpclient[.exe] -version`

> INFO[0000] 0.1.4.202001

## TODO

- [x] Progress bar
- [x] Settings file
- [x] Improve download process
- [x] Ignore mechanism
- [ ] Encode password in settings file
- [x] Test settings file executing from shortcut (works)
- [x] Improve fetching data
- [ ] Retry failed downloads
- [ ] GUI 
- [ ] Support SSH
- [ ] Upload files

### Dependencies

- [FTP][1]
- [Logger][2]
- [Progress Bar][4]

[1]: https://github.com/jlaffaye/ftp#goftp
[2]: https://github.com/sirupsen/logrus#logrus-
[3]: https://github.com/cheggaaa/pb#terminal-progress-bar-for-go
[4]: https://github.com/gosuri/uiprogress
