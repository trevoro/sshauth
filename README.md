# SSHAuth

SSHAuth is a tool that lets you control and manage SSH `authorized_keys` files
using the Github Teams API.

You probably don't want to use this in production, but this was a fun and (very)
quick experiment. If you seriously want to use it, make sure you combine this
tool with something like
[DenyHosts](http://denyhosts.sourceforge.net/ssh_config.html) so that
brute-force attempts don't hammer the Github API.

## Getting Started

Building `sshauth` requires `go`.

You then need to install dependencies and build:

    $ make setup
    $ make && sudo make install

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
