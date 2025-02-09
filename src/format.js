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
	return `link-${link.SHA256}`;
}

function commitMessage(link) {
	return dedent(`
		${pullTitle()}
		
		From: ${link.from.toString(true)}
		To:   ${link.to.toString(true)}
	`);
}

function pullTitle() {
	return `auto(link): update links`;
}

function pullBody(link, config) {
	return dedent(`
		This automated PR updates the following link as defined in [the configuration](${config.URL}).
		
		## Link
		From: ${link.from.toString(true)}
		To:   ${link.to.toString(true)}
		
		---
		
		See [\`action-ln\`](https://github.com/nobe4/action-ln).
	`);
}

module.exports = {
	indent,
	dedent,
	prettify,
	branchName,
	commitMessage,
	pullBody,
};
