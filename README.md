## Notes 
1. Something like this is also a valid route stucture `app/(site)/[slug]/suman` but it's almost bogus because it allows pages under route `/anything/suman` where anything could any any valid pathname.
1. In preview template, to write a variable inside literal braces, for example `{ {{.CamelCmp}} }` make sure there's a space between the first and the second brace. `{ {{.C` and not `{{{.C`.

## Installation

If you're on a macOS/Linux, run the folllowing command in your terminal to install the latest version of ngo.

```sh
curl -s https://ng-inky.vercel.app/install.txt | sh
```
