import PocketBase, { LocalAuthStore } from "pocketbase";
import { dataProvider, authProvider, liveProvider } from "refine-pocketbase";

// Reuse the PocketBase admin UI's auth token from localStorage so that
// superusers logged into /_/ are automatically authenticated here.
const store = new LocalAuthStore("__pb_superuser_auth__");

// In dev, Vite proxies /api to PB. In prod, same origin.
const pb = new PocketBase("/", store);

export const pbClient = pb;
export const pbDataProvider = dataProvider(pb);
export const pbAuthProvider = authProvider(pb, {
  collection: "_superusers",
});
export const pbLiveProvider = liveProvider(pb);
