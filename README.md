# traefik-plugin-requestid

Traefik middleware plugin that gives every request a random id: 32 hex
characters, the same format as nginx's `$request_id`.

The id is set on the request, so the backend and the access log see it, and on
the response, so a user can quote it. An `X-Request-Id` sent by the client is
replaced, unless `trustIncoming: true` is set (only behind a proxy you trust).

## Use

Traefik loads it as a local plugin. Put this repo in
`plugins-local/src/github.com/windhoos-bv/requestid` under Traefik's working
directory. The module name is the path Traefik looks in, it doesn't have to
match where the code is hosted.

Static config:

```yaml
experimental:
  localPlugins:
    requestid:
      moduleName: github.com/windhoos-bv/requestid
entryPoints:
  websecure:
    http:
      # Before crowdsec, so its block page knows the id
      middlewares: [request-id@file]
accessLog:
  fields:
    headers:
      names:
        X-Request-Id: keep  # logged as request_X-Request-Id
```

Dynamic config:

```yaml
http:
  middlewares:
    request-id:
      plugin:
        requestid:
          headerName: X-Request-Id  # default
          trustIncoming: false      # default
```

## Used by

- `windhoos/infra/containers`: the `ghcr.io/windhoos-bv/traefik` image.
- `windhoos/ansible-managed`: `roles/lib/traefik`.

## Develop

Traefik runs plugins with [Yaegi](https://github.com/traefik/yaegi), so only
the Go standard library can be used.

```sh
go test ./...
```

Release by pushing a tag `vX.Y.Z`. The users above pin a tag.
