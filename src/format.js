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
	return JSON.stringify(o, Object.getOwnPropertyNames(o));
}

function branchName(link) {
	return `ln-${link.SHA256.substring(0, 8)}`;
}

function commitMessage(link) {
	return dedent(`
		${pullTitle()}
		
		From: ${link.from.toString(true)}
		To:   ${link.to.toString(true)}
	`);
}

function pullTitle() {
	return `auto(ln): update links`;
}

function pullBody(link, config) {
	return dedent(`
		This automated PR updates the following file.
		
		From | To
		--- | ---
		\`${link.from.toString(true)}\` | \`${link.to.toString(true)}\`
		
		---
		
		Configuration: [\`${config.path}\`](${config.URL}).
		Powered by [\`action-ln\`](https://github.com/nobe4/action-ln).
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
