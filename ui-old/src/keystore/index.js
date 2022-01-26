class Keystore {
  constructor() {
    // this.debug = true;
  }

  encrypt(data, password) {
    return new Promise((resolve, reject) => {
      window.keystore._encrypt(data, password, (data, error) => {
        if (error) {
          reject(error);
          return;
        }
        resolve(data);
      });
    });
  }
  // eslint-disable-next-line no-unused-vars
  decrypt(data, password) {
    return new Promise((resolve, reject) => {
      window.keystore._encrypt(data, password, (data, error) => {
        if (error) {
          reject(error);
          return;
        }
        resolve(data);
      });
    });
  }
}

export default new Keystore();
