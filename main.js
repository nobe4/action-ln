const core = require("@actions/core");
const github = require("@actions/github");
const yaml = require("js-yaml");
const fs = require("fs");

try {
  const payload = JSON.stringify(github.context.payload, undefined, 2);
  console.log(`The event payload: ${payload}`);

  const configPath = core.getInput("config-path", { required: true });
  const config = yaml.load(fs.readFileSync(configPath, "utf8"));
  console.log(config);
} catch (error) {
  core.setFailed(error.message);
}
