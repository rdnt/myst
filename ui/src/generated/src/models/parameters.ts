import { OperationParameter, OperationURLParameter } from "@azure/core-client";
import {
  CreateKeystoreRequest as CreateKeystoreRequestMapper,
  UnlockKeystoreRequest as UnlockKeystoreRequestMapper,
  CreateEntryRequest as CreateEntryRequestMapper
} from "../models/mappers";

export const contentType: OperationParameter = {
  parameterPath: ["options", "contentType"],
  mapper: {
    defaultValue: "application/json",
    isConstant: true,
    serializedName: "Content-Type",
    type: {
      name: "String"
    }
  }
};

export const body: OperationParameter = {
  parameterPath: "body",
  mapper: CreateKeystoreRequestMapper
};

export const accept: OperationParameter = {
  parameterPath: "accept",
  mapper: {
    defaultValue: "application/json",
    isConstant: true,
    serializedName: "Accept",
    type: {
      name: "String"
    }
  }
};

export const $host: OperationURLParameter = {
  parameterPath: "$host",
  mapper: {
    serializedName: "$host",
    required: true,
    type: {
      name: "String"
    }
  },
  skipEncoding: true
};

export const body1: OperationParameter = {
  parameterPath: "body",
  mapper: UnlockKeystoreRequestMapper
};

export const keystoreId: OperationURLParameter = {
  parameterPath: "keystoreId",
  mapper: {
    serializedName: "keystoreId",
    required: true,
    type: {
      name: "String"
    }
  }
};

export const body2: OperationParameter = {
  parameterPath: "body",
  mapper: CreateEntryRequestMapper
};
