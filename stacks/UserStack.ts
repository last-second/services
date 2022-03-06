import * as sst from "@serverless-stack/resources";

export class UserStack extends sst.Stack {
  constructor(scope: sst.App, id: string, props?: sst.StackProps) {
    super(scope, id, props);

    // Create a DynamoDB table
    const table = new sst.Table(this, "Users", {
      primaryIndex: { partitionKey: "id" },
      fields: { id: sst.TableFieldType.STRING },
    });

    this.addDefaultFunctionEnv({
      LOGLEVEL: "debug",
      USERTABLE_NAME: table.dynamodbTable.tableName,
    });

    // Create a HTTP API
    const api = new sst.Api(this, "Api", {
      routes: {
        "POST /api/user": "lambda/create_user/main.go",
        "GET /api/user": "lambda/get_user/main.go",
        "PATCH /api/user": "lambda/update_user/main.go",
      },
    });

    api.attachPermissions([table]);

    // Show the endpoint in the output
    this.addOutputs({
      ApiEndpoint: api.url,
    });
  }
}
