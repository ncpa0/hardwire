export function SetupClientHelpers() {
  class Hardwire {
    static formHeaders(
      currentRouter: string,
      islands: string[],
      items: string[],
    ) {
      let presentIslands = "";
      for (let i = 0; i < islands.length; i++) {
        const islandID = islands[i];
        const islandEl = document.getElementById(islandID);
        if (islandEl) {
          presentIslands += ";" + islandID;
        }
      }

      let listItems = "";
      for (let i = 0; i < items.length; i++) {
        const itemKey = items[i];
        listItems += ";" + itemKey;
      }

      return {
        "Hardwire-Islands-Update": presentIslands.slice(1),
        "Hardwire-Dynamic-Fragment-Request": currentRouter,
        "Hardwire-Dynamic-List-Patch": listItems.slice(1),
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
