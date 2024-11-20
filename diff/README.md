# Immutable Geth Fork Diff

You can use this file to generate an `index.html` containing a structured breakdown of all the diffs introduced by Immutable Geth from upstream.

## How To

First install the tool forkdiff:
```
go install github.com/protolambda/forkdiff
```

The tool does not seem to work with remote repositories. So you must manually clone and set up the tags you want to work with.

For example, if you fetch and tag appropriately in this repository, you can use this command:
```
forkdiff -fork ./diff/fork.yaml -repo . -upstream-repo . -out /tmp/index.html
```

Check the expected tags in the `diff/fork.yaml` file.
