// Code generated by qtc from "login.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

package templates

import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

func StreamLogin(qw422016 *qt422016.Writer, clientID string, clientEndpoint string) {
	qw422016.N().S(`
<!doctype html>
<html lang="en">

<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@next/css/pico.min.css" />
  <title>Hello, world!</title>
</head>

<body>
  <main class="container">
    <a href="/auth/github/login" role="button">Login with Github</a>
  </main>
  </script>
</body>

</html>
`)
}

func WriteLogin(qq422016 qtio422016.Writer, clientID string, clientEndpoint string) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	StreamLogin(qw422016, clientID, clientEndpoint)
	qt422016.ReleaseWriter(qw422016)
}

func Login(clientID string, clientEndpoint string) string {
	qb422016 := qt422016.AcquireByteBuffer()
	WriteLogin(qb422016, clientID, clientEndpoint)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}
