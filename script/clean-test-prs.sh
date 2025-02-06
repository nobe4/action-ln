#!/usr/bin/env bash

set -e

for i in $(seq 0 2); do
	gh pr list --repo "frozen-fishsticks/action-ln-test-${i}" \
		--json url \
		-q '.[].url' \
		| xargs --no-run-if-empty -n1 gh pr close --delete-branch
done
