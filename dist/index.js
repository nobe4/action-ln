const childProcess = require("child_process");
const os = require("os");
const process = require("process");

const version = "v0.0.19";

const platform = os.platform();
const arch = {
	x64: "amd64",
	arm64: "arm64",
}[os.arch()];

const binary = `main-${platform}-${arch}-${version}`;
console.log(`Binary: ${binary}`);

function main() {
	try {
		const proc = childProcess.spawnSync(`${__dirname}/${binary}`, {
			stdio: "inherit",
		});

		if (typeof proc.status === "number") {
			process.exit(proc.status);
		}

		process.exit(1);
	} catch (error) {
		console.error(error);
		process.exit(1);
	}
}

if (require.main === module) {
	main();
}
