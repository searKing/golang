#!/bin/bash

if [[ ! ( $?SCRIPTS_LIB_STRINGS_SH_ ) ]]; then
  return
fi

export SCRIPTS_LIB_STRINGS_SH_="strings.sh"
readonly SCRIPTS_LIB_STRINGS_SH_

# strings::trim(str)
# remove blank space in both side
strings::trim()
{
    echo "$@"
}
export -f strings::trim