# prmstore-exec

This is wrapper command to exec a command with ENV vars that are fetched from Amazon SSM Parameter Store.

## Usage

Saved Parameters:

```
/staging/database/host = "database.mydomain.local"
/staging/database/user = "dbuser"
/staging/database/password = "password"
```

```sh
$ prmstore-exec --path /staging --with-clean-env -- env
DATABASE_HOST=database.mydomain.local
DATABASE_USER=dbuser
DATABASE_PASSWORD=password
```

unless `--with-clean-env`, also display system ENV vars.

## Help

```
Usage:
  prmstore-exec [OPTIONS] -- command

Options:
      --path=PATH                            path for ssm:GetParametersByPath
      --no-recursive                         get parameters not recuvsively
      --no-omit-path-prefix                  No omit path prefix from parameter name
      --no-uppercase                         No convert parameter name to uppercase
      --with-clean-env                       No takeover OS Environment Variables
      --replace-map=OLD_SUBSTR:NEW_SUBSTR    Pattern Table for parameter name replacement

Help Options:
  -h, --help                                 Show this help message
```

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/joker1007/prmstore-exec.

## License

The gem is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
