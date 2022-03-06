import { SharedStack } from "./SharedStack";
import { UserStack } from "./UserStack";
import { TaskStack } from "./TaskStack";
import * as sst from "@serverless-stack/resources";
import path from "path";

const stripHandler = (handler: string) =>
  path.basename(handler.replace(/\/main\.go$/, ""));

export default function main(app: sst.App): void {
  // Set default runtime for all functions
  app.setDefaultFunctionProps({
    functionName: ({ stack: { stackName }, functionProps: { handler } }) =>
      `${stackName}-${stripHandler(handler ?? "")}`,
    runtime: "go1.x",
    environment: {
      CGO_ENABLED: "0",
      GOOS: "linux",
      GOARCH: "amd64",
    },
  });

  const sharedStack = new SharedStack(app, "SharedStack");

  const userStack = new UserStack(app, "UserStack", { api: sharedStack.api });
  userStack.addDependency(sharedStack);

  const taskStack = new TaskStack(app, "TaskStack", {
    api: sharedStack.api,
    userTable: userStack.table,
  });
  taskStack.addDependency(sharedStack);
  taskStack.addDependency(userStack);

  // Add more stacks
}
