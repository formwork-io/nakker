#/usr/bin/env bash

# The next three lines are for the go shell.
export SCRIPT_NAME="lint"
export SCRIPT_HELP="Lint the source."
[[ "$GOGO_GOSH_SOURCE" -eq 1 ]] && return 0

# Normal script execution starts here.
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"/../
source "$DIR"/env.sh || exit 1
cd "$DIR" || exit 1

CMD="golint"
which $CMD >/dev/null 2>&1
if [ $? -eq 1 ]; then
    echo "$CMD: command not found" >&2
    exit 1
fi
echo "[LINTING]"
for x in *.go; do
    echo -e "[\e[1m${x}\e[0m]"
    golint "$x"
    echo
done

