import { useState } from "react";
import { useList } from "@refinedev/core";
import { useSearchParams } from "react-router";
import { pbClient } from "@/providers/pocketbase";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { CheckCircle2, XCircle } from "lucide-react";

interface SimulateResult {
  allowed: boolean;
  reason: string;
  matched_statement?: {
    sid: string;
    effect: string;
    policy_name: string;
  };
  trace?: string[];
}

export function Simulator() {
  const [searchParams] = useSearchParams();
  const [userId, setUserId] = useState(searchParams.get("user") ?? "");
  const [action, setAction] = useState("");
  const [resource, setResource] = useState("*");
  const [result, setResult] = useState<SimulateResult | null>(null);

  // Fetch users for dropdown
  const { result: usersResult } = useList({
    resource: "users",
    pagination: { pageSize: 100 },
  });

  // Fetch policies to extract action suggestions
  const { result: policiesResult } = useList({
    resource: "iam_policies",
    pagination: { pageSize: 100 },
  });
  const actionSuggestions = Array.from(
    new Set(
      (policiesResult.data ?? []).flatMap((p: Record<string, unknown>) => {
        const doc = p.document as { statement?: { action?: string[] }[] } | undefined;
        return (doc?.statement ?? []).flatMap((s) => s.action ?? []);
      })
    )
  ).sort();

  const [isPending, setIsPending] = useState(false);

  const handleSimulate = async () => {
    setIsPending(true);
    try {
      const data = await pbClient.send<SimulateResult>("/api/iam/simulate", {
        method: "POST",
        body: { user_id: userId, action, resource },
      });
      setResult(data);
    } catch {
      setResult({ allowed: false, reason: "Request failed", trace: ["Error calling /api/iam/simulate"] });
    } finally {
      setIsPending(false);
    }
  };

  return (
    <div className="max-w-2xl space-y-6">
      <h1 className="text-2xl font-bold">Policy Simulator</h1>
      <p className="text-muted-foreground">
        Test whether a user is allowed to perform an action.
      </p>

      <div className="space-y-4">
        <div className="space-y-2">
          <Label>User</Label>
          <Select value={userId} onValueChange={setUserId}>
            <SelectTrigger>
              <SelectValue placeholder="Select user..." />
            </SelectTrigger>
            <SelectContent>
              {(usersResult.data ?? []).map((u: Record<string, unknown>) => (
                <SelectItem key={u.id as string} value={u.id as string}>
                  {u.email as string}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-2">
          <Label>Action</Label>
          <Input
            value={action}
            onChange={(e) => setAction(e.target.value)}
            placeholder="e.g. collections:read"
            list="action-suggestions"
          />
          <datalist id="action-suggestions">
            {actionSuggestions.map((a) => (
              <option key={a as string} value={a as string} />
            ))}
          </datalist>
        </div>

        <div className="space-y-2">
          <Label>Resource</Label>
          <Input
            value={resource}
            onChange={(e) => setResource(e.target.value)}
            placeholder="e.g. posts, * (default)"
          />
        </div>

        <Button
          onClick={handleSimulate}
          disabled={!userId || !action || isPending}
        >
          {isPending ? "Simulating..." : "Simulate"}
        </Button>
      </div>

      {result && (
        <Card>
          <CardContent className="space-y-4 pt-6">
            <div className="flex items-center gap-2">
              {result.allowed ? (
                <>
                  <CheckCircle2 className="h-6 w-6 text-green-600" />
                  <Badge variant="default" className="text-lg">
                    ALLOWED
                  </Badge>
                </>
              ) : (
                <>
                  <XCircle className="h-6 w-6 text-red-600" />
                  <Badge variant="destructive" className="text-lg">
                    DENIED
                  </Badge>
                </>
              )}
            </div>

            {result.matched_statement && (
              <div className="space-y-1 text-sm">
                <p>
                  <span className="text-muted-foreground">Matched: </span>
                  {result.matched_statement.sid} ({result.matched_statement.effect})
                </p>
                <p>
                  <span className="text-muted-foreground">Policy: </span>
                  {result.matched_statement.policy_name}
                </p>
              </div>
            )}

            {result.trace && (
              <div>
                <h3 className="mb-2 text-sm font-medium">Evaluation Trace</h3>
                <ol className="list-decimal space-y-1 pl-5 text-sm text-muted-foreground">
                  {result.trace.map((step, i) => (
                    <li key={i}>{step}</li>
                  ))}
                </ol>
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  );
}
