#!/bin/sh

has() {
    command -v "$1" >/dev/null 2>&1
}

die() {
    echo >&2 "$@"
    exit 1
}

if ! has yue; then
    while true; do
        printf "Yuescript is not installed, install it with luarocks? [Y/n] "
        read -r yn
        case "$yn" in
        [Yy]*) break ;;
        [Nn]*) die "Install Yuescript to proceed" ;;
        *) echo "Invalid input" ;;
        esac
    done

    if ! has luarocks; then
        die "luarocks is required to install yuescript"
    fi

    luarocks install yuescript
fi

yue --target=5.1 .
