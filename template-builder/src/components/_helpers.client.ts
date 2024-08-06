export function SetupClientHelpers() {
  class Hardwire {
    static formHeaders(currentRouter: string, islands: string[]) {
      let presentIslands = "";
      for (let i = 0; i < islands.length; i++) {
        const islandID = islands[i];
        const islandEl = document.getElementById(islandID);
        if (islandEl) {
          presentIslands += "," + islandID;
        }
      }
      return {
        "Hardwire-Islands-Update": presentIslands.slice(1),
        "Hardwire-Dynamic-Fragment-Request": currentRouter,
      };
    }
  }

  Object.defineProperty(window, "__hardwire", {
    value: Hardwire,
    writable: false,
    configurable: false,
    enumerable: false,
  });

  return Hardwire;
}

export type HardwireClientApi = ReturnType<typeof SetupClientHelpers>;
export type ClientApiKeys = Exclude<keyof HardwireClientApi, "prototype">;

export class Client {
  static call<K extends ClientApiKeys>(
    method: K,
    ...args: Parameters<HardwireClientApi[K]>
  ) {
    return `__hardwire.${method}(${args.map((v) => JSON.stringify(v)).join(", ")})`;
  }
}
