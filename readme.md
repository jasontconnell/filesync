# Filesync

Usage:

```filesync -c config.json```

## Reader

The reader role will listen on the host operating system for updated files, and send those files to the writer.

## Writer

A writer will be given updated files from the reader and update the local file system with the new contents.

### Supported Operations

File creation, modification, and deletion.

Folder creation, modification, and deletion.

Renaming is not supported at this time.

## Config File Example

### Reader

```
{
    "role": "reader",
    "path": ".\\test\\source",
    "clients": [
        "127.0.0.1:57575"
    ]
}
```

### Writer

```
{
    "role": "writer",
    "path": ".\\test\\dest",
    "bind": "127.0.0.1:57575"
}
```