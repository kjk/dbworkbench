# pgweb

Web-based PostgreSQL database browser written in Go.

Based on http://github.com/sosedoff/pgweb

## Usage

Start server:

```
pgweb
```

You can also provide connection flags:

```
pgweb --host localhost --user myuser --db mydb
```

Connection URL scheme is also supported:

```
pgweb --url postgres://user:password@host:port/database?sslmode=[mode]
```

## Testing

Run tests:

```
make test
```

## License

The MIT License (MIT)

Copyright (c) 2014-2015 Dan Sosedoff, <dan.sosedoff@gmail.com>
