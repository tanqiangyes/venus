package cmd

import (
	"github.com/filecoin-project/venus/app/node"
	"strings"

	cmds "github.com/ipfs/go-ipfs-cmds"
)

var configCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Get and set filecoin config values",
		ShortDescription: `
venus config controls configuration variables. These variables are stored
in a config file inside your filecoin repo. When getting values, a key should be
provided, like so:

venus config KEY

When setting values, the key should be given first followed by the value and
separated by a space, like so:

venus config KEY VALUE

The key should be specified as a period separated string of keys. The value may
be either a bare string or any valid json compatible with the given key.`,
		LongDescription: `
venus config controls configuration variables. The configuration values
are stored as a JSON config file in your filecoin repo. When using venus
config, a key and value may be provided to set variables, or just a key may be
provided to fetch it's associated value without modifying it.

Keys should be listed with a dot separation for each layer of nesting within The
JSON config. For example, the "addresses" key resides within an object under the
"bootstrap" key, therefore it should be addressed with the string
"bootstrap.addresses" like so:

$ venus config bootstrap.addresses
[
	"newaddr"
]

Values may be either bare strings (be sure to quote said string if they contain
spaces to avoid arguments being separated by your shell) or as encoded JSON
compatible with the associated keys. For example, "bootstrap.addresses" expects
an array of strings, so it should be set with something like so:

$ venus config bootstrap.addresses '["newaddr"]'

When setting keys with subkeys, such as the "bootstrap" key which has 3 keys
underneath it, period, minPeerThreshold, and addresses, the given JSON value
will be merged with existing values to avoid unintentionally resetting other
configuration variables under "bootstrap". For example, setting period then
setting addresses, like so, will not change the value of "period":

$ venus config bootstrap
{
	"addresses": [],
	"minPeerThreshold": 0,
	"period": "1m"
}
$ venus config bootstrap '{"period": "5m"}'
$ venus config bootstrap '{"addresses": ["newaddr"]}'
$ venus config bootstrap
{
	"addresses": ["newaddr"],
	"minPeerThreshold": 0,
	"period": "5m"
}
`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("key", true, false, "The key of the config entry (e.g. \"api.address\")"),
		cmds.StringArg("value", false, false, "Optionally, a value with which to set the config entry"),
	},
	Run: func(req *cmds.Request, re cmds.ResponseEmitter, env cmds.Environment) error {
		api := env.(*node.Env).ConfigAPI
		key := req.Arguments[0]
		var value string

		if len(req.Arguments) == 2 {
			value = req.Arguments[1]
		} else if strings.Contains(key, "=") {
			args := strings.Split(key, "=")
			key = args[0]
			value = args[1]
		}

		if value != "" {
			err := api.ConfigSet(req.Context, key, value)
			if err != nil {
				return err
			}
		}
		res, err := api.ConfigGet(req.Context, key)
		if err != nil {
			return err
		}

		return re.Emit(res)
	},
}
