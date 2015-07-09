#/usr/bin/env bash

# The next three lines are for the go shell.
export SCRIPT_NAME="clean"
export SCRIPT_HELP="Cleans the build."
[[ "$GOGO_GOSH_SOURCE" -eq 1 ]] && return 0

# Normal script execution starts here.
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"/../
source "$DIR"/env.sh || exit 1
cd "$DIR" || exit 1
assert-env-or-die BUILD_OUTPUT
assert-env-or-die TEST_OUTPUT

rm -f "$BUILD_OUTPUT" "$TEST_OUTPUT"

