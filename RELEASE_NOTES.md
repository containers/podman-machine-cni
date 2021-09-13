# Release Notes

## 0.2
- Error messages from `gvproxy` are now returned to the caller, instead of being dropped.
- Fixed a bug where calling the plugin with no ports specified would cause an `unexpected end of JSON input` error.
- Fixed a bug where the plugin could return an error on a CNI `DEL` command, in violation of the CNI spec.
- The plugin will now print its version when called directly, without arguments.

## 0.1
- Initial release