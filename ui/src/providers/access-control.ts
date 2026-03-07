import type { AccessControlProvider } from "@refinedev/core";
import { pbClient } from "./pocketbase";

const ACTION_MAP: Record<string, string> = {
  list: "list",
  create: "create",
  edit: "update",
  delete: "delete",
  show: "view",
};

export const accessControlProvider: AccessControlProvider = {
  can: async ({ resource, action }) => {
    const iamAction = ACTION_MAP[action];
    if (!resource || !iamAction) {
      return { can: true };
    }

    try {
      const data = await pbClient.send<{ allowed: boolean }>(
        "/api/iam/check",
        {
          method: "POST",
          body: {
            action: `collections:${iamAction}`,
            resource: resource,
          },
        },
      );
      return {
        can: data.allowed,
        reason: data.allowed ? undefined : "Access denied by IAM policy",
      };
    } catch {
      return { can: false, reason: "Failed to check permissions" };
    }
  },
  options: {
    buttons: {
      enableAccessControl: true,
      hideIfUnauthorized: true,
    },
  },
};
