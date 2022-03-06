import * as sst from "@serverless-stack/resources";

export class SharedStack extends sst.Stack {
  api: sst.Api;
  constructor(scope: sst.App, id: string, props?: sst.StackProps) {
    super(scope, id, props);

    // Create a HTTP API
    this.api = new sst.Api(this, "Api");
  }
}
