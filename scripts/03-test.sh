#/usr/bin/env bash

# The next three lines are for the go shell.
export SCRIPT_NAME="test"
export SCRIPT_HELP="Run tests."
[[ "$GOGO_GOSH_SOURCE" -eq 1 ]] && return 0

# Normal script execution starts here.
dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"/../
source "$dir"/env.sh || exit 1
cd "$dir" || exit 1
assert-env-or-die TEST_OUTPUT

subdirs=($(find -name "*.go" -exec dirname {} \; | sort | uniq))
for subdir in ${subdirs[@]}; do
    cd "$subdir" || exit 1
    echo -e "[\e[1mTESTING $(basename $(pwd))\e[0m]"
    go test -o "$TEST_OUTPUT" $@
    cd "$dir" || exit 1
    echo
done

