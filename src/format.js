function indent(str, indent = "    ") {
	return str
		.split("\n")
		.map((l) => indent + l)
		.join("\n");
}

function dedent(str, trim = true) {
	if (!str) {
		return str;
	}

	const smallestWhitespacePrefixLen = str
		.split("\n")
		.filter((l) => l.trim())
		.map((l) => l.match(/^\s*/).slice(0)[0])
		.map((l) => l.length)
		.filter((l) => l >= 0)
		.reduce((a, b) => Math.min(a, b));

	let out = str
		.split("\n")
		.map((l) => l.substring(smallestWhitespacePrefixLen, l.length))
		.join("\n");

	if (trim) {
		out = out.trim();
	}

	return out;
}

function prettify(o) {
	try {
		return JSON.stringify(o, Object.getOwnPropertyNames(o));
	} catch {
		return `${o}`;
	}
}

function commitMessage(link) {
	return dedent(`
		${pullTitle()}
		
		From: ${link.from.toString(true)}
		To:   ${link.to.toString(true)}
	`);
}

function branchName() {
	return "auto-action-ln";
}

function pullTitle() {
	return "auto(ln): update links";
}

function pullBody(group, config, context) {
	const execution = `${context.serverUrl}/${context.repo.owner}/${context.repo.repo}/actions/runs/${context.runId}`;

	let table = [];
	for (let link of group) {
		table.push(
			`\`${link.from.toString(true)}\` | \`${link.to.toString(true)}\``,
		);
	}

	// TODO: make this not bad
	return dedent(`
		This automated PR updates the following file.
		
		From | To
		--- | ---
		${table.join("\n		")}
		
		---
		
		| Quick links | [execution](${execution}) | [configuration](${config.URL}) | [action-ln](https://github.com/nobe4/action-ln) |
		| --- | --- | --- | --- |
	`);
}

module.exports = {
	indent,
	dedent,
	prettify,
	branchName,
	commitMessage,
	pullBody,
	pullTitle,
};
