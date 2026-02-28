import PocketBase from "pocketbase";
import { dataProvider, authProvider, liveProvider } from "refine-pocketbase";

// In dev, Vite proxies /api to PB. In prod, same origin.
const pb = new PocketBase("/");

export const pbClient = pb;
export const pbDataProvider = dataProvider(pb);
export const pbAuthProvider = authProvider(pb, {
  collection: "users",
});
export const pbLiveProvider = liveProvider(pb);
