# Opener

Quickly open and close the applications you use everyday.

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

In that file, list the names of the files you want to use in an array, like [this](https://github.com/JonathanWThom/my-tab/blob/master/applications.json).

To open all the applications in your list, simply run `opener`.
When you want to close them, run `opener c`.

### Notes

This program was built on a my Mac, so probably only works on Unix systems.

If it doesn't work as intended, please open an issue!

### TODO

Better error handling/more informative messages.

### License

MIT
