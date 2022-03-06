import * as sst from "@serverless-stack/resources";

export interface UserStackProps extends sst.StackProps {
  api: sst.Api;
}

export class UserStack extends sst.Stack {
  table: sst.Table;

  constructor(scope: sst.App, id: string, { api, ...props }: UserStackProps) {
    super(scope, id, props);

    // Create a DynamoDB table
    const table = new sst.Table(this, "Users", {
      primaryIndex: { partitionKey: "id" },
      fields: { id: sst.TableFieldType.STRING },
    });

    this.addDefaultFunctionEnv({
      LOGLEVEL: "debug",
      USERTABLE_NAME: table.tableName,
    });

    // Register lambdas
    const routes: sst.ApiProps["routes"] = {
      "POST /api/user": "lambda/create_user/main.go",
      "GET /api/user": "lambda/get_user/main.go",
      "PATCH /api/user": "lambda/update_user/main.go",
    };

    api.addRoutes(this, routes);
    Object.keys(routes).forEach((route) => {
      api.getFunction(route)?.attachPermissions([table]);
    });

    // Show the endpoint in the output
    this.addOutputs({
      ApiEndpoint: api.url,
      UserTableName: table.tableName,
    });

    this.table = table;
  }
}
