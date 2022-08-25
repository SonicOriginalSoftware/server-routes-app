# `app`

For file/web-app serving, the `app` subdomain routing is appropriate.

This will serve content, by default, from a `public` folder located in whatever the current directory context is. It can be configured by setting/`export`ing a `SERVE_PATH` `env` variable.

Requests to the `/` root path will respond automatically with the `/index.html` file.
