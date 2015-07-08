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
default GL_RAIL_0_PATTERN       "pub/sub"
default GL_RAIL_0_INGRESS_PORT  9002
default GL_RAIL_0_EGRESS_PORT   9003
# RAIL 1
default GL_RAIL_1_NAME          "reqrep1"
default GL_RAIL_1_PATTERN       "req/rep"
default GL_RAIL_1_INGRESS_PORT  9004
default GL_RAIL_1_EGRESS_PORT   9005
# RAIL 2
default GL_RAIL_2_NAME          "reqrep2"
default GL_RAIL_2_PATTERN       "req/rep"
default GL_RAIL_2_INGRESS_PORT  9006
default GL_RAIL_2_EGRESS_PORT   9007
# RAIL 3
default GL_RAIL_3_NAME          "reqrep3"
default GL_RAIL_3_PATTERN       "req/rep"
default GL_RAIL_3_INGRESS_PORT  9008
default GL_RAIL_3_EGRESS_PORT   9009

### PATHS ###
default BUILD           "$DIR"/build
default TEST_OUTPUT     "$BUILD/tests"

### GOLANG ###
default GL_BUILD_ARGS   "-o $BUILD/greenline"
default GL_INSTALL_ARGS ""

### THE GO SHELL ###
default GOSH_SCRIPTS    "$DIR"/scripts
default GOSH_PROMPT     "gosh \e[0;32mgreenline\e[0m (?|#|#?)> "

