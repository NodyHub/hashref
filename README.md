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
  -h, --help             Show context-sensitive help.
  -c, --config=STRING    Path to hashref config (default: ~/.hashref), can be overwritten by environment.
  -d, --details          Show details to hash. The fiel "success" will be added to the json, even if not stored on the
                         server.
  -g, --generate         Generate client configuration
  -m, --meta=STRING      Read metadata from JSON file, comma separated file list, existing keys are overwritten. Empty
                         values are removed from metadata.
  -r, --remove           Remove hash from db
  -s, --set              Set metadata for input. (Extends/update existing)
  -o, --output=STRING    Specify output
  -p, --publisher        Set metadata to publisher
  -v, --verbose          Show verbose output
  -y, --yes              Always confirm
```
### Example output

```shell
~% date > test
~% shasum -a 256 test
38324ddbb291ce99893bd2488d0ae740ae71b8334511f4cf141e3f953bb32d79  test
~% hashref -s -d -v test
2022/12/12 17:17:22 Flags: {Config: Details:true Generate:false Meta: Remove:false Set:true Output: Publisher:false Verbose:true Yes:false Input:[test]}
2022/12/12 17:17:22 Try to load hashref configuration json
2022/12/12 17:17:22 Loading config /Users/jan/.hashref successfull!
2022/12/12 17:17:22 Check env for configuration
2022/12/12 17:17:22 Process input test
2022/12/12 17:17:22 Input is a File!
2022/12/12 17:17:22 Calculate hash from 29 bytes
2022/12/12 17:17:22 Collect metadata for test
2022/12/12 17:17:22 Remove empty fields from metadata
2022/12/12 17:17:22 Push test to hashref
2022/12/12 17:17:22 Set data for file 38324ddbb291ce99893bd2488d0ae740ae71b8334511f4cf141e3f953bb32d79
2022/12/12 17:17:22 Convert metada to json
{
    "file": "test",
    "hash": "38324ddbb291ce99893bd2488d0ae740ae71b8334511f4cf141e3f953bb32d79",
    "input": "test",
    "last_modified": "2022-12-12 17:17:02.956557832 +0100 CET",
    "last_updated": "2022-12-12 17:17:22.546749 +0100 CET m=+0.001374459",
    "path": ".",
    "permission": "-rw-r--r--",
    "size": "29",
    "type": "file"
}
~% hashref test
test found :)
~% hashref test -d
{
    "file": "test",
    "hash": "38324ddbb291ce99893bd2488d0ae740ae71b8334511f4cf141e3f953bb32d79",
    "input": "test",
    "last_modified": "2022-12-12 17:17:02.956557832 +0100 CET",
    "last_updated": "2022-12-12 17:17:22.546749 +0100 CET m=+0.001374459",
    "path": ".",
    "permission": "-rw-r--r--",
    "publisher": "2f183a4e64493af3f377f745eda502363cd3e7ef6e4d266d444758de0a85fcc8",
    "size": "29",
    "type": "file"
}
```
