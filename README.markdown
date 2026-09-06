# Shorts

Poor man's [bit.ly](https://bitly.com). Given a map of urls in the config file, it will respond to a known (short) key with a redirect to the long URL.

It's pricacy-friendly. There is no tracking, no logging, and no stats.

# Build

```command
$ GOOS=linux GOARCH=amd64 go build
```

# Release

Pushing a `v*` tag triggers the [release workflow](.github/workflows/release.yml), which builds `shorts` as a statically linked Linux binary for `amd64` and `arm64` and uploads the tarballs as GitHub release assets. To cut a new release:

```command
$ git tag v0.1.0
$ git push --tags
```

Then point the deployment at the new release by bumping the tag in `service_binary` (`deployment/playbook.yml`) and deploy as described below.

# Deploy

The deployment downloads the release asset referenced by `service_binary` (`deployment/playbook.yml`), installs it as a hardened systemd service, and exposes it through Caddy. Run it from `deployment` with the target host, e.g.:

```command
$ ansible-playbook playbook.yml -i neon,
```

# Development

* Build and run with

  ```command
  $ go build && ./shorts
  ```

# TODO

* Create pipeline that deploys on changes / releases
* [fswatch](https://github.com/fsnotify/fsnotify) the config file and reload the URLs when changed (saves a restart when just the URLs are re-deployed)
