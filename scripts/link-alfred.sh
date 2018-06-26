#!/bin/bash
_thisdir=$(realpath $(dirname "$0"))
_repodir=$(realpath "${_thisdir}/../")

echo "Linking \"$(basename ${_repodir})\" to Alfred 3 workflows"
ln -s "$_repodir" "$HOME/Library/Application Support/Alfred 3/Alfred.alfredpreferences/workflows/."

if [ $? != 0 ]
then
    echo ""                                            >> /dev/stderr
    echo "ERROR or WARNING or SOMETHING"               >> /dev/stderr
    echo "  Alfred 3 is not installed on this machine" >> /dev/stderr
    echo "  or"                                        >> /dev/stderr
    echo "  this Alfred script is already installed"   >> /dev/stderr
fi
