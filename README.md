## Notes

1. Something like this is also a valid route stucture `app/(site)/[slug]/suman`
   but it's almost bogus because it allows pages under route `/anything/suman`
   where anything could any any valid pathname.
1. In preview template, to write a variable inside literal braces, for example
   `{ {{.CamelCmp}} }` make sure there's a space between the first and the
   second brace. `{ {{.C` and not `{{{.C`.
1. If you have want to install a route (say static) into a pre-existing dynamic
   route that isn't named slug (say [index]) then you can do it it as follows:
   if `/products/[pid]` already exists, and you want to install
   `/product/[pid]/archive` as a static route, then you can do it by specifying
   the route name as `/product/[slug]/archive` and it would still create the
   static route inside existing `/product/[pid]` route page. This is a
   limitation that might be fixed in the future.

## Installation

If you're on a macOS/Linux, run the folllowing command in your terminal to
install the latest version of ngo. Note that if you don't have Go installed, you
might need superuser password as part of the installation process. This is just
to add `~/.local/bin` to your path. If you would rather not provide superuser
password, you can hit control+c at the point where the password is asked to quit
the process and manually add `~/.local/bin` to your PATH environment variable.

```sh
rm -f `which ngo`
curl -s https://ng-inky.vercel.app/install.txt | sh
```
