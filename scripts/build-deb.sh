#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
VERSION=${VERSION:-1.0.0}
ARCH=${ARCH:-amd64}
PKG="liveryapi_${VERSION}_${ARCH}"
STAGE="$ROOT/dist/$PKG"

rm -rf "$STAGE"
mkdir -p "$STAGE/DEBIAN" "$STAGE/usr/bin" "$STAGE/lib/systemd/system" "$STAGE/var/lib/liveryapi"

GOOS=linux GOARCH="$ARCH" CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o "$STAGE/usr/bin/liveryapi" "$ROOT"
cp "$ROOT/debian/control" "$STAGE/DEBIAN/control"
cp "$ROOT/debian/postinst" "$STAGE/DEBIAN/postinst"
cp "$ROOT/debian/prerm" "$STAGE/DEBIAN/prerm"
cp "$ROOT/debian/liveryapi.service" "$STAGE/lib/systemd/system/liveryapi.service"

sed -i "s/^Version: .*/Version: $VERSION/" "$STAGE/DEBIAN/control"
chmod 0755 "$STAGE/usr/bin/liveryapi" "$STAGE/DEBIAN/postinst" "$STAGE/DEBIAN/prerm"
chmod 0755 "$STAGE/var/lib/liveryapi"
dpkg-deb --build "$STAGE" "$ROOT/dist/$PKG.deb"
echo "Created $ROOT/dist/$PKG.deb"
