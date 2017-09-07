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
./markdown-to-json > a
diff testdata/error.go a >b
cat b
```

We will see the diff like this
```
122a123,130
> 	// TODO: newly added code, need to update.
> 	// SpecificationVersionOciVersion represents "* **`ociVersion`** (string, Required) MUST be in [SemVer v2.0.0][semver-v2.0.0] format and specifies the version of the Open Container Runtime Specification with which the bundle complies."
>         SpecificationVersionOciVersion
> 
> 	// TODO: newly added code, need to update.
> 	// Extensibility represents "Instead they MUST ignore unknown properties."
>         Extensibility
> 
178a187,191
> 	// TODO: newly added code, need to update.
> 	specificationVersionRefRef = func(version string) (reference string, err error) {
> 		return fmt.Sprintf(referenceTemplate, version, "config.md#specification-version"), nil
> 	}
> 
266a280,289
> 	// Specification version
> 
> 	// TODO: newly added code, need to update.
> 	SpecificationVersionOciVersion: {Level: rfc2119.Must, Reference: specificationVersionRef},
> 
> 	// Extensibility
> 
> 	// TODO: newly added code, need to update.
> 	Extensibility: {Level: rfc2119.Must, Reference: extensibilityRef},
> 
319a343
> 
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
