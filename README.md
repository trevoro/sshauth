# SSHAuth

SSHAuth is a tool that lets you control SSH login access using the Github Team
API. It makes it easy to

* Control access to a machine using SSH
* Revoke access to a machine or group of machines
* Eliminate the need to manage `authorized_keys` files

## Getting Started

Building `sshauth` requires `go`.

    $ make setup
    $ make && sudo make install

You can optionally use `go get` to download and build sshauth. It's up to you
what system you choose to use. See the `Makefile` for more information on how to
setup the binary so that its permissions satisfy the requirements of `sshd`.

## Configuration 

You need to do **4 things** to make this work

1. Create a Github OAuth Token with a descriptive name.
2. Create a team with a descriptive name and add members to that team.
3. Copy the `config.example` to `/etc/sshauth/config.json` and fill out the
appropriate config items. This includes the `token`, `owner`, and `team`
4. Edit your `/etc/ssh/sshd_config` file and add the following two stanzas.

<!-- code block fix -->
    AuthorizedKeysCommand /usr/local/sbin/sshauth
    AuthorizedKeysCommandUser deploy
    # or root if you're feelin' gutsy

## Notes

- Keys aren't cached, so every SSH authentication request makes several API
  requests Github. If you have very large groups that are given SSH access this
  _could_ run through your API requests (you get 5000/hour)
- There is no run-time enforcement on permissions of your config file. Be
  careful.

## Background

OpenSSH version 6.6 introduced this nifty config option that makes it possible
to run a command that output a list of valid SSH public keys, similar to the
`authorized_keys` file. This means that instead of manually managing your
`authorized_keys` file on a server SSH can just run a command that does it for
you. You could collapse a directory layout so that keys are kept per file, hit a
remote endpoint or API, or something else creative.

Building on this idea is simple enough: Create a team in your Github
Organization called “ssh” or something, and then get all the SSH keys for the
users in that team.  This way, when you need to revoke access to a machine you
can just remove someone from the “ssh” group and you’re done.

The `AuthorizedKeysCommand` option in SSH is just a first pass. If the keys it
returns are not present for a user, it will continue to use the default
`authorized_keys` file. That means you could have a backup or master key on all
the servers, but individual user keys could still be fetched from Github.

If you use this in production, make sure you combine this tool with something
like [DenyHosts](http://denyhosts.sourceforge.net/ssh_config.html) so that
brute-force attempts don't hammer the Github API.

You should also make sure you have a backup key in your
`/root/.ssh/authorized_keys` file, so that if the command doesn't work, or if
the Github API is unavailable, you aren't locked out of your machine.
