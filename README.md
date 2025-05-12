# LeCoin: One Coin to Rule Them All

LeCoin is a toy cryptocurrency modeled after Bitcoin. Its purpose is to practically
educate us on blockchains and cryptocurrencies, and so it is not intended for the 
final product to be fully functional. It is named after LeBron James, the greatest
basketball player of all time and king of the sport. We looked to his greatness to 
guide us through development.

## Running Instructions

First, make sure Go is installed, and then build the project by running `make clean all`.
To then run, simply use the command

```bash
util/vnet_run <config_dir>
```

or, to just use our pre-made configs,

```bash
util/vnet_run configs
```

The way this works is that in the configs directory you specify, you should have a set of JSON files with the following form

```json
{
	"name": "host-1",
	"port": 8000,
	"args": ["miner"]
}

```

The `args` key can be either `"miner"` or `"nonminer"` to designate that host as a machine that should do mining.
None of the hosts should have overlapping `name` or `port`.

After running `util/vnet_run <config_dir>` this will launch a tmux session where you can swap between
host panes using the macro `<Ctrl>B + O`. Once the CLI is running, you can run `help` to see the available
commands. We give an overview of these below

## LeCoin Commands
- `ss`: create a transaction sending yourself 0.0 LeCoin. This needs to be done to get the blockchain going, with the miner who
mines this first block getting the first LeCoin in circulation.
- `balance`: show your wallet (public key) plus balance
- `list`: list all known wallets and corresponding balances
- `users`: shows the index of each know wallet, usefel for the `send` command
- `send <usr_idx> <amt>`: send wallet indexed by `usr_idx` a transaction with `amt` LeCoins.
