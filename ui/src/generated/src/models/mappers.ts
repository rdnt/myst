import * as coreClient from "@azure/core-client";

export const CreateKeystoreRequest: coreClient.CompositeMapper = {
  type: {
    name: "Composite",
    className: "CreateKeystoreRequest",
    modelProperties: {
      name: {
        serializedName: "name",
        required: true,
        type: {
          name: "String"
        }
      },
      password: {
        serializedName: "password",
        required: true,
        type: {
          name: "String"
        }
      }
    }
  }
};

export const Keystore: coreClient.CompositeMapper = {
  type: {
    name: "Composite",
    className: "Keystore",
    modelProperties: {
      id: {
        serializedName: "id",
        required: true,
        type: {
          name: "String"
        }
      },
      name: {
        serializedName: "name",
        required: true,
        type: {
          name: "String"
        }
      },
      entries: {
        serializedName: "entries",
        required: true,
        type: {
          name: "Sequence",
          element: {
            type: {
              name: "Composite",
              className: "Entry"
            }
          }
        }
      }
    }
  }
};

export const Entry: coreClient.CompositeMapper = {
  type: {
    name: "Composite",
    className: "Entry",
    modelProperties: {
      id: {
        serializedName: "id",
        required: true,
        type: {
          name: "String"
        }
      },
      label: {
        serializedName: "label",
        required: true,
        type: {
          name: "String"
        }
      },
      username: {
        serializedName: "username",
        required: true,
        type: {
          name: "String"
        }
      },
      password: {
        serializedName: "password",
        required: true,
        type: {
          name: "String"
        }
      }
    }
  }
};

export const ErrorModel: coreClient.CompositeMapper = {
  type: {
    name: "Composite",
    className: "ErrorModel",
    modelProperties: {
      code: {
        serializedName: "code",
        required: true,
        type: {
          name: "String"
        }
      },
      message: {
        serializedName: "message",
        required: true,
        type: {
          name: "String"
        }
      }
    }
  }
};

export const UnlockKeystoreRequest: coreClient.CompositeMapper = {
  type: {
    name: "Composite",
    className: "UnlockKeystoreRequest",
    modelProperties: {
      password: {
        serializedName: "password",
        required: true,
        type: {
          name: "String"
        }
      }
    }
  }
};

export const CreateEntryRequest: coreClient.CompositeMapper = {
  type: {
    name: "Composite",
    className: "CreateEntryRequest",
    modelProperties: {
      label: {
        serializedName: "label",
        required: true,
        type: {
          name: "String"
        }
      },
      username: {
        serializedName: "username",
        required: true,
        type: {
          name: "String"
        }
      },
      password: {
        serializedName: "password",
        required: true,
        type: {
          name: "String"
        }
      }
    }
  }
};
