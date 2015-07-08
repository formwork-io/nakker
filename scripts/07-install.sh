#/usr/bin/env bash

# The next three lines are for the go shell.
export SCRIPT_NAME="install"
export SCRIPT_HELP="Run go install."
[[ "$GOGO_GOSH_SOURCE" -eq 1 ]] && return 0

# Normal script execution starts here.
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"/../
source "$DIR"/env.sh || exit 1
cd "$DIR" || exit 1

CMD="go install"
which colorgo >/dev/null
if [ $? -eq 0 ]; then
    echo "[INSTALLING]"
    CMD="colorgo install"
fi
$CMD $GL_INSTALL_ARGS && echo "[COMPLETE]"

