go-runtime
==========

The `go-runtime` project is a proof of concept, which aims to verify whether it is possible to delegate handling of http requests to JS code. This would allow for a go-based runtime for Dirigible.

**Getting the project**

```bash
go get -v github.com/dirigiblelabs/go-runtime
```

**Running the code**

```bash
go run main.go
```

By default, this will start an http server on `127.0.0.1:8080`. You can use the `HOST` and `PORT` environment variables to change that.

You can test that it works by calling the example javascript endpoint.

```
curl http://127.0.0.1:8080/example.js
```

If you want to run your own javascript code, just place it somewhere in the directory of the project, or just edit the `example.js` file.

**Building a standalone executable**

You can build and use the executable in any arbitrary directory.

```bash
go buld -o go-runtime main.go
```

```bash
cp go-runtime ~/my-backend
cd ~/my-backend
./go-runtime
```

**Remarks**

* This is a PoC project, so don't use it for productive use.
* Only files that end in `.js` can be called.
* Only files that are in the directory or subdirectories of where the executable is run can be called.
* The project does not come with vendored dependencies, so if it remains unmaintained for a long period of time, it might fail to compile. Regardless, the idea here is to preserve the design concept.
