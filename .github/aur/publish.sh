#!/usr/bin/env bash

set -e

WD=$(cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd)
ROOT=${WD%/.github/aur}

export VERSION=$1
echo "Publishing to AUR as version ${VERSION}"

cd $WD

export GIT_SSH_COMMAND="ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"

rm -rf .pkg
git clone ssh://aur@aur.archlinux.org/typioca-git.git .pkg 2>&1

export PKGVER=$VERSION

envsubst '$PKGVER' < .SRCINFO.template > .pkg/.SRCINFO
envsubst '$PKGVER' < PKGBUILD.template > .pkg/PKGBUILD

cd .pkg

git config user.name "bloznelis"
git config user.email "bloznelis05@gmail.com"
git add -A

git commit -m "Release $VERSION"
git push origin master
