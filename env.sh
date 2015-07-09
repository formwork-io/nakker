#!/usr/bin/env bash
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Pull in standard functions, e.g., default.
source "$DIR/.gosh.sh" || return 1
default CUSTOM_ENV_SH "$DIR/env.sh.custom"
assert-source "$CUSTOM_ENV_SH" || return 1

### GENERAL ENV VARS ###
default DIR             "$DIR"

### GREENLINE DEFAULT RAILS ###
# RAIL 0
default GL_RAIL_0_NAME          "broadcast"
default GL_RAIL_0_PROTOCOL      "broadcast"
default GL_RAIL_0_INGRESS       9002
default GL_RAIL_0_EGRESS        9003
# RAIL 1
default GL_RAIL_1_NAME          "reqrep1"
default GL_RAIL_1_PROTOCOL      "request"
default GL_RAIL_1_INGRESS       9004
default GL_RAIL_1_EGRESS        9005
# RAIL 2
default GL_RAIL_2_NAME          "reqrep2"
default GL_RAIL_2_PROTOCOL      "request"
default GL_RAIL_2_INGRESS       9006
default GL_RAIL_2_EGRESS        9007
# RAIL 3
default GL_RAIL_3_NAME          "reqrep3"
default GL_RAIL_3_PROTOCOL      "request"
default GL_RAIL_3_INGRESS       9008
default GL_RAIL_3_EGRESS        9009

### PATHS ###
default BUILD           "$DIR"/build
default TEST_OUTPUT     "$BUILD/tests"

### GOLANG ###
default GL_BUILD_ARGS   "-o $BUILD/greenline"
default GL_INSTALL_ARGS ""

### THE GO SHELL ###
default GOSH_SCRIPTS    "$DIR"/scripts
default GOSH_PROMPT     "gosh \e[0;32mgreenline\e[0m (?|#|#?)> "

