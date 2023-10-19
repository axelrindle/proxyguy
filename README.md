# proxyguy

> Dynamic proxy generator for corporate environments.

Primarily interesting for WSL users unable to use the [forced proxy setting](https://learn.microsoft.com/de-de/windows/wsl/wsl-config#experimental-settings).

## Usage

> [!IMPORTANT]
> I only support Linux systems. There is no guarantee that this will work on Windows too.

Download the binary from the [latest release](https://github.com/axelrindle/proxyguy/releases/latest) and place it in the `PATH`, e.g. in `/usr/local/bin`.

## Bash Integration

Place the following line somewhere in your `.bashrc` file:

```shell
eval $( proxyguy -quiet )
```

The following environment variables will be automatically configured in every shell session:

- `http_proxy`
- `https_proxy`
- `no_proxy`
- `HTTP_PROXY`
- `HTTPS_PROXY`
- `NO_PROXY`

## Server Mode (WIP)

I'm working on a custom proxy server implementation which decides based on the connectivity to the corporate network
whether to forward requests to another proxy or directly to the internet.

Start the server using

```shell
proxyguy -server
```

## Todos

- [ ] Documentation
- [ ] Tests

## License

[MIT](LICENSE)
