#!/bin/bash
# Designed to be executed via svg-term
set -e
set -u

PROMPT="▶"

enter() {
  INPUT=$1
  DELAY=1

  prompt
  sleep "$DELAY"
  type "$INPUT"
  sleep 0.5
  printf '%b' "\\n"
  eval "$INPUT"
  type "\\n"
}

prompt() {
  printf '%b ' "$PROMPT" | pv -q
}

type() {
  printf '%b' "$1" | pv -qL $((10 + (-2 + RANDOM % 5)))
}

main() {
  IFS='%'

  enter "krew list"

  enter "cat <<EOF > krewfile
explore
modify-secret
access-matrix
krew
EOF"

  enter "krewfile --help"

  enter "krewfile -command \"krew\" -file ./krewfile"

  enter "krew list"

  prompt

  sleep 3

  echo ""

  unset IFS
}

main
