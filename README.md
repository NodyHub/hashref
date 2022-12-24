# hashref
Cli to create/check/update/delete all the hashes

## Output

### Usage

```shell
% hashref -h
Usage: hashref [<input> ...]

Arguments:
  [<input> ...]    Files, strings, hashes

Flags:
  -h, --help                Show context-sensitive help.
  -c, --config=STRING       Path to hashref config (default: ~/.hashref), can be overwritten
                            by environment.
  -d, --details             Show details to hash. The fiel "success" will be added to the
                            json, even if not stored on the server.
  -g, --generate            Generate client configuration
  -m, --meta=STRING         Read metadata from JSON file, comma separated file list, existing
                            keys are overwritten. Empty values are removed from metadata.
  -r, --remove              Remove hash from db
  -s, --set                 Set metadata for input. (Extends/update existing)
      --self                Set/get metadata to yourself
  -o, --output=STRING       Specify output
  -p, --publisher=STRING    Limit request to data from publisher
  -v, --verbose             Show verbose output
  -y, --yes                 Always confirm
```
