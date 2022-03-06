import * as sst from "@serverless-stack/resources";

export interface TaskStackProps extends sst.StackProps {
  api: sst.Api;
  userTable: sst.Table;
}

export class TaskStack extends sst.Stack {
  table: sst.Table;

  constructor(
    scope: sst.App,
    id: string,
    { api, userTable, ...props }: TaskStackProps
  ) {
    super(scope, id, props);

    // Create a DynamoDB table
    const table = new sst.Table(this, "Tasks", {
      primaryIndex: { partitionKey: "id" },
      fields: { id: sst.TableFieldType.STRING },
    });

    this.addDefaultFunctionEnv({
      LOGLEVEL: "debug",
      USERTABLE_NAME: userTable.tableName,
      TASKTABLE_NAME: table.tableName,
    });

    // Create a HTTP API
    const routes: sst.ApiProps["routes"] = {
      "POST /api/task": "lambda/create_task/main.go",
      "GET /api/task": "lambda/get_task/main.go",
      "PATCH /api/task": "lambda/update_task/main.go",
    };

    api.addRoutes(this, routes);
    Object.keys(routes).forEach((route) => {
      api.getFunction(route)?.attachPermissions([table, userTable]);
    });

    // Show the endpoint in the output
    this.addOutputs({
      ApiEndpoint: api.url,
      TaskTableName: table.tableName,
    });

    this.table = table;
  }
}
