#!/bin/bash

if [[ ! ( $?SCRIPTS_LIB_LOGS_SH_ ) ]]; then
  return
fi

export SCRIPTS_LIB_LOGS_SH_="log.sh"
readonly SCRIPTS_LIB_LOGS_SH_

# @param_in loglevel
# @param_in message
function log() {
  datetime=[$(date +'%Y-%m-%dT%H:%M:%S%z')]
  local loglevel=$1
  shift 1

  if [[ -z "$loglevel" || ""x == "$loglevel"x ]]; then
    loglevel="INFO"
  fi
  echo "$datetime [$0] $loglevel :: $*"
}

function log::info() {
  log "INFO" "$@" >&1
}
export -f log::info

function log::debug() {
  log "DEBUG" "$@" >&1
}
export -f log::debug

function log::warn() {
  log "WARN" "$@" >&2
}
export -f log::warn

function log::error() {
  log "ERROR" "$@" >&2
}
export -f log::error