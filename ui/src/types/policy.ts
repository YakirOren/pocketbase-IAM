export interface PolicyStatement {
  sid: string;
  effect: "Allow" | "Deny";
  action: string[];
  resource: string[];
}

export interface PolicyDocument {
  version: string;
  statement: PolicyStatement[];
}

export const DEFAULT_STATEMENT: PolicyStatement = {
  sid: "",
  effect: "Allow",
  action: [""],
  resource: ["*"],
};

export const DEFAULT_POLICY_DOCUMENT: PolicyDocument = {
  version: "2024-01-01",
  statement: [{ ...DEFAULT_STATEMENT }],
};
