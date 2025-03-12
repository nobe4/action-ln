#!/usr/bin/env bash

set +e

for i in $(seq 0 2); do
	repo="frozen-fishsticks/action-ln-test-${i}"
	echo "Cleanig up ${repo}"

	gh api -X DELETE \
		"/repos/${repo}/git/refs/heads/test"

	gh pr list --repo "${repo}" \
		--json url \
		-q '.[].url' \
		| xargs --no-run-if-empty -n1 gh pr close --delete-branch
done
