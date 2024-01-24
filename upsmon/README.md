# UpsMon service

How to use:

  - Download binary from releases section, rename to `upsmon` and give exec permissions (`chmod +x upsmon`).
  - Copy `upsmon.service` to `/etc/systemd/system/upsmon.service`.
  - Modify `ExecStart` binary location and passed flags as needed.
  - Example script to run on ac fail timeout can be found in `acfail.sh`. Give it exec permissions (`chmod +x acfail.sh`).

UpsMon command usage:

```
  -l string
        log level, default WARN. Values: DEBUG, INFO, WARN (default "INFO")
  -p string
        post to remote url (optional)
  -r string
        rate to poll ups status, golang duration format, default 1s (default "1s")
  -s string
        run script after timeout (optional)
  -t string
        timeout to run script (optional, mandatory if -s), golang duration
  -u string
        ups json url, e.g. http://localhost:8888/0/json
```