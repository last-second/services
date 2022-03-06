import { UserStack } from "./UserStack";
import * as sst from "@serverless-stack/resources";

export default function main(app: sst.App): void {
  // Set default runtime for all functions
  app.setDefaultFunctionProps({
    runtime: "go1.x",
    environment: {
      CGO_ENABLED: "0",
      GOOS: "linux",
      GOARCH: "amd64",
    },
  });

  new UserStack(app, "UserStack");

  // Add more stacks
}
