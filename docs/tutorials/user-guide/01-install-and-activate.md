# Install and Activate Acta

Install the command and activate the per-user socket and verification timer.
This chapter uses a local source checkout because Acta does not yet define a
release distribution channel.

## Check the prerequisites

Acta requires Go, Git, and a Linux user session with systemd. Check them from
the repository root:

```sh
go version
git --version
systemctl --user --version
```

## Install the command

Install `acta` in the user-local binary directory used by the supplied systemd
units:

```sh
mkdir -p "$HOME/.local/bin"
GOBIN="$HOME/.local/bin" go install .
```

Confirm that the shell can find the command:

```sh
command -v acta
```

The result should identify `$HOME/.local/bin/acta`. Add `$HOME/.local/bin` to
`PATH` when the command is not found.

## Install the user units

Copy the supplied units from the checkout:

```sh
install -Dm0644 assets/systemd/user/acta.service \
  "$HOME/.config/systemd/user/acta.service"
install -Dm0644 assets/systemd/user/acta.socket \
  "$HOME/.config/systemd/user/acta.socket"
install -Dm0644 assets/systemd/user/acta-verify.service \
  "$HOME/.config/systemd/user/acta-verify.service"
install -Dm0644 assets/systemd/user/acta-verify.timer \
  "$HOME/.config/systemd/user/acta-verify.timer"
```

Reload the user manager and activate socket-based startup and scheduled
verification:

```sh
systemctl --user daemon-reload
systemctl --user enable --now acta.socket acta-verify.timer
```

## Check the activation

```sh
systemctl --user is-active acta.socket acta-verify.timer
```

Both units should report `active`. The service starts on the first client
request. The timer runs verification at 03:00 in the machine's local timezone.

Continue with [Connect a workspace agent](02-connect-a-workspace-agent.md).
