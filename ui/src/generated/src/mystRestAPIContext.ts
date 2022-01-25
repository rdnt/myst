import * as coreClient from "@azure/core-client";
import { MystRestAPIOptionalParams } from "./models";

export class MystRestAPIContext extends coreClient.ServiceClient {
  $host: string;

  /**
   * Initializes a new instance of the MystRestAPIContext class.
   * @param options The parameter options
   */
  constructor(options?: MystRestAPIOptionalParams) {
    // Initializing default values for options
    if (!options) {
      options = {};
    }
    const defaults: MystRestAPIOptionalParams = {
      requestContentType: "application/json; charset=utf-8"
    };

    const packageDetails = `azsdk-js-mystRestAPI/1.0.0-beta.1`;
    const userAgentPrefix =
      options.userAgentOptions && options.userAgentOptions.userAgentPrefix
        ? `${options.userAgentOptions.userAgentPrefix} ${packageDetails}`
        : `${packageDetails}`;

    const optionsWithDefaults = {
      ...defaults,
      ...options,
      userAgentOptions: {
        userAgentPrefix
      },
      baseUri: options.endpoint || "http://localhost:8081/api"
    };
    super(optionsWithDefaults);

    // Assigning values to Constant parameters
    this.$host = options.$host || "http://localhost:8081/api";
  }
}
