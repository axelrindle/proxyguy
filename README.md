# proxyguy

[![Codacy Code Quality](https://app.codacy.com/project/badge/Grade/3dec0f6cb90948418d71add803076eb8)](https://app.codacy.com/gh/axelrindle/proxyguy/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Coverage](https://app.codacy.com/project/badge/Coverage/3dec0f6cb90948418d71add803076eb8)](https://app.codacy.com/gh/axelrindle/proxyguy/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

> Dynamic proxy generator for corporate environments.

Primarily interesting for WSL users unable to use the [forced proxy setting](https://learn.microsoft.com/de-de/windows/wsl/wsl-config#experimental-settings).

## Usage

> [!IMPORTANT]
> I only support Linux/WSL systems. There is no guarantee that this will work on Windows too.

Download the binary from the [latest release](https://github.com/axelrindle/proxyguy/releases/latest) and place it in the `PATH`, e.g. in `/usr/local/bin`.

Help is available using

```shell
proxyguy --help
```

## Configuration

Configuration can be done by either using a YAML file or environment variables.

You should use the YAML file for local installations and environment variables for the Docker image.

| Key | Environment Variable | Default | Description |
| --- | --- | --- | --- |
| pac | PAC | | The URL to the .pac file. |
| timeout | TIMEOUT | `1000` | A timeout after which proxy resolving will fail. |
| proxy.override | PROXY_OVERRIDE | | Defines a static proxy endpoint. Will disable the PAC resolution. |
| proxy.ignore | PROXY_IGNORE | `localhost,127.0.0.1` | Defines the value for the `NO_PROXY` variable, urls and hosts to directly connect to. |
| proxy.determine-url | PROXY_DETERMINE | `https://ubuntu.com` | An url used to find the proxy endpoint to use. Should be a publicly available address. |
| server.address | SERVER_ADDRESS | `0.0.0.0` | On which address the server should bind. |
| server.port | SERVER_PORT | `1337` | The port to listen on. |

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
