# wxa-cli

The `wxa-cli` utility provides various helper functions. More will be added as required.

Running `wxa-cli --version` will display the following:

```sh
Usage: wxa-cli [--version] [--help] <command> [<args>]

Available commands are:
    create-skill       List skills configured on the skills service.
    delete-skill       Delete skill on the skills service.
    generate-keys      Generate an RSA keypair in pem format.
    generate-secret    Generate a secret token for signing requests.
    list-skills        List skills configured on the skills service.
    version            Show version information.
```

You are able to get assistance for each command by providing the `--help` option for each one, e.g.:

```sh
$ wxa-cli generate-keys --help

Usage: wxa-cli [global options] generate-keys [options]

  Generate an RSA keypair in pem format.

Options:
  -public=FILENAME    Specify the filename for the generated public key. Default "public.pem".

  -private=FILENAME   Specify the filename for the generated private key. Default "private.pem".
```
