import ColorHash from "color-hash";

const ch = new ColorHash({lightness: 0.6, saturation: 0.3});

export function hash(s: string): string {
  return ch.hex(s);
}
