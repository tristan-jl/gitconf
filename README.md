# Gitconf

Simple Go module to manage multiple git configs.

Note: you probably want to just alias some sort of `git config set ...` command instead
of using this - I wrote primarily to start learning Go.

### Configuration

Programme expects your `.gitconfig` file to be in your `$HOME` directory. This file is
atomically copied over when the programme is used to modify it.

This programme expects the `.gitconfig` files to be stored in the
`$HOME/.config/gitconf` directory. These files should are expected to have a name such
as `config1.gitconfig`, `config2.gitconfig`. The current git profile can then be set to
config1 by:
```
gitconf set config1
```

This will atomically copy the `config1.gitconfig` over the `.gitconfig` file.

Additionally, the current git config file can be viewed using:
```
gitconf show
```

The current state is stored in the the `$HOME/.config/gitconf.config` file - this
shouldn't be manually modified.
