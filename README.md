# Opener

Quickly open and close the applications you use every day.

### Installation

If you have Go installed:

```
go get -u github.com/jonathanwthom/opener
cd $GOPATH/src/github.com/jonathanwthom/opener
go install
```

Otherwise, move the binary to somewhere in your PATH:

```
git clone https://github.com/JonathanWThom/opener
cd opener
mv opener /usr/local/bin/opener
```

### Usage

In the root of your filesystem, create a file called applications.json.

`touch ~/applications.json`

In that file, list the names of the files you want to use like [this](https://github.com/JonathanWThom/opener/blob/master/applications.json).

You will need to have one key that is "default" - these are the applications that will open and close with no additional arguments passed.

To open all the applications in your list, simply run `opener`.
When you want to close them, run `opener -c`.

To create extra groups, simply add another key with a corresponding list to `applications.json`. You can then open and close that group
of apps by passing the `-g` flag. For example: `opener -g weekend` and `opener -c -g weekend`.

##### Browsers

When including web browsers in your list of applications, you can specify sites to open by default. Add another group with the name of
the browser, and an array of sites you want to open, like this:

```
{
    "default": [
        "Google Chrome"
    ],
    "Google Chrome": [
        "https://stackoverflow.com",
        "https://twitter.com"
    ]
}
```

### Notes

This program was built on a my Mac, so probably only works on Unix systems.

If it doesn't work as intended, please open an issue!

### License

MIT
