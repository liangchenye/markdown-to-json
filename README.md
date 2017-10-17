# markdown-to-json
make OCI RFC error guardian easier ...

This tool is much like i18n project, 'make update-po'
will parse the codes and generate a po file.
If a new translation string was detected, the po file will updated either.

## How does it works
We provide two demo files:
- `testdata/error.go` is an old version of `specerror/error.go` file in `opencontainers/runtime-tool` project.
- `testdata/config.md` is a new version of `config.md` file in `opencontainers/runtime-spec` project.
   We assume that it is updated, so there are some RFC errors missing in `error.go`

Then we can run
```
go build
./markdown-to-json testdata/config.md > config.go
```

We got two RFC code updated: `SpecificationVersionOciVersion` and `Extensibility`.
So now, as the runtime-tool developer, we need to modify them by using a suitable variable name and remove the `TODO` lines.
Don't forget to update the real runtime validation code.
Then we can 
```
mv b ${yourpath}/runtime-tool/specerror/error.go
git add ${yourpath}/runtime-tool/specerror/error.go
git commit .. ; git push
```

Now we are done!

## TODO
- update
  I did not check the new version, so the update function is not tested.
  If we find there is update, we can add 'icon' and mark there is an update.
- percentage 
  for example, if we support 40% RFC codes, we can tell and add to icon just like UT coverage.
