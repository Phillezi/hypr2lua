# hypr2lua

An attempt of a conversion tool for migrating from `hyprlang`-based configs to the new `lua` syntax.

>[!WARNING]
>Not everything is implemented, for example 
> - `col.*` in the `general` config is emitted wrong
> - `animations` are not even implemented
> - strings get spaces removed which will break commands (if you use this, check your execs)
> - bind keys are missing " + " between them
> - some types are wrong, for example some bools get emitted as str
> 
>and much more.
>Feel free to open PRs if you want to add support!

>[!IMPORTANT]
>Since it is in a half-implemented state, you can use this tool like a starter for migrating to lua, and then do some manual fixes where the tool screws up.

## Quickstart

### Installing

#### Using `go`

>[!NOTE]
>Make sure you have the go tools dir added to your `PATH` env, to be able to use it (add this to your `~/.bashrc` and source it after)
>```bash
>export GOBIN=$(go env GOPATH)/bin
>export PATH=$GOBIN:$PATH
>```

```bash
go install github.com/phillezi/hypr2lua/cmd/hypr2lua@latest
```

### Usage

```bash
hypr2lua <hyprland config file>
# example
# hypr2lua ~/.config/hypr/hyprland.conf
#
# Output is emitted to STDOUT, you can redirect this to a file with `>`
# example
# hypr2lua ~/.config/hypr/hyprland.conf > example.lua
# NOTE: `>` will override the contents of the file you redirect it to
```
