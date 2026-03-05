import type { PolicyDocument } from "@/types/policy";

export interface EffectivePermission {
  action: string;
  effect: "Allow" | "Deny";
  source: string; // e.g., "Role: Editor", "Direct: PolicyName", "Group: Devs"
  policyName: string;
  sid: string;
}

interface PolicyWithSource {
  policyName: string;
  document: PolicyDocument;
  source: string;
}

/**
 * Merges policies and computes effective permissions.
 * Deny overrides Allow (AWS IAM model).
 */
export function computeEffectivePermissions(
  policies: PolicyWithSource[]
): EffectivePermission[] {
  const permissions = new Map<string, EffectivePermission>();

  for (const { policyName, document: doc, source } of policies) {
    for (const stmt of doc.statement) {
      for (const action of stmt.action) {
        const existing = permissions.get(action);
        // Deny always wins
        if (existing && existing.effect === "Deny") continue;
        // New Deny overrides existing Allow
        if (stmt.effect === "Deny" || !existing) {
          permissions.set(action, {
            action,
            effect: stmt.effect as "Allow" | "Deny",
            source,
            policyName,
            sid: stmt.sid,
          });
        }
      }
    }
  }

  return Array.from(permissions.values()).sort((a, b) =>
    a.action.localeCompare(b.action)
  );
}
