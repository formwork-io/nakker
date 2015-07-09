#/usr/bin/env bash

# The next three lines are for the go shell.
export SCRIPT_NAME="build-loop"
export SCRIPT_HELP="Run go build everytime files are modified."
[[ "$GOGO_GOSH_SOURCE" -eq 1 ]] && return 0

# Normal script execution starts here.
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"/../
source "$DIR"/env.sh || exit 1
cd "$DIR" || exit 1

CMD="go build"
which colorgo >/dev/null
if [ $? -eq 0 ]; then
    CMD="colorgo build"
fi
while true; do
    "$GOSH_SCRIPTS"/gate || exit 1
    clear
    date +"%a %b %m %T"
    echo "[BUILDING]"
    $CMD $GL_BUILD_ARGS && echo "[COMPLETE]"
done
exit 0

