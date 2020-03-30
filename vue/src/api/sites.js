/**
 * Mocking client-server processing
 */
const _sites = [
  { id: 0 },
  { id: 1 },
  { id: 2 },
  { id: 3 },
  { id: 4 },
  { id: 5 }
];

export default {
  getSites(cb) {
    setTimeout(() => cb(_sites), 1000);
  }
};
