#!/bin/bash
_thisdir=$(realpath "$(dirname "$0")")
_repodir=$(realpath "${_thisdir}/../")
_target="$HOME/Library/Application Support/Alfred 3/Alfred.alfredpreferences/workflows/$(basename ${_repodir})"

echo "Linking \"$(basename "${_repodir}")\" to Alfred 3 workflows to ${_target}"

# Check if target is an actual file
if [ -f "$_target" ] || [ -d "$_target" ] && ! [ -L "$_target" ]
then
    echo "ERROR: Target exists and is an actual file or directory" >>/dev/stderr
    exit 1
fi

if ! ln -hfs "$_repodir" "$_target"
then
    echo ""                                            >> /dev/stderr
    echo "ERROR or WARNING or SOMETHING"               >> /dev/stderr
    echo "  Alfred 3 is not installed on this machine" >> /dev/stderr
    echo "  or"                                        >> /dev/stderr
    echo "  this Alfred script is already installed"   >> /dev/stderr
    exit 1
fi
