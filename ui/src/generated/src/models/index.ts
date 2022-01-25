import * as coreClient from "@azure/core-client";

export interface CreateKeystoreRequest {
  name: string;
  password: string;
}

export interface Keystore {
  id: string;
  name: string;
  entries: Entry[];
}

export interface Entry {
  id: string;
  label: string;
  username: string;
  password: string;
}

export interface ErrorModel {
  code: string;
  message: string;
}

export interface UnlockKeystoreRequest {
  password: string;
}

export interface CreateEntryRequest {
  label: string;
  username: string;
  password: string;
}

/** Optional parameters. */
export interface CreateKeystoreOptionalParams
  extends coreClient.OperationOptions {}

/** Contains response data for the createKeystore operation. */
export type CreateKeystoreResponse = Keystore[];

/** Optional parameters. */
export interface KeystoreIdsOptionalParams
  extends coreClient.OperationOptions {}

/** Contains response data for the keystoreIds operation. */
export type KeystoreIdsResponse = {
  /** The parsed response body. */
  body: string[];
};

/** Optional parameters. */
export interface UnlockKeystoreOptionalParams
  extends coreClient.OperationOptions {}

/** Contains response data for the unlockKeystore operation. */
export type UnlockKeystoreResponse = Keystore[];

/** Optional parameters. */
export interface KeystoreOptionalParams extends coreClient.OperationOptions {}

/** Contains response data for the keystore operation. */
export type KeystoreResponse = Keystore[];

/** Optional parameters. */
export interface CreateEntryOptionalParams
  extends coreClient.OperationOptions {}

/** Contains response data for the createEntry operation. */
export type CreateEntryResponse = Entry;

/** Optional parameters. */
export interface GetKeystoreOptionalParams
  extends coreClient.OperationOptions {}

/** Contains response data for the getKeystore operation. */
export type GetKeystoreResponse = Keystore[];

/** Optional parameters. */
export interface MystRestAPIOptionalParams
  extends coreClient.ServiceClientOptions {
  /** server parameter */
  $host?: string;
  /** Overrides client endpoint. */
  endpoint?: string;
}
