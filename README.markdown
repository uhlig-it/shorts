# Shorts

Poor man's [bit.ly](https://bitly.com). Given a map of urls in the config file, it will respond to a known (short) key with a redirect to the long URL.

It's pricacy-friendly. There is no tracking, no logging, and no stats.

# Deploy

The binary as present in the root directory of this project is copied to the target machine by the deployment code. Run it from `deployment` with:

```command
$ ansible-playbook playbook.yml -i foo.example.com,
```

# Development

* Build and run with

  ```command
  $ go build && ./shorts`
  ```

# TODO

* fswatch the config file and reload the URLs when changed (saves a restart when just the URLs are re-deployed)
