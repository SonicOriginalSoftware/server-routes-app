# `app`

For file/web-app serving, the `app` routing is appropriate.

This will serve content from an `fs.FS` filesystem you create and pass through `app.New()`.

Requests to the `/` root path will respond automatically with the `/index.html` file.
