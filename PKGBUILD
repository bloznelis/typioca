# Maintainer: Energi <bloznelis05@gmail.com>
pkgname=typioca-git
name=typioca
pkgver=1.0.1
pkgrel=1
pkgdesc="Minimal, terminal based typing speed tester"
arch=(x86_64)
url="https://github.com/bloznelis/typioca"
license=(MIT)
groups=()
depends=()
makedepends=(git make go)
provides=("$name")
conflicts=("$name")
source=("git+$url")
sha256sums=('SKIP')

pkgver() {
	cd "$srcdir/$name"
  printf "%s" "$(git describe --abbrev=0 --tags)"
}

build() {
	cd "$srcdir/$name"
	make VERSION=$pkgver build
}

package() {
	cd "$srcdir/$name"

  install -Dm755 execs/$name "$pkgdir/usr/bin/$name"
  install -Dm644 LICENSE -t "$pkgdir/usr/share/licenses/$name/"
}
