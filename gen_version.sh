#!/bin/bash
# This'll be replaced by `qlu version` soon...
if (git describe --abbrev=0 --exact-match &>/dev/null); then
  # We're on a tagged commit - use that as the version
  git describe --abbrev=0 --exact-match | sed 's/v\(.*\)/\1/'
else
  # Get the latest tagged version (if there is one)
  tags=$(git rev-list --tags --max-count=1 2>/dev/null)
  if [ "$tags" == "" ]; then
    v="0.0.0"
  else
    v=$(git describe --abbrev=0 --tags $tags 2>/dev/null | sed 's/v\(.*\)/\1/')
  fi
  # Split by period into an array
  a=( ${v//./ } )
  # Increment the patch-level number
  (( a[2]++ ))
  # This is a pre-release - locally it gets '-dev', in TravisCI it gets the build number
  echo "${a[0]}.${a[1]}.${a[2]}-${TRAVIS_BUILD_NUMBER:-dev}"
fi
