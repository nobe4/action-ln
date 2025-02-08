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

function jsonError(e) {
	return JSON.stringify(e, Object.getOwnPropertyNames(e));
}

module.exports = { indent, dedent, jsonError };
