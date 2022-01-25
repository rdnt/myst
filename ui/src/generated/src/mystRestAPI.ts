import * as coreClient from "@azure/core-client";
import * as Parameters from "./models/parameters";
import * as Mappers from "./models/mappers";
import { MystRestAPIContext } from "./mystRestAPIContext";
import {
  MystRestAPIOptionalParams,
  CreateKeystoreRequest,
  CreateKeystoreOptionalParams,
  CreateKeystoreResponse,
  KeystoreIdsOptionalParams,
  KeystoreIdsResponse,
  UnlockKeystoreRequest,
  UnlockKeystoreOptionalParams,
  UnlockKeystoreResponse,
  KeystoreOptionalParams,
  KeystoreResponse,
  CreateEntryRequest,
  CreateEntryOptionalParams,
  CreateEntryResponse,
  GetKeystoreOptionalParams,
  GetKeystoreResponse
} from "./models";

export class MystRestAPI extends MystRestAPIContext {
  /**
   * Initializes a new instance of the MystRestAPI class.
   * @param options The parameter options
   */
  constructor(options?: MystRestAPIOptionalParams) {
    super(options);
  }

  /**
   * Creates a new encrypted keystore with the given password
   * @param body
   * @param options The options parameters.
   */
  createKeystore(
    body: CreateKeystoreRequest,
    options?: CreateKeystoreOptionalParams
  ): Promise<CreateKeystoreResponse> {
    return this.sendOperationRequest(
      { body, options },
      createKeystoreOperationSpec
    );
  }

  /**
   * Returns all keystore ids
   * @param options The options parameters.
   */
  keystoreIds(
    options?: KeystoreIdsOptionalParams
  ): Promise<KeystoreIdsResponse> {
    return this.sendOperationRequest({ options }, keystoreIdsOperationSpec);
  }

  /**
   * Unlocks a keystore if it exists and returns it
   * @param keystoreId unique identifier for a keystore
   * @param body
   * @param options The options parameters.
   */
  unlockKeystore(
    keystoreId: string,
    body: UnlockKeystoreRequest,
    options?: UnlockKeystoreOptionalParams
  ): Promise<UnlockKeystoreResponse> {
    return this.sendOperationRequest(
      { keystoreId, body, options },
      unlockKeystoreOperationSpec
    );
  }

  /**
   * Get a keystore if it exists and return it
   * @param keystoreId unique identifier for a keystore
   * @param options The options parameters.
   */
  keystore(
    keystoreId: string,
    options?: KeystoreOptionalParams
  ): Promise<KeystoreResponse> {
    return this.sendOperationRequest(
      { keystoreId, options },
      keystoreOperationSpec
    );
  }

  /**
   * Creates a new entry and adds it to the keystore
   * @param keystoreId unique identifier for a keystore
   * @param body
   * @param options The options parameters.
   */
  createEntry(
    keystoreId: string,
    body: CreateEntryRequest,
    options?: CreateEntryOptionalParams
  ): Promise<CreateEntryResponse> {
    return this.sendOperationRequest(
      { keystoreId, body, options },
      createEntryOperationSpec
    );
  }

  /**
   * Get a keystore if it exists and returns it
   * @param keystoreId unique identifier for a keystore
   * @param options The options parameters.
   */
  getKeystore(
    keystoreId: string,
    options?: GetKeystoreOptionalParams
  ): Promise<GetKeystoreResponse> {
    return this.sendOperationRequest(
      { keystoreId, options },
      getKeystoreOperationSpec
    );
  }
}
// Operation Specifications
const serializer = coreClient.createSerializer(Mappers, /* isXml */ false);

const createKeystoreOperationSpec: coreClient.OperationSpec = {
  path: "/keystores",
  httpMethod: "POST",
  responses: {
    200: {
      bodyMapper: {
        type: {
          name: "Sequence",
          element: { type: { name: "Composite", className: "Keystore" } }
        }
      }
    },
    default: {
      bodyMapper: Mappers.ErrorModel
    }
  },
  requestBody: Parameters.body,
  urlParameters: [Parameters.$host],
  headerParameters: [Parameters.contentType, Parameters.accept],
  mediaType: "json",
  serializer
};
const keystoreIdsOperationSpec: coreClient.OperationSpec = {
  path: "/keystores",
  httpMethod: "GET",
  responses: {
    200: {
      bodyMapper: {
        type: { name: "Sequence", element: { type: { name: "String" } } }
      }
    },
    default: {
      bodyMapper: Mappers.ErrorModel
    }
  },
  urlParameters: [Parameters.$host],
  headerParameters: [Parameters.accept],
  serializer
};
const unlockKeystoreOperationSpec: coreClient.OperationSpec = {
  path: "/keystore/{keystoreId}",
  httpMethod: "POST",
  responses: {
    200: {
      bodyMapper: {
        type: {
          name: "Sequence",
          element: { type: { name: "Composite", className: "Keystore" } }
        }
      }
    },
    default: {
      bodyMapper: Mappers.ErrorModel
    }
  },
  requestBody: Parameters.body1,
  urlParameters: [Parameters.$host, Parameters.keystoreId],
  headerParameters: [Parameters.contentType, Parameters.accept],
  mediaType: "json",
  serializer
};
const keystoreOperationSpec: coreClient.OperationSpec = {
  path: "/keystore/{keystoreId}",
  httpMethod: "GET",
  responses: {
    200: {
      bodyMapper: {
        type: {
          name: "Sequence",
          element: { type: { name: "Composite", className: "Keystore" } }
        }
      }
    },
    default: {
      bodyMapper: Mappers.ErrorModel
    }
  },
  urlParameters: [Parameters.$host, Parameters.keystoreId],
  headerParameters: [Parameters.accept],
  serializer
};
const createEntryOperationSpec: coreClient.OperationSpec = {
  path: "/keystore/{keystoreId}/entries",
  httpMethod: "POST",
  responses: {
    200: {
      bodyMapper: Mappers.Entry
    },
    default: {
      bodyMapper: Mappers.ErrorModel
    }
  },
  requestBody: Parameters.body2,
  urlParameters: [Parameters.$host, Parameters.keystoreId],
  headerParameters: [Parameters.contentType, Parameters.accept],
  mediaType: "json",
  serializer
};
const getKeystoreOperationSpec: coreClient.OperationSpec = {
  path: "/keystore/{keystoreId}/entries",
  httpMethod: "GET",
  responses: {
    200: {
      bodyMapper: {
        type: {
          name: "Sequence",
          element: { type: { name: "Composite", className: "Keystore" } }
        }
      }
    },
    default: {
      bodyMapper: Mappers.ErrorModel
    }
  },
  urlParameters: [Parameters.$host, Parameters.keystoreId],
  headerParameters: [Parameters.accept],
  serializer
};
