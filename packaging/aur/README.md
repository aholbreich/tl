# AUR package: `tl-bin`

Arch Linux / Omarchy packaging for **tl** — a binary (`-bin`) package that
downloads prebuilt Go binaries from GitHub releases (no source build needed).

Files:

- `PKGBUILD` — the build recipe
- `.SRCINFO` — machine-readable metadata (regenerate with `makepkg --printsrcinfo > .SRCINFO`)

## Test the build locally

```sh
cd packaging/aur
makepkg -f            # build the .pkg.tar.zst without installing
namcap tl-bin-*.pkg.tar.zst   # optional: lint the package
```

## Install the built package

```sh
makepkg -si           # build + install (asks for sudo/pkexec via pacman)
# or from the archive:
sudo pacman -U tl-bin-0.9.0-1-x86_64.pkg.tar.zst
```

Verify:

```sh
pacman -Ql tl-bin     # list installed files
tl --version
tl completion --install
```

## Publishing to the AUR

The AUR uses one separate Git repository per package. Keep this project's
`packaging/aur/` directory as the packaging source, and publish copies of
`PKGBUILD` and `.SRCINFO` from a separate AUR checkout.

1. Create an account at <https://aur.archlinux.org/register> and add your SSH
   public key under **My Account**.
2. Verify SSH access:

   ```sh
   ssh aur@aur.archlinux.org help
   ```

3. Clone the new package repository (the initial clone may report that it is
   empty), then copy in the package files:

   ```sh
   git clone ssh://aur@aur.archlinux.org/tl-bin.git ~/git/tl-bin-aur
   cp packaging/aur/PKGBUILD packaging/aur/.SRCINFO ~/git/tl-bin-aur/
   cd ~/git/tl-bin-aur
   ```

4. Test, commit, and push the AUR repository:

   ```sh
   makepkg --verifysource
   makepkg -f
   namcap PKGBUILD tl-bin-*.pkg.tar.zst   # optional, requires namcap
   git add PKGBUILD .SRCINFO
   git commit -m "Initial import: tl-bin 0.9.0-1"
   git push -u origin HEAD:master
   ```

The first successful push creates <https://aur.archlinux.org/packages/tl-bin>.
After that, users can install it with `omarchy pkg aur add tl-bin`, `yay -S
tl-bin`, or another AUR helper.

## Bumping the version

1. Edit `pkgver` in `PKGBUILD`, set `pkgrel=1`.
2. Update the `sha256sums_*` (recompute: `sha256sum tl-linux-amd64.tar.gz tl-linux-arm64.tar.gz LICENSE`).
3. Regenerate `.SRCINFO` and commit.
